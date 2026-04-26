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

const hexChars = "0123456789abcdef"

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
	if gunsClearance != "" {
		req.AddCookie(&http.Cookie{Name: "guns_clearance", Value: gunsClearance})
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to submit link click, status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func SubmitSolution(payload SolutionPayload) error {
	insertPos := int(payload.Timestamp % 10)
	seal := payload.Seal[:insertPos] + string(hexChars[rand.Intn(16)]) + payload.Seal[insertPos:]

	p := map[string]interface{}{
		"_t": payload.TurnstileResponse,
		"_gpp_ch": []interface{}{
			payload.Underscore2xa,
			payload.Timestamp,
			payload.O09,
			payload.Nonce,
			seal,
			payload.Oo,
		},
		"username":   payload.Username,
		"deviceType": []string{"desktop", "mobile"}[rand.Intn(2)],
		"event":      "view",
		"linkId":     nil,
		"referrer":   "https://miwa.lol/tenshii",
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://guns.lol/api/analytics/record", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
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
	if gunsClearance != "" {
		req.AddCookie(&http.Cookie{Name: "guns_clearance", Value: gunsClearance})
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit solution, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
