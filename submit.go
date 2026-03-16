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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
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
	p := map[string]interface{}{
		"_t": payload.TurnstileResponse,
		"_gpp_ch": []interface{}{
			payload.Underscore2xa,
			payload.Timestamp,
			payload.O09,
			payload.Nonce,
			payload.Seal,
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
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
		return fmt.Errorf("failed to submit solution, status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}
