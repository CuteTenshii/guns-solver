package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type SolutionPayload struct {
	Username          string
	Underscore2xa     string
	Nonce             string
	O09               string
	Timestamp         int64
	Oo                string
	Seal              string
	TurnstileResponse string
}

// truncate shortens s to at most n characters for readable error messages.
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

// analyticsViewURL is the endpoint that records a profile view. Cloudflare's
// Managed Challenge is attached to this path (not the profile page), so it is
// also the URL FlareSolverr must hit to mint a cf_clearance cookie.
const analyticsViewURL = "https://guns.lol/api/analytics/view"

// gppFormatVersion is the constant "_c14f" value in the analytics payload.
// It is 2 in every observed request and page (RSC) payload; its exact meaning
// (schema/format version) is unconfirmed but it does not vary per request.
const gppFormatVersion = 2

// gppChallenge is the "_gpp_ch" object of the /api/analytics/view payload. The
// server expects these obfuscated, build-stable key names; Go marshals struct
// fields in declaration order, matching the observed request field order.
type gppChallenge struct {
	Version   int    `json:"_c14f"`
	Timestamp int64  `json:"_e92a"`
	O09       string `json:"_b73d"`
	Nonce     string `json:"_a18c"`
	Salt      string `json:"_f52b"` // _2xa (public salt) from the profile page
	Oo        string `json:"_7bc1"` // _oo produced by the WASM solver
}

// recordBody is the full JSON body POSTed to /api/analytics/view.
type recordBody struct {
	Turnstile  string       `json:"_t"`
	GppCh      gppChallenge `json:"_gpp_ch"`
	Username   string       `json:"username"`
	DeviceType string       `json:"deviceType"`
	Referrer   string       `json:"referrer"`
}

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

func SubmitSolution(payload SolutionPayload) error {
	body := recordBody{
		Turnstile: payload.TurnstileResponse,
		GppCh: gppChallenge{
			Version:   gppFormatVersion,
			Timestamp: payload.Timestamp,
			O09:       payload.O09,
			Nonce:     payload.Nonce,
			Salt:      payload.Underscore2xa,
			Oo:        payload.Oo,
		},
		Username:   payload.Username,
		DeviceType: "desktop", //[]string{"desktop", "mobile"}[rand.Intn(2)],
		Referrer:   "",        //https://miwa.lol/tenshii", //"https://guns.lol/" + payload.Username,
	}
	jsonPayload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Cloudflare gates this endpoint by client fingerprint: a raw Go client is
	// challenged (403) while a real browser passes. When a FlareSolverr endpoint
	// is configured, POST the payload through its browser instead of directly.
	if flaresolverrEndpoint != "" {
		status, respText, err := SubmitViaFlareSolverr(flaresolverrEndpoint, analyticsViewURL, string(jsonPayload), proxyURL)
		if err != nil {
			return fmt.Errorf("submit via flaresolverr: %w", err)
		}
		if status != http.StatusOK {
			return fmt.Errorf("submit via flaresolverr returned origin status %d: %s", status, truncate(respText, 300))
		}
		return nil
	}

	req, err := http.NewRequest("POST", analyticsViewURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", "https://guns.lol")
	req.Header.Set("Referer", "https://guns.lol/"+payload.Username)
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
		return err
	}
	defer resp.Body.Close()
	captureCfClearance(resp)

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit solution, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
