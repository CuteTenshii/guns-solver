package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	gunsClearance = ""
	httpClient    = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	nonceRegex         = regexp.MustCompile(`_n: '([a-zA-Z0-9]{31,32})',`)
	o09Regex           = regexp.MustCompile(`o09: '([a-f0-9]{64})',`)
	underscore2xaRegex = regexp.MustCompile(`_2xa: '([a-zA-Z0-9_-]{118})',`)
	timestampRegex     = regexp.MustCompile(`org_ts: \\"(\d+)\\",`)
)

// TODO: find better names
type WorkerData struct {
	Nonce             string
	O09               string
	Underscore2xa     string
	OriginalTimestamp int64
}

func FetchWorkerData(username string) (*WorkerData, error) {
	req, err := http.NewRequest("GET", "https://guns.lol/"+username, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.AddCookie(&http.Cookie{Name: "GUNS_LOCALE", Value: "en"})
	req.AddCookie(&http.Cookie{Name: "GUNS_PATH_LOCALE", Value: "en"})
	if gunsClearance != "" {
		req.AddCookie(&http.Cookie{Name: "guns_clearance", Value: gunsClearance})
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTemporaryRedirect {
		if gunsClearance != "" {
			// Already have a clearance cookie but still got redirected, that means the username is an alias
			location := resp.Header.Get("Location")
			return FetchWorkerData(location[1:]) // Remove leading slash
		}
		foundClearance := false
		for _, cookie := range resp.Cookies() {
			// Extract the guns_clearance cookie
			if cookie.Name == "guns_clearance" {
				log.Println("Obtained \"guns_clearance\" cookie")
				gunsClearance = cookie.Value
				foundClearance = true
				break
			}
		}
		if !foundClearance {
			return nil, err
		}

		return FetchWorkerData(username)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	nonce := nonceRegex.FindStringSubmatch(string(bodyBytes))
	o09 := o09Regex.FindStringSubmatch(string(bodyBytes))
	underscore2xa := underscore2xaRegex.FindStringSubmatch(string(bodyBytes))
	timestamp := timestampRegex.FindStringSubmatch(string(bodyBytes))

	if nonce == nil || o09 == nil || underscore2xa == nil || timestamp == nil {
		return nil, errors.New("failed to parse worker data from response")
	}

	originalTs, err := strconv.ParseInt(timestamp[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &WorkerData{
		Nonce:             nonce[1],
		O09:               o09[1],
		Underscore2xa:     underscore2xa[1],
		OriginalTimestamp: originalTs,
	}, nil
}
