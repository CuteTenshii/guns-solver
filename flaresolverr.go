package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type flareSolverrRequest struct {
	Cmd        string      `json:"cmd"`
	URL        string      `json:"url"`
	PostData   string      `json:"postData"`
	MaxTimeout int         `json:"maxTimeout"`
	Proxy      *flareProxy `json:"proxy,omitempty"`
}

type flareProxy struct {
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type flareSolverrResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Solution struct {
		Status   int    `json:"status"`
		Response string `json:"response"`
	} `json:"solution"`
}

// setFlareProxy attaches proxyURL to reqBody, splitting any embedded credentials
// into separate username/password fields — Chrome's --proxy-server rejects
// credentials inline in the URL (ERR_NO_SUPPORTED_PROXIES).
func setFlareProxy(reqBody *flareSolverrRequest, proxyURL string) error {
	if proxyURL == "" {
		return nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("parse proxy url: %w", err)
	}
	p := &flareProxy{}
	if u.User != nil {
		p.Username = u.User.Username()
		if pw, ok := u.User.Password(); ok {
			p.Password = pw
		}
		u.User = nil
	}
	p.URL = u.String()
	reqBody.Proxy = p
	return nil
}

// SubmitViaFlareSolverr POSTs postData to targetURL through a FlareSolverr
// instance (endpoint, e.g. http://localhost:8191/v1). The request originates
// from FlareSolverr's real browser, which passes Cloudflare's bot check that
// a raw Go client fails. It returns the origin HTTP status and response body
// FlareSolverr observed.
//
// FlareSolverr v3 sends postData with Content-Type
// application/x-www-form-urlencoded and cannot set custom request headers;
// guns.lol reads the raw body as JSON regardless of Content-Type.
func SubmitViaFlareSolverr(endpoint, targetURL, postData, proxyURL string) (int, string, error) {
	reqBody := flareSolverrRequest{
		Cmd:        "request.post",
		URL:        targetURL,
		PostData:   postData,
		MaxTimeout: 60000,
	}
	if err := setFlareProxy(&reqBody, proxyURL); err != nil {
		return 0, "", err
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return 0, "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Reached directly, never through the guns.lol proxy (that proxy is passed
	// to FlareSolverr in the request body).
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	var fsResp flareSolverrResponse
	if err := json.NewDecoder(resp.Body).Decode(&fsResp); err != nil {
		return 0, "", fmt.Errorf("decode flaresolverr response: %w", err)
	}
	if fsResp.Status != "ok" {
		return 0, "", fmt.Errorf("flaresolverr status %q: %s", fsResp.Status, fsResp.Message)
	}
	return fsResp.Solution.Status, fsResp.Solution.Response, nil
}
