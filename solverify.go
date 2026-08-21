package main

import (
	"context"
	"fmt"
	"strconv"
)

const solverifyURL = "https://solver.solverify.net"

// solverifyProvider solves captchas through Solverify.
type solverifyProvider struct {
	apiKey string
}

func (p *solverifyProvider) Name() string { return providerSolverify }

type solverifyCreateRequest struct {
	ClientKey string `json:"clientKey"`
	Task      any    `json:"task"`
}

// solverifyTurnstileTask solves the Turnstile widget on Solverify's own
// infrastructure; the token is not IP-bound, so no proxy is passed.
type solverifyTurnstileTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	Action     string `json:"action,omitempty"`
	CData      string `json:"cdata,omitempty"`
}

// solverifyInterstitialTask clears a Cloudflare interstitial in a real browser
// and returns its cookies. Solverify navigates the target page itself, so it
// needs no captured HTML — but it does require the proxy the cookie will be
// bound to. Ports are sent as numeric strings, which is what the API expects.
type solverifyInterstitialTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	UserAgent     string `json:"useragent,omitempty"`
	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     string `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type solverifyCreateResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskID           string `json:"taskId"`
}

type solverifyResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskID    string `json:"taskId"`
}

type solverifyResultResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
	Solution         struct {
		Value     string            `json:"value"`
		Cookies   map[string]string `json:"cookies"`
		UserAgent string            `json:"useragent"`
	} `json:"solution"`
}

// Solverify task states. Anything other than these two means the task is still
// queued or running.
const (
	solverifyCompleted = "completed"
	solverifyFailed    = "failed"
)

// solverifyError renders an API-level failure. Solverify answers with HTTP 200
// even for errors, signalling them through errorId.
func solverifyError(op, code, description string) error {
	detail := code
	if description != "" {
		if detail != "" {
			detail += ": "
		}
		detail += description
	}
	if detail == "" {
		detail = "unknown error"
	}
	return fmt.Errorf("solverify %s failed (%s)", op, detail)
}

// createTask submits a task and returns its id.
func (p *solverifyProvider) createTask(ctx context.Context, task any) (string, error) {
	var resp solverifyCreateResponse
	if err := postJSON(ctx, solverifyURL+"/createTask",
		solverifyCreateRequest{ClientKey: p.apiKey, Task: task}, &resp); err != nil {
		return "", err
	}
	if resp.ErrorID != 0 {
		return "", solverifyError("createTask", resp.ErrorCode, resp.ErrorDescription)
	}
	if resp.TaskID == "" {
		return "", fmt.Errorf("solverify createTask returned no task id")
	}
	return resp.TaskID, nil
}

// taskResult fetches a task's current state, reporting whether it has finished.
func (p *solverifyProvider) taskResult(ctx context.Context, taskID string) (*solverifyResultResponse, bool, error) {
	var resp solverifyResultResponse
	if err := postJSON(ctx, solverifyURL+"/getTaskResult",
		solverifyResultRequest{ClientKey: p.apiKey, TaskID: taskID}, &resp); err != nil {
		return nil, false, err
	}
	if resp.ErrorID != 0 || resp.Status == solverifyFailed {
		return nil, false, solverifyError("getTaskResult", resp.ErrorCode, resp.ErrorDescription)
	}
	return &resp, resp.Status == solverifyCompleted, nil
}

func (p *solverifyProvider) SolveTurnstile(ctx context.Context, r turnstileRequest) (string, error) {
	taskID, err := p.createTask(ctx, solverifyTurnstileTask{
		Type:       "turnstile",
		WebsiteURL: r.PageURL,
		WebsiteKey: r.SiteKey,
		Action:     r.Action,
		CData:      r.CData,
	})
	if err != nil {
		return "", err
	}

	var token string
	err = pollTask(ctx, func() (bool, error) {
		res, done, err := p.taskResult(ctx, taskID)
		if err != nil || !done {
			return false, err
		}
		if res.Solution.Value == "" {
			return false, fmt.Errorf("solverify returned an empty turnstile token")
		}
		token = res.Solution.Value
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *solverifyProvider) SolveCfClearance(ctx context.Context, r clearanceRequest) (clearanceResult, error) {
	if r.ProxyURL == "" {
		return clearanceResult{}, fmt.Errorf("cf_clearance requires a proxy (pass -proxy)")
	}
	proxy, err := parseProxyParts(r.ProxyURL)
	if err != nil {
		return clearanceResult{}, err
	}
	if proxy.Scheme != "http" {
		return clearanceResult{}, fmt.Errorf("solverify only accepts http proxies, got %q", proxy.Scheme)
	}

	taskID, err := p.createTask(ctx, solverifyInterstitialTask{
		Type:          "interstitial",
		WebsiteURL:    r.PageURL,
		UserAgent:     r.UserAgent,
		ProxyType:     proxy.Scheme,
		ProxyAddress:  proxy.Host,
		ProxyPort:     strconv.Itoa(proxy.Port),
		ProxyLogin:    proxy.Login,
		ProxyPassword: proxy.Password,
	})
	if err != nil {
		return clearanceResult{}, err
	}

	var result clearanceResult
	err = pollTask(ctx, func() (bool, error) {
		res, done, err := p.taskResult(ctx, taskID)
		if err != nil || !done {
			return false, err
		}
		cookie := res.Solution.Cookies["cf_clearance"]
		if cookie == "" {
			return false, fmt.Errorf("solverify solved the interstitial but returned no cf_clearance cookie")
		}
		result = clearanceResult{Cookie: cookie, UserAgent: res.Solution.UserAgent}
		return true, nil
	})
	if err != nil {
		return clearanceResult{}, err
	}
	return result, nil
}
