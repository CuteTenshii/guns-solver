package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	turnstileSiteKey = "0x4AAAAAAAgU7T2niLQD-TLm"
	capmonsterURL    = "https://api.capmonster.cloud"
)

type capmonsterCreateRequest struct {
	ClientKey string         `json:"clientKey"`
	Task      capmonsterTask `json:"task"`
}

type capmonsterTask struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	WebsiteKey    string `json:"websiteKey"`
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

type capmonsterCreateResponse struct {
	ErrorID int `json:"errorId"`
	TaskID  int `json:"taskId"`
}

type capmonsterResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

type capmonsterResultResponse struct {
	ErrorID  int    `json:"errorId"`
	Status   string `json:"status"`
	Solution struct {
		Token string `json:"token"`
	} `json:"solution"`
}

func SolveTurnstile(ctx context.Context, apiKey, pageURL, proxyURL string) (string, error) {
	task := capmonsterTask{
		Type:       "TurnstileTaskProxyless",
		WebsiteURL: pageURL,
		WebsiteKey: turnstileSiteKey,
	}

	if proxyURL == "" {
		pURL, err := url.Parse(proxyURL)
		if err != nil {
			return "", fmt.Errorf("invalid proxy URL: %w", err)
		}
		port, _ := strconv.Atoi(pURL.Port())
		task.Type = "TurnstileTask"
		task.ProxyType = pURL.Scheme
		task.ProxyAddress = pURL.Hostname()
		task.ProxyPort = port
		if pURL.User != nil {
			task.ProxyLogin = pURL.User.Username()
			task.ProxyPassword, _ = pURL.User.Password()
		}
	}

	createBody, _ := json.Marshal(capmonsterCreateRequest{ClientKey: apiKey, Task: task})

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

		req, err := http.NewRequestWithContext(ctx, "POST", capmonsterURL+"/getTaskResult", bytes.NewReader(resultBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		var result capmonsterResultResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.ErrorID != 0 {
			return "", fmt.Errorf("capmonster getTaskResult error %d", result.ErrorID)
		}
		if result.Status == "ready" {
			return result.Solution.Token, nil
		}
	}
}
