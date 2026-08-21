package main

import (
	"context"
	"encoding/base64"
	"fmt"
)

const capmonsterURL = "https://api.capmonster.cloud"

// capmonsterProvider solves captchas through CapMonster Cloud.
type capmonsterProvider struct {
	apiKey string
}

func (p *capmonsterProvider) Name() string { return providerCapmonster }

type capmonsterCreateRequest struct {
	ClientKey string `json:"clientKey"`
	Task      any    `json:"task"`
}

// capmonsterTurnstileTask solves the Turnstile widget itself. It runs on
// CapMonster's own infrastructure: the token is not bound to our IP, so no
// proxy is passed.
type capmonsterTurnstileTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction"`
	Data       string `json:"data"`
}

// capmonsterClearanceTask is the Cloudflare Challenge job that returns a
// cf_clearance cookie. Unlike the Turnstile widget task, it bypasses the
// Managed Challenge gating an endpoint and requires our own proxy — the
// resulting cookie is bound to that proxy's IP and to the user-agent.
type capmonsterClearanceTask struct {
	Type               string `json:"type"`
	WebsiteURL         string `json:"websiteURL"`
	WebsiteKey         string `json:"websiteKey"`
	CloudflareTaskType string `json:"cloudflareTaskType"`
	HTMLPageBase64     string `json:"htmlPageBase64"`
	UserAgent          string `json:"userAgent"`
	ProxyType          string `json:"proxyType"`
	ProxyAddress       string `json:"proxyAddress"`
	ProxyPort          int    `json:"proxyPort"`
	ProxyLogin         string `json:"proxyLogin,omitempty"`
	ProxyPassword      string `json:"proxyPassword,omitempty"`
}

type capmonsterCreateResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskID           int    `json:"taskId"`
}

type capmonsterResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

type capmonsterResultResponse struct {
	ErrorID          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
	Solution         struct {
		Token       string `json:"token"`
		CfClearance string `json:"cf_clearance"`
	} `json:"solution"`
}

// capmonsterReady is the status CapMonster reports once a task has a solution.
const capmonsterReady = "ready"

// capmonsterError renders a non-zero errorId with whatever detail CapMonster
// supplied alongside it.
func capmonsterError(op string, errorID int, code, description string) error {
	detail := code
	if description != "" {
		if detail != "" {
			detail += ": "
		}
		detail += description
	}
	if detail == "" {
		return fmt.Errorf("capmonster %s error %d", op, errorID)
	}
	return fmt.Errorf("capmonster %s error %d (%s)", op, errorID, detail)
}

// createTask submits a task and returns its id.
func (p *capmonsterProvider) createTask(ctx context.Context, task any) (int, error) {
	var resp capmonsterCreateResponse
	if err := postJSON(ctx, capmonsterURL+"/createTask",
		capmonsterCreateRequest{ClientKey: p.apiKey, Task: task}, &resp); err != nil {
		return 0, err
	}
	if resp.ErrorID != 0 {
		return 0, capmonsterError("createTask", resp.ErrorID, resp.ErrorCode, resp.ErrorDescription)
	}
	return resp.TaskID, nil
}

// taskResult fetches a task's current state.
func (p *capmonsterProvider) taskResult(ctx context.Context, taskID int) (*capmonsterResultResponse, error) {
	var resp capmonsterResultResponse
	if err := postJSON(ctx, capmonsterURL+"/getTaskResult",
		capmonsterResultRequest{ClientKey: p.apiKey, TaskID: taskID}, &resp); err != nil {
		return nil, err
	}
	if resp.ErrorID != 0 {
		return nil, capmonsterError("getTaskResult", resp.ErrorID, resp.ErrorCode, resp.ErrorDescription)
	}
	return &resp, nil
}

func (p *capmonsterProvider) SolveTurnstile(ctx context.Context, r turnstileRequest) (string, error) {
	taskID, err := p.createTask(ctx, capmonsterTurnstileTask{
		Type:       "TurnstileTaskProxyless",
		WebsiteURL: r.PageURL,
		WebsiteKey: r.SiteKey,
		PageAction: r.Action,
		Data:       r.CData,
	})
	if err != nil {
		return "", err
	}

	var token string
	err = pollTask(ctx, func() (bool, error) {
		res, err := p.taskResult(ctx, taskID)
		if err != nil {
			return false, err
		}
		if res.Status != capmonsterReady {
			return false, nil
		}
		if res.Solution.Token == "" {
			return false, fmt.Errorf("capmonster returned an empty turnstile token")
		}
		token = res.Solution.Token
		return true, nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (p *capmonsterProvider) SolveCfClearance(ctx context.Context, r clearanceRequest) (clearanceResult, error) {
	if r.ProxyURL == "" {
		return clearanceResult{}, fmt.Errorf("cf_clearance requires a proxy (pass -proxy)")
	}
	proxy, err := parseProxyParts(r.ProxyURL)
	if err != nil {
		return clearanceResult{}, err
	}

	taskID, err := p.createTask(ctx, capmonsterClearanceTask{
		Type:               "TurnstileTask",
		WebsiteURL:         r.PageURL,
		WebsiteKey:         "xxxxxx",
		CloudflareTaskType: "cf_clearance",
		HTMLPageBase64:     base64.StdEncoding.EncodeToString(r.HTMLPage),
		UserAgent:          r.UserAgent,
		ProxyType:          proxy.Scheme,
		ProxyAddress:       proxy.Host,
		ProxyPort:          proxy.Port,
		ProxyLogin:         proxy.Login,
		ProxyPassword:      proxy.Password,
	})
	if err != nil {
		return clearanceResult{}, err
	}

	var cookie string
	err = pollTask(ctx, func() (bool, error) {
		res, err := p.taskResult(ctx, taskID)
		if err != nil {
			return false, err
		}
		if res.Status != capmonsterReady {
			return false, nil
		}
		if res.Solution.CfClearance == "" {
			return false, fmt.Errorf("capmonster returned an empty cf_clearance")
		}
		cookie = res.Solution.CfClearance
		return true, nil
	})
	if err != nil {
		return clearanceResult{}, err
	}
	// CapMonster solves with the user-agent it is handed and reports none back.
	return clearanceResult{Cookie: cookie}, nil
}
