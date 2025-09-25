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
	publicSaltRegex = regexp.MustCompile(`\s+ps: '([a-f0-9]{59,64})',`)
	challengeRegex  = regexp.MustCompile(`\s+c: '([a-f0-9]{64})',`)
	difficultyRegex = regexp.MustCompile(`\s+d: '(\d{1,3})',`)
	nonceRegex      = regexp.MustCompile(`\s+_n: '([a-zA-Z0-9]{16})',`)
	timestampRegex  = regexp.MustCompile(`\s+org_ts: '(\d+)'\\n`)
)

type WorkerData struct {
	PublicSalt        string
	Challenge         string
	Difficulty        int
	Nonce             string
	OriginalTimestamp int64
}

func FetchWorkerData(username string) (*WorkerData, error) {
	req, err := http.NewRequest("GET", "https://guns.lol/"+username, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	if gunsClearance != "" {
		req.Header.Set("Cookie", fmt.Sprintf("gunslol_clearance=%s", gunsClearance))
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
			// Extract the gunslol_clearance cookie
			if cookie.Name == "gunslol_clearance" {
				log.Println("Obtained \"gunslol_clearance\" cookie")
				gunsClearance = cookie.Value
				foundClearance = true
				break
			}
		}
		if !foundClearance {
			return nil, err
		} else {
			return FetchWorkerData(username)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	publicSalt := publicSaltRegex.FindStringSubmatch(string(bodyBytes))
	challenge := challengeRegex.FindStringSubmatch(string(bodyBytes))
	difficulty := difficultyRegex.FindStringSubmatch(string(bodyBytes))
	nonce := nonceRegex.FindStringSubmatch(string(bodyBytes))
	timestamp := timestampRegex.FindStringSubmatch(string(bodyBytes))

	if publicSalt == nil || challenge == nil || difficulty == nil || nonce == nil || timestamp == nil {
		return nil, errors.New("failed to parse worker data from response")
	}

	diff, err := strconv.Atoi(difficulty[1])
	if err != nil {
		return nil, err
	}
	originalTs, err := strconv.ParseInt(timestamp[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &WorkerData{
		PublicSalt:        publicSalt[1],
		Challenge:         challenge[1],
		Difficulty:        diff,
		Nonce:             nonce[1],
		OriginalTimestamp: originalTs,
	}, nil
}
