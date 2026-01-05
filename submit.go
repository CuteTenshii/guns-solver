package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SolutionPayload struct {
	Username          string
	Underscore2xa     string
	Nonce             string
	O09               string
	Timestamp         int64
	ResultHash        string
	TurnstileResponse string
}

func SubmitSolution(payload SolutionPayload) error {
	sealObject := map[string]string{
		"_oo":  "",
		"seal": "",
	}
	p := map[string]interface{}{
		"_t": payload.TurnstileResponse,
		"_gpp_ch": []interface{}{
			payload.Underscore2xa,
			payload.Timestamp,
			payload.O09,
			payload.Nonce,
			sealObject,
		},
		"username":   payload.Username,
		"deviceType": "desktop",
		"event":      "view",
		"linkId":     nil,
		"referrer":   "",
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	fmt.Println("=== PAYLOAD STRUCTURE ===")
	prettyPayload, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(prettyPayload))
	fmt.Println("=== RAW JSON ===")
	fmt.Println(string(jsonPayload))
	fmt.Println("=========================")
	req, err := http.NewRequest("POST", "https://guns.lol/api/analytics/record", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	if gunsClearance != "" {
		req.Header.Set("Cookie", fmt.Sprintf("gunslol_clearance=%s", gunsClearance))
	}
	//resp, err := httpClient.Do(req)
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	body, _ := io.ReadAll(resp.Body)
	//	return fmt.Errorf("failed to submit solution, status: %d, body: %s", resp.StatusCode, string(body))
	//}
	return nil
}
