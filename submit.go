package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

// SolutionPayload holds everything needed to record a profile view: the scraped
// `_gpp_ch` challenge fields, the WASM proof, and the Turnstile token.
type SolutionPayload struct {
	Username  string
	Version   int    // v
	ID        string // e
	Timestamp int64  // t
	Nonce     string // n
	Seal      string // s
	Challenge string // c
	Proof     string // _oo produced by the WASM solver
	Token     string // Turnstile token
	Referrer  string
	// DeviceNum is the numeric device enum guns.lol now expects:
	// desktop=0, mobile=1, tablet=2.
	DeviceNum int
}

// truncate shortens s to at most n characters for readable error messages.
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

// analyticsViewURL is the endpoint that records a profile view. Cloudflare's
// Managed Challenge is attached to this path (not the profile page), so a raw
// client without a valid cf_clearance cookie is answered with a 403 challenge.
const analyticsViewURL = "https://guns.lol/api/analytics/view"

func SubmitLinkClick(username, linkID string) error {
	p := map[string]interface{}{
		"username":   username,
		"event":      "click",
		"linkId":     linkID,
		"referrer":   "https://guns.lol/" + username,
		"deviceType": []string{"desktop", "mobile", "tablet"}[rand.Intn(3)],
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://guns.lol/api/analytics/record", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	addClearanceCookies(req)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	captureCfClearance(resp)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to submit link click, status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

// buildViewPayload marshals the positional array guns.lol's rewritten view
// beacon expects:
//
//	[token, [v, e, t, n, s, c, proof], username, deviceNum, referrer]
func buildViewPayload(p SolutionPayload) ([]byte, error) {
	challenge := []any{p.Version, p.ID, p.Timestamp, p.Nonce, p.Seal, p.Challenge, p.Proof}
	return json.Marshal([]any{p.Token, challenge, p.Username, p.DeviceNum, p.Referrer})
}

func SubmitSolution(ctx context.Context, p SolutionPayload, capmonsterKey string) error {
	jsonPayload, err := buildViewPayload(p)
	if err != nil {
		return err
	}

	status, respBody, err := postAnalyticsView(jsonPayload, p.Username)
	if err != nil {
		return err
	}

	// Cloudflare gates this endpoint by client fingerprint: a raw Go client is
	// answered with a 403 challenge unless it presents a valid cf_clearance
	// cookie. Mint one with CapMonster from the interstitial we just got, store
	// it, and retry the POST once.
	if status == http.StatusForbidden {
		warnln("Cloudflare 403 on submit — minting cf_clearance via CapMonster")
		cookie, err := SolveCfClearance(ctx, capmonsterKey, "https://guns.lol/"+p.Username, respBody, proxyURL)
		if err != nil {
			return fmt.Errorf("mint cf_clearance: %w", err)
		}
		cfClearance = cookie
		saveCfClearance()

		status, respBody, err = postAnalyticsView(jsonPayload, p.Username)
		if err != nil {
			return err
		}
	}

	if status != http.StatusOK {
		return fmt.Errorf("failed to submit solution, status: %d, body: %s", status, truncate(string(respBody), 300))
	}
	return nil
}

// postAnalyticsView sends the view payload to the analytics endpoint with the
// browser-like headers and clearance cookies, returning the origin status and
// response body. Any cf_clearance the response sets is captured for reuse.
func postAnalyticsView(jsonPayload []byte, username string) (int, []byte, error) {
	req, err := http.NewRequest("POST", analyticsViewURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://guns.lol")
	req.Header.Set("Referer", "https://guns.lol/"+username)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("sec-ch-ua", `"Chromium";v="146", "Google Chrome";v="146", "Not.A/Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")

	req.AddCookie(&http.Cookie{Name: "GUNS_LOCALE", Value: "en"})
	req.AddCookie(&http.Cookie{Name: "GUNS_PATH_LOCALE", Value: "en"})
	addClearanceCookies(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	captureCfClearance(resp)

	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, respBody, nil
}
