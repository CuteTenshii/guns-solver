package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// capmonsterClearanceTask is the "task" object for a CapMonster Cloudflare
// Challenge job that returns a cf_clearance cookie (cloudflareTaskType
// "cf_clearance"). Unlike the Turnstile widget task, this one bypasses the
// Managed Challenge gating an endpoint and requires your own proxy — the
// resulting cookie is bound to that proxy's IP and to userAgent.
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

type capmonsterClearanceCreateRequest struct {
	ClientKey string                  `json:"clientKey"`
	Task      capmonsterClearanceTask `json:"task"`
}

type capmonsterClearanceResultResponse struct {
	ErrorID  int    `json:"errorId"`
	Status   string `json:"status"`
	Solution struct {
		CfClearance string `json:"cf_clearance"`
	} `json:"solution"`
}

// capmonsterProxy holds a proxy split into the discrete fields CapMonster
// expects, since it cannot accept inline credentials in a single URL.
type capmonsterProxy struct {
	Type     string
	Address  string
	Port     int
	Login    string
	Password string
}

// parseCapmonsterProxy splits a proxy URL (e.g. http://user:pass@host:port)
// into the discrete fields CapMonster's Cloudflare Challenge task requires.
func parseCapmonsterProxy(rawURL string) (capmonsterProxy, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return capmonsterProxy{}, fmt.Errorf("parse proxy url: %w", err)
	}

	p := capmonsterProxy{Type: u.Scheme, Address: u.Hostname()}
	if p.Type == "" || p.Address == "" {
		return capmonsterProxy{}, fmt.Errorf("proxy url %q must include scheme and host", rawURL)
	}

	port := u.Port()
	if port == "" {
		return capmonsterProxy{}, fmt.Errorf("proxy url %q must include a port", rawURL)
	}
	p.Port, err = strconv.Atoi(port)
	if err != nil {
		return capmonsterProxy{}, fmt.Errorf("invalid proxy port %q: %w", port, err)
	}

	if u.User != nil {
		p.Login = u.User.Username()
		if pw, ok := u.User.Password(); ok {
			p.Password = pw
		}
	}
	return p, nil
}

// SolveCfClearance asks CapMonster to solve the Cloudflare Managed Challenge for
// pageURL and returns a cf_clearance cookie. htmlPage is the raw "Just a moment"
// interstitial (status 403) captured from a direct request; proxyURL is the same
// upstream proxy the tool egresses through, so the cookie's bound IP matches.
func SolveCfClearance(ctx context.Context, apiKey, pageURL string, htmlPage []byte, proxyURL string) (string, error) {
	if proxyURL == "" {
		return "", fmt.Errorf("cf_clearance requires a proxy (pass -proxy)")
	}
	proxy, err := parseCapmonsterProxy(proxyURL)
	if err != nil {
		return "", err
	}

	task := capmonsterClearanceTask{
		Type:               "TurnstileTask",
		WebsiteURL:         pageURL,
		WebsiteKey:         "xxxxxx",
		CloudflareTaskType: "cf_clearance",
		HTMLPageBase64:     base64.StdEncoding.EncodeToString(htmlPage),
		UserAgent:          userAgent,
		ProxyType:          proxy.Type,
		ProxyAddress:       proxy.Address,
		ProxyPort:          proxy.Port,
		ProxyLogin:         proxy.Login,
		ProxyPassword:      proxy.Password,
	}

	createBody, _ := json.Marshal(capmonsterClearanceCreateRequest{ClientKey: apiKey, Task: task})
	req, err := http.NewRequestWithContext(ctx, "POST", capmonsterURL+"/createTask", bytes.NewReader(createBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var createResp capmonsterCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return "", err
	}
	if createResp.ErrorID != 0 {
		return "", fmt.Errorf("capmonster createTask error %d", createResp.ErrorID)
	}

	resultBody, _ := json.Marshal(capmonsterResultRequest{ClientKey: apiKey, TaskID: createResp.TaskID})

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(3 * time.Second):
		}

		req, err = http.NewRequestWithContext(ctx, "POST", capmonsterURL+"/getTaskResult", bytes.NewReader(resultBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		var result capmonsterClearanceResultResponse
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return "", fmt.Errorf("capmonster getTaskResult decode: %w", err)
		}
		resp.Body.Close()

		if result.ErrorID != 0 {
			return "", fmt.Errorf("capmonster getTaskResult error %d", result.ErrorID)
		}
		if result.Status == "ready" {
			if result.Solution.CfClearance == "" {
				return "", fmt.Errorf("capmonster returned empty cf_clearance")
			}
			return result.Solution.CfClearance, nil
		}
	}
}
