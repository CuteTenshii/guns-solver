package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const uncaptchaURL = "https://api.uncaptcha.io/v1"

// The uncaptcha task types this tool uses. The service also offers `wafauto`,
// which makes the post-challenge request itself; that cannot be used here
// because the request being cleared has to carry our own signed view payload.
const (
	uncaptchaTaskTurnstile = "turnstile"
	uncaptchaTaskWAF       = "waf"
)

// uncaptchaMaxAttempts bounds how many times a retryable failure is re-sent
// before the solve is abandoned.
const uncaptchaMaxAttempts = 3

// uncaptchaRetryBackoff is the wait before re-sending a retryable failure; it
// doubles with each attempt. A 429 means the key's concurrency limit is already
// saturated, so the first wait has to be long enough for an in-flight solve to
// land — with -concurrency above the key's limit that is the common case.
const uncaptchaRetryBackoff = 5 * time.Second

// uncaptchaProvider solves captchas through uncaptcha.io.
//
// Its API is synchronous, unlike CapMonster's and Solverify's: a single request
// carries the task and blocks until the solution is ready or the service's
// 30-second execution limit is reached. There is no task id and nothing to
// poll, so this provider uses neither pollTask nor a task-result endpoint.
type uncaptchaProvider struct {
	apiKey string
}

func (p *uncaptchaProvider) Name() string { return providerUncaptcha }

type uncaptchaExecuteRequest struct {
	TaskType string `json:"task_type"`
	TaskData any    `json:"task_data"`
}

// uncaptchaTurnstileTask solves a standalone Turnstile widget. As with the
// other providers no proxy is sent: the token is not bound to our IP.
type uncaptchaTurnstileTask struct {
	URL     string `json:"url"`
	SiteKey string `json:"sitekey"`
	Action  string `json:"action,omitempty"`
	CData   string `json:"cdata,omitempty"`
}

// uncaptchaWAFTask clears a Cloudflare challenge and returns the identity to
// replay the challenged request with. The proxy is required — the clearance is
// bound to the IP that solved the challenge — and the supplied HTML must be the
// interstitial that same proxy and User-Agent were served, or the service
// rejects the solve as a mismatched client identity.
type uncaptchaWAFTask struct {
	URL       string `json:"url"`
	Proxy     string `json:"proxy"`
	UserAgent string `json:"user_agent,omitempty"`
	HTML      string `json:"html,omitempty"`
}

// uncaptchaResponse is the envelope every uncaptcha endpoint answers with.
// Failures carry success:false and a human-readable message in place of data.
type uncaptchaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Solution struct {
			// Token is the solved Turnstile response token.
			Token string `json:"token"`
			// Clearance is the cf_clearance value from a WAF solve.
			Clearance string `json:"clearance"`
			// Headers are the browser headers the WAF challenge was solved
			// with, including the User-Agent the clearance is bound to.
			Headers map[string]string `json:"headers"`
		} `json:"solution"`
	} `json:"data"`
}

// solvedUserAgent picks the User-Agent out of the headers a WAF solve reports.
// Header names are case-insensitive on the wire, so the documented lowercase
// key is matched loosely rather than relied on — missing it would silently skip
// the identity-mismatch warning in SubmitSolution.
func solvedUserAgent(headers map[string]string) string {
	for name, value := range headers {
		if strings.EqualFold(name, "user-agent") {
			return value
		}
	}
	return ""
}

// uncaptchaRetryable reports whether a failed execute is worth re-sending. The
// documented retryable statuses are a saturated concurrency limit (429), a
// solver-side error (500) and the execution limit (504); any other 5xx is a
// transient gateway failure and is treated the same way. None of them consume
// balance — a failed or timed-out solve is refunded — so a retry is free.
func uncaptchaRetryable(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

// uncaptchaError renders a failed execute. The service reports failures through
// the HTTP status line and a human-readable message rather than an error code.
func uncaptchaError(taskType string, status int, message string) error {
	if message == "" {
		message = "no message"
	}
	return fmt.Errorf("uncaptcha %s task failed (status %d): %s", taskType, status, message)
}

// execute sends one task and returns its solution, retrying transient failures
// with a doubling backoff.
func (p *uncaptchaProvider) execute(ctx context.Context, taskType string, taskData any) (*uncaptchaResponse, error) {
	backoff := uncaptchaRetryBackoff
	var lastErr error

	for attempt := 1; attempt <= uncaptchaMaxAttempts; attempt++ {
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		var resp uncaptchaResponse
		status, err := postJSONWithHeaders(ctx, uncaptchaURL+"/task/execute",
			map[string]string{"X-Api-Key": p.apiKey},
			uncaptchaExecuteRequest{TaskType: taskType, TaskData: taskData}, &resp)

		switch {
		case err != nil:
			lastErr = err
		case !resp.Success:
			lastErr = uncaptchaError(taskType, status, resp.Message)
		default:
			return &resp, nil
		}

		// A cancelled or expired context will not be any healthier next time,
		// and a permanent rejection (bad key, no balance, blacklisted URL) only
		// wastes the caller's remaining time budget if it is repeated.
		if ctx.Err() != nil || !uncaptchaRetryable(status) {
			return nil, lastErr
		}
	}
	return nil, fmt.Errorf("uncaptcha %s task failed after %d attempts: %w", taskType, uncaptchaMaxAttempts, lastErr)
}

func (p *uncaptchaProvider) SolveTurnstile(ctx context.Context, r turnstileRequest) (string, error) {
	resp, err := p.execute(ctx, uncaptchaTaskTurnstile, uncaptchaTurnstileTask{
		URL:     r.PageURL,
		SiteKey: r.SiteKey,
		Action:  r.Action,
		CData:   r.CData,
	})
	if err != nil {
		return "", err
	}
	if resp.Data.Solution.Token == "" {
		return "", fmt.Errorf("uncaptcha reported success but returned an empty turnstile token")
	}
	return resp.Data.Solution.Token, nil
}

func (p *uncaptchaProvider) SolveCfClearance(ctx context.Context, r clearanceRequest) (clearanceResult, error) {
	if r.ProxyURL == "" {
		return clearanceResult{}, fmt.Errorf("cf_clearance requires a proxy (pass -proxy)")
	}
	// uncaptcha takes the proxy as a single URL rather than split fields, but
	// validating it here turns a malformed value into a local error instead of
	// a round-trip that fails inside the service.
	if _, err := parseProxyParts(r.ProxyURL); err != nil {
		return clearanceResult{}, err
	}

	// Clear the request that was actually challenged — the analytics POST, not
	// the profile page. Unlike Solverify, uncaptcha is given the interstitial
	// we were served rather than navigating the target itself, so the method
	// never has to be replayed here.
	target := r.ChallengeURL
	if target == "" {
		target = r.PageURL
	}

	resp, err := p.execute(ctx, uncaptchaTaskWAF, uncaptchaWAFTask{
		URL:       target,
		Proxy:     r.ProxyURL,
		UserAgent: r.UserAgent,
		HTML:      string(r.HTMLPage),
	})
	if err != nil {
		return clearanceResult{}, err
	}

	solution := resp.Data.Solution
	if solution.Clearance == "" {
		return clearanceResult{}, fmt.Errorf("uncaptcha reported success for %s but returned no clearance cookie", target)
	}
	return clearanceResult{Cookie: solution.Clearance, UserAgent: solvedUserAgent(solution.Headers)}, nil
}
