package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	if data, err := os.ReadFile("clearance.txt"); err == nil {
		gunsClearance = strings.TrimSpace(string(data))
	}
}

func saveClearance() {
	os.WriteFile("clearance.txt", []byte(gunsClearance), 0600)
}

var (
	gunsClearance = ""
	httpClient    = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// Profile page regexes (single-quoted JS values)
	nonceRegex         = regexp.MustCompile(`_n: '([a-zA-Z0-9]{31,32})',`)
	o09Regex           = regexp.MustCompile(`o09: '([a-f0-9]{64})',`)
	underscore2xaRegex = regexp.MustCompile(`_2xa: '([a-zA-Z0-9_-]{80,})',`)
	timestampRegex     = regexp.MustCompile(`org_ts: \\"(\d+)\\",`)
	// Challenge page regexes (double-quoted JS object fields in _gs_sets)
	challengeNonceRegex = regexp.MustCompile(`[{,]_n:"([^"]+)"`)
	challengeO09Regex   = regexp.MustCompile(`[{,]o09:"([^"]+)"`)
	challenge2xaRegex   = regexp.MustCompile(`[{,]_2xa:"([^"]+)"`)
	challengeOrgTsRegex = regexp.MustCompile(`[{,]_org_ts:"([^"]+)"`)
	challengeDRegex     = regexp.MustCompile(`,d:"([^"]+)"`)
	challengeSRegex     = regexp.MustCompile(`[{,]__s:"([^"]+)"`)
)

func SetProxy(rawURL string) error {
	proxyURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return nil
}

// TODO: find better names
type WorkerData struct {
	Nonce             string
	O09               string
	Underscore2xa     string
	OriginalTimestamp int64
}

func FetchWorkerData(ctx context.Context, username string) (*WorkerData, error) {
	req, err := http.NewRequest("GET", "https://guns.lol/"+username, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36")

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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	if resp.StatusCode == http.StatusUnauthorized {
		if gunsClearance != "" {
			gunsClearance = ""
			os.Remove("clearance.txt")
		}
		log.Println("Got 401, solving clearance challenge...")
		if err = solveChallenge(ctx, body); err != nil {
			return nil, fmt.Errorf("challenge: %w", err)
		}
		return FetchWorkerData(ctx, username)
	}

	if resp.StatusCode == http.StatusTemporaryRedirect {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "guns_clearance" {
				gunsClearance = cookie.Value
				saveClearance()
				return FetchWorkerData(ctx, username)
			}
		}
		if gunsClearance != "" {
			location := resp.Header.Get("Location")
			return FetchWorkerData(ctx, location[1:]) // remove leading slash
		}
		return nil, errors.New("307 redirect without clearance cookie")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d %s", resp.StatusCode, resp.Status)
	}

	nonce := nonceRegex.FindStringSubmatch(body)
	o09 := o09Regex.FindStringSubmatch(body)
	underscore2xa := underscore2xaRegex.FindStringSubmatch(body)
	timestamp := timestampRegex.FindStringSubmatch(body)

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

func solveChallenge(ctx context.Context, body string) error {
	nonce := challengeNonceRegex.FindStringSubmatch(body)
	o09 := challengeO09Regex.FindStringSubmatch(body)
	twoXa := challenge2xaRegex.FindStringSubmatch(body)
	orgTs := challengeOrgTsRegex.FindStringSubmatch(body)
	d := challengeDRegex.FindStringSubmatch(body)
	s := challengeSRegex.FindStringSubmatch(body)

	if nonce == nil || o09 == nil || twoXa == nil || orgTs == nil || d == nil || s == nil {
		return errors.New("failed to parse challenge data from page")
	}

	difficulty, err := strconv.Atoi(d[1])
	if err != nil {
		return fmt.Errorf("invalid difficulty %q: %w", d[1], err)
	}

	log.Printf("Challenge params: o09=%s nonce=%s difficulty=%d", o09[1], nonce[1], difficulty)
	res, err := SolveWithWasm(ctx, o09[1], difficulty, orgTs[1], nonce[1], twoXa[1])
	if err != nil {
		return fmt.Errorf("wasm solve: %w", err)
	}
	log.Printf("Challenge solved: _oo=%s", res.Oo)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("_o", res.Oo)
	w.WriteField("_s", s[1])
	w.WriteField("_u", nonce[1])
	w.WriteField("_i", twoXa[1])
	w.WriteField("_x", o09[1])
	w.WriteField("_t", orgTs[1])
	w.Close()

	req, err := http.NewRequest("POST", "https://guns.lol/_challenge/verify", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36")

	vresp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer vresp.Body.Close()
	io.ReadAll(vresp.Body) // drain

	if vresp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify returned %d", vresp.StatusCode)
	}

	for _, cookie := range vresp.Cookies() {
		if cookie.Name == "guns_clearance" {
			gunsClearance = cookie.Value
			saveClearance()
			return nil
		}
	}
	return errors.New("no clearance cookie in verify response")
}
