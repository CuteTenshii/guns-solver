package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
// and returns its cookies. Solverify navigates the target itself, so it needs
// no captured HTML — but it does require the proxy the cookie will be bound to,
// and the method the challenged request used. Ports are sent as numeric
// strings, which is what the API expects.
type solverifyInterstitialTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	UserAgent  string `json:"useragent,omitempty"`
	// Source asks for the landed page's HTML in solution.html. It is the only
	// way to tell a session that was never challenged from one that failed the
	// challenge, so it is always requested and reported when no cookie comes
	// back.
	Source        bool   `json:"source"`
	Method        string `json:"method,omitempty"`
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
		Value     string           `json:"value"`
		Cookies   solverifyCookies `json:"cookies"`
		UserAgent string           `json:"useragent"`
		HTML      string           `json:"html"`
	} `json:"solution"`
}

// solverifyCookies is the solved session's cookie jar, documented as a
// name-to-value object. Cookie jars are also commonly serialised as a list of
// cookie objects, so both shapes are accepted rather than failing the whole
// decode over the difference.
type solverifyCookies map[string]string

func (c *solverifyCookies) UnmarshalJSON(data []byte) error {
	var byName map[string]string
	if err := json.Unmarshal(data, &byName); err == nil {
		*c = byName
		return nil
	}

	var list []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("cookies is neither an object nor a list: %w", err)
	}
	jar := make(solverifyCookies, len(list))
	for _, cookie := range list {
		jar[cookie.Name] = cookie.Value
	}
	*c = jar
	return nil
}

// describeLandedPage summarises the page the solver's browser ended up on.
// That is what separates "Cloudflare never challenged this request" from "the
// request never reached the challenge in the first place".
func describeLandedPage(page string) string {
	if page == "" {
		return "no page captured"
	}
	lower := strings.ToLower(page)
	for _, marker := range []string{"just a moment", "cf-chl", "challenge-platform"} {
		if strings.Contains(lower, marker) {
			return "a Cloudflare challenge page"
		}
	}
	return truncate(strings.Join(strings.Fields(page), " "), 160)
}

// names lists the cookies present, for error messages that would otherwise not
// say why a cookie is missing.
func (c solverifyCookies) names() string {
	if len(c) == 0 {
		return "none"
	}
	names := make([]string, 0, len(c))
	for name := range c {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
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

	// Navigate the request that actually carried the challenge — same URL and
	// same method. A real browser is served the profile page unchallenged, and
	// the analytics endpoint only answers POST, so a GET there misses the rule
	// too; either way an unchallenged session is issued no cf_clearance.
	target := r.ChallengeURL
	if target == "" {
		target = r.PageURL
	}
	method := strings.ToUpper(r.ChallengeMethod)
	switch method {
	case "":
		method = http.MethodGet
	case http.MethodGet, http.MethodPost:
	default:
		return clearanceResult{}, fmt.Errorf("solverify navigates with GET or POST only, got %q", r.ChallengeMethod)
	}

	taskID, err := p.createTask(ctx, solverifyInterstitialTask{
		Type:          "interstitial",
		WebsiteURL:    target,
		UserAgent:     r.UserAgent,
		Source:        true,
		Method:        method,
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
			// Cloudflare only mints cf_clearance when it actually challenges,
			// so an unchallenged navigation completes with a cookie jar that
			// simply has none. Name what did come back — that distinguishes it
			// from a genuinely failed solve.
			return false, fmt.Errorf("solverify navigated %s %s without being issued a cf_clearance cookie "+
				"(cookies: %s; landed on: %s)",
				method, target, res.Solution.Cookies.names(), describeLandedPage(res.Solution.HTML))
		}
		result = clearanceResult{Cookie: cookie, UserAgent: res.Solution.UserAgent}
		return true, nil
	})
	if err != nil {
		return clearanceResult{}, err
	}
	return result, nil
}
