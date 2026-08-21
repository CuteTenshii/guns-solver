package main

import (
	"encoding/json"
	"testing"
)

func TestCaptchaFlagsResolve(t *testing.T) {
	tests := []struct {
		name         string
		flags        captchaFlags
		wantProvider string
		wantKey      string
	}{
		{
			name:         "defaults to capmonster",
			flags:        captchaFlags{key: "k"},
			wantProvider: providerCapmonster,
			wantKey:      "k",
		},
		{
			name:         "explicit solverify",
			flags:        captchaFlags{provider: providerSolverify, key: "k"},
			wantProvider: providerSolverify,
			wantKey:      "k",
		},
		{
			name:         "explicit uncaptcha",
			flags:        captchaFlags{provider: providerUncaptcha, key: "k"},
			wantProvider: providerUncaptcha,
			wantKey:      "k",
		},
		{
			name:         "no key resolves but stays empty",
			flags:        captchaFlags{},
			wantProvider: providerCapmonster,
		},
		{
			// An unrecognised name is passed through for newCaptchaProvider to
			// reject, so the error names the supported set in one place.
			name:         "unknown provider passes through",
			flags:        captchaFlags{provider: "nope", key: "k"},
			wantProvider: "nope",
			wantKey:      "k",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			provider, key := tc.flags.resolve()
			if provider != tc.wantProvider {
				t.Errorf("provider = %q, want %q", provider, tc.wantProvider)
			}
			if key != tc.wantKey {
				t.Errorf("key = %q, want %q", key, tc.wantKey)
			}
		})
	}
}

func TestNewCaptchaProvider(t *testing.T) {
	for _, name := range supportedCaptchaProviders {
		p, err := newCaptchaProvider(name, "key")
		if err != nil {
			t.Fatalf("newCaptchaProvider(%q): %v", name, err)
		}
		if p.Name() != name {
			t.Errorf("provider %q reports name %q", name, p.Name())
		}
	}
	if _, err := newCaptchaProvider("nope", "key"); err == nil {
		t.Error("expected an error for an unknown provider")
	}
	if _, err := newCaptchaProvider(providerSolverify, ""); err == nil {
		t.Error("expected an error for a missing API key")
	}
}

func TestParseProxyParts(t *testing.T) {
	got, err := parseProxyParts("http://user:pa:ss@proxy.example.com:8080")
	if err != nil {
		t.Fatalf("parseProxyParts: %v", err)
	}
	want := proxyParts{Scheme: "http", Host: "proxy.example.com", Port: 8080, Login: "user", Password: "pa:ss"}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

	if _, err := parseProxyParts("proxy.example.com:8080"); err == nil {
		t.Error("expected an error for a proxy without a scheme")
	}
	if _, err := parseProxyParts("http://proxy.example.com"); err == nil {
		t.Error("expected an error for a proxy without a port")
	}
}

// TestSolverifyTaskShape pins the JSON field names Solverify's API expects;
// they differ from CapMonster's (cdata/action, lowercase useragent, and a
// string proxyPort), so a rename here silently breaks solving.
func TestSolverifyTaskShape(t *testing.T) {
	turnstile, err := json.Marshal(solverifyTurnstileTask{
		Type:       "turnstile",
		WebsiteURL: "https://guns.lol/user",
		WebsiteKey: turnstileSiteKey,
		Action:     "guns_view",
		CData:      "cd",
	})
	if err != nil {
		t.Fatalf("marshal turnstile task: %v", err)
	}
	wantTurnstile := `{"type":"turnstile","websiteURL":"https://guns.lol/user","websiteKey":"` +
		turnstileSiteKey + `","action":"guns_view","cdata":"cd"}`
	if string(turnstile) != wantTurnstile {
		t.Errorf("turnstile task =\n%s\nwant\n%s", turnstile, wantTurnstile)
	}

	// action and cdata must be omitted rather than sent empty.
	bare, err := json.Marshal(solverifyTurnstileTask{Type: "turnstile", WebsiteURL: "u", WebsiteKey: "k"})
	if err != nil {
		t.Fatalf("marshal bare turnstile task: %v", err)
	}
	if string(bare) != `{"type":"turnstile","websiteURL":"u","websiteKey":"k"}` {
		t.Errorf("bare turnstile task = %s", bare)
	}

	interstitial, err := json.Marshal(solverifyInterstitialTask{
		Type:          "interstitial",
		WebsiteURL:    "https://guns.lol/api/analytics/view",
		UserAgent:     "UA",
		Source:        true,
		Method:        "POST",
		ProxyType:     "http",
		ProxyAddress:  "1.1.1.1",
		ProxyPort:     "8080",
		ProxyLogin:    "u",
		ProxyPassword: "p",
	})
	if err != nil {
		t.Fatalf("marshal interstitial task: %v", err)
	}
	wantInterstitial := `{"type":"interstitial","websiteURL":"https://guns.lol/api/analytics/view","useragent":"UA",` +
		`"source":true,"method":"POST","proxyType":"http","proxyAddress":"1.1.1.1","proxyPort":"8080",` +
		`"proxyLogin":"u","proxyPassword":"p"}`
	if string(interstitial) != wantInterstitial {
		t.Errorf("interstitial task =\n%s\nwant\n%s", interstitial, wantInterstitial)
	}
}

// TestSolverifyResultDecode decodes the documented completed-task responses and
// checks the fields the solver reads out of them.
func TestSolverifyResultDecode(t *testing.T) {
	const turnstileBody = `{"errorId":0,"status":"completed","solution":{"value":"0.MpR_tg"},` +
		`"errorCode":null,"errorDescription":null}`
	var turnstile solverifyResultResponse
	if err := json.Unmarshal([]byte(turnstileBody), &turnstile); err != nil {
		t.Fatalf("decode turnstile result: %v", err)
	}
	if turnstile.Status != solverifyCompleted {
		t.Errorf("status = %q, want %q", turnstile.Status, solverifyCompleted)
	}
	if turnstile.Solution.Value != "0.MpR_tg" {
		t.Errorf("token = %q", turnstile.Solution.Value)
	}

	const interstitialBody = `{"errorId":0,"status":"completed","solution":{"value":"dce67487",` +
		`"cookies":{"cf_clearance":".AdOHyy60coo8q0wL6TAJFu2V1"},"useragent":"UA","html":"<html></html>"},` +
		`"errorCode":null,"errorDescription":null}`
	var interstitial solverifyResultResponse
	if err := json.Unmarshal([]byte(interstitialBody), &interstitial); err != nil {
		t.Fatalf("decode interstitial result: %v", err)
	}
	if got := interstitial.Solution.Cookies["cf_clearance"]; got != ".AdOHyy60coo8q0wL6TAJFu2V1" {
		t.Errorf("cf_clearance = %q", got)
	}
	if interstitial.Solution.UserAgent != "UA" {
		t.Errorf("useragent = %q", interstitial.Solution.UserAgent)
	}

	// The jar is also accepted as a list of cookie objects.
	const listBody = `{"errorId":0,"status":"completed","solution":{"cookies":` +
		`[{"name":"__cf_bm","value":"a"},{"name":"cf_clearance","value":"b"}],"useragent":"UA"}}`
	var asList solverifyResultResponse
	if err := json.Unmarshal([]byte(listBody), &asList); err != nil {
		t.Fatalf("decode list-shaped cookies: %v", err)
	}
	if got := asList.Solution.Cookies["cf_clearance"]; got != "b" {
		t.Errorf("cf_clearance from list = %q, want %q", got, "b")
	}
	if got := asList.Solution.Cookies.names(); got != "__cf_bm, cf_clearance" {
		t.Errorf("cookie names = %q", got)
	}

	var empty solverifyCookies
	if got := empty.names(); got != "none" {
		t.Errorf("empty cookie names = %q, want %q", got, "none")
	}

	const failedBody = `{"errorId":1,"status":"failed","errorCode":"ERROR_TASK_FAILED",` +
		`"errorDescription":"solver could not complete the task"}`
	var failed solverifyResultResponse
	if err := json.Unmarshal([]byte(failedBody), &failed); err != nil {
		t.Fatalf("decode failed result: %v", err)
	}
	if failed.ErrorID == 0 || failed.Status != solverifyFailed {
		t.Errorf("failed result decoded as errorId %d status %q", failed.ErrorID, failed.Status)
	}
}

// TestCapmonsterTaskShape pins CapMonster's task field names, which the
// provider split must not have altered.
func TestCapmonsterTaskShape(t *testing.T) {
	turnstile, err := json.Marshal(capmonsterTurnstileTask{
		Type:       "TurnstileTaskProxyless",
		WebsiteURL: "https://guns.lol/user",
		WebsiteKey: turnstileSiteKey,
		PageAction: "guns_view",
		Data:       "cd",
	})
	if err != nil {
		t.Fatalf("marshal turnstile task: %v", err)
	}
	want := `{"type":"TurnstileTaskProxyless","websiteURL":"https://guns.lol/user","websiteKey":"` +
		turnstileSiteKey + `","pageAction":"guns_view","data":"cd"}`
	if string(turnstile) != want {
		t.Errorf("turnstile task =\n%s\nwant\n%s", turnstile, want)
	}

	clearance, err := json.Marshal(capmonsterClearanceTask{
		Type:               "TurnstileTask",
		WebsiteURL:         "https://guns.lol/user",
		WebsiteKey:         "xxxxxx",
		CloudflareTaskType: "cf_clearance",
		HTMLPageBase64:     "aHRtbA==",
		UserAgent:          "UA",
		ProxyType:          "http",
		ProxyAddress:       "1.1.1.1",
		ProxyPort:          8080,
	})
	if err != nil {
		t.Fatalf("marshal clearance task: %v", err)
	}
	wantClearance := `{"type":"TurnstileTask","websiteURL":"https://guns.lol/user","websiteKey":"xxxxxx",` +
		`"cloudflareTaskType":"cf_clearance","htmlPageBase64":"aHRtbA==","userAgent":"UA",` +
		`"proxyType":"http","proxyAddress":"1.1.1.1","proxyPort":8080}`
	if string(clearance) != wantClearance {
		t.Errorf("clearance task =\n%s\nwant\n%s", clearance, wantClearance)
	}
}

// TestSolveCfClearanceRequiresProxy covers the guard both providers share:
// cf_clearance is IP-bound, so it cannot be minted without a proxy.
func TestSolveCfClearanceRequiresProxy(t *testing.T) {
	for _, name := range supportedCaptchaProviders {
		p, err := newCaptchaProvider(name, "key")
		if err != nil {
			t.Fatalf("newCaptchaProvider(%q): %v", name, err)
		}
		if _, err := p.SolveCfClearance(t.Context(), clearanceRequest{PageURL: "https://guns.lol/user"}); err == nil {
			t.Errorf("%s: expected an error without a proxy", name)
		}
	}
}

func TestDescribeLandedPage(t *testing.T) {
	tests := []struct {
		name string
		page string
		want string
	}{
		{name: "empty", page: "", want: "no page captured"},
		{name: "challenge", page: "<html><title>Just a moment...</title></html>", want: "a Cloudflare challenge page"},
		{name: "challenge script", page: `<script src="/cdn-cgi/challenge-platform/h/b/x"></script>`, want: "a Cloudflare challenge page"},
		{name: "origin error", page: "<html>\n  Method   Not  Allowed\n</html>", want: "<html> Method Not Allowed </html>"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := describeLandedPage(tc.page); got != tc.want {
				t.Errorf("describeLandedPage(%q) = %q, want %q", tc.page, got, tc.want)
			}
		})
	}
}

// TestSolverifyRejectsUnsupportedMethod covers the GET/POST-only navigation
// rule, since the challenged request's method is replayed verbatim.
func TestSolverifyRejectsUnsupportedMethod(t *testing.T) {
	p, err := newCaptchaProvider(providerSolverify, "key")
	if err != nil {
		t.Fatalf("newCaptchaProvider: %v", err)
	}
	_, err = p.SolveCfClearance(t.Context(), clearanceRequest{
		PageURL:         "https://guns.lol/user",
		ChallengeURL:    "https://guns.lol/api/analytics/view",
		ChallengeMethod: "PUT",
		ProxyURL:        "http://1.1.1.1:8080",
	})
	if err == nil {
		t.Error("expected an error for a PUT challenge method")
	}
}

// TestSolverifyRejectsNonHTTPProxy covers Solverify's http-only proxy rule.
func TestSolverifyRejectsNonHTTPProxy(t *testing.T) {
	p, err := newCaptchaProvider(providerSolverify, "key")
	if err != nil {
		t.Fatalf("newCaptchaProvider: %v", err)
	}
	_, err = p.SolveCfClearance(t.Context(), clearanceRequest{
		PageURL:  "https://guns.lol/user",
		ProxyURL: "socks5://1.1.1.1:1080",
	})
	if err == nil {
		t.Error("expected an error for a socks5 proxy")
	}
}

// TestUncaptchaTaskShape pins uncaptcha.io's request shape. Its fields are
// snake_case and nested under task_data — nothing like the other two services —
// so a rename here silently breaks solving.
func TestUncaptchaTaskShape(t *testing.T) {
	turnstile, err := json.Marshal(uncaptchaExecuteRequest{
		TaskType: uncaptchaTaskTurnstile,
		TaskData: uncaptchaTurnstileTask{
			URL:     "https://guns.lol/user",
			SiteKey: turnstileSiteKey,
			Action:  "guns_view",
			CData:   "cd",
		},
	})
	if err != nil {
		t.Fatalf("marshal turnstile task: %v", err)
	}
	wantTurnstile := `{"task_type":"turnstile","task_data":{"url":"https://guns.lol/user","sitekey":"` +
		turnstileSiteKey + `","action":"guns_view","cdata":"cd"}}`
	if string(turnstile) != wantTurnstile {
		t.Errorf("turnstile task =\n%s\nwant\n%s", turnstile, wantTurnstile)
	}

	// action and cdata must be omitted rather than sent empty: the service
	// rejects a solve whose action/cdata do not match the widget.
	bare, err := json.Marshal(uncaptchaTurnstileTask{URL: "u", SiteKey: "k"})
	if err != nil {
		t.Fatalf("marshal bare turnstile task: %v", err)
	}
	if string(bare) != `{"url":"u","sitekey":"k"}` {
		t.Errorf("bare turnstile task = %s", bare)
	}

	waf, err := json.Marshal(uncaptchaExecuteRequest{
		TaskType: uncaptchaTaskWAF,
		TaskData: uncaptchaWAFTask{
			URL:       analyticsViewURL,
			Proxy:     "http://u:p@1.1.1.1:8080",
			UserAgent: "UA",
			HTML:      "<html></html>",
		},
	})
	if err != nil {
		t.Fatalf("marshal waf task: %v", err)
	}
	// encoding/json escapes <, > and & by default, so the challenge HTML goes
	// out as \u003c... — a different encoding of the same string, which the
	// service decodes back to the original. CapMonster never hits this because
	// it takes the interstitial base64-encoded.
	wantWAF := `{"task_type":"waf","task_data":{"url":"` + analyticsViewURL +
		`","proxy":"http://u:p@1.1.1.1:8080","user_agent":"UA","html":"\u003chtml\u003e\u003c/html\u003e"}}`
	if string(waf) != wantWAF {
		t.Errorf("waf task =\n%s\nwant\n%s", waf, wantWAF)
	}

	var decoded struct {
		TaskData struct {
			HTML string `json:"html"`
		} `json:"task_data"`
	}
	if err := json.Unmarshal(waf, &decoded); err != nil {
		t.Fatalf("decode waf task: %v", err)
	}
	if decoded.TaskData.HTML != "<html></html>" {
		t.Errorf("challenge html round-tripped as %q", decoded.TaskData.HTML)
	}
}

// TestUncaptchaResponseDecode decodes the documented success and failure
// envelopes and checks the fields the provider reads out of them.
func TestUncaptchaResponseDecode(t *testing.T) {
	const turnstileBody = `{"success":true,"data":{"solution":{"token":"0.AuR5g8kP"}}}`
	var turnstile uncaptchaResponse
	if err := json.Unmarshal([]byte(turnstileBody), &turnstile); err != nil {
		t.Fatalf("decode turnstile response: %v", err)
	}
	if !turnstile.Success {
		t.Error("turnstile response decoded as unsuccessful")
	}
	if turnstile.Data.Solution.Token != "0.AuR5g8kP" {
		t.Errorf("token = %q", turnstile.Data.Solution.Token)
	}

	const wafBody = `{"success":true,"data":{"solution":{"clearance":"abc123","cf_bm":"xyz789",` +
		`"cf_rt":"token123","headers":{"user-agent":"Mozilla/5.0 UA","sec-ch-ua-mobile":"?0"},` +
		`"attributes":{"a":"value1"}}}}`
	var waf uncaptchaResponse
	if err := json.Unmarshal([]byte(wafBody), &waf); err != nil {
		t.Fatalf("decode waf response: %v", err)
	}
	if waf.Data.Solution.Clearance != "abc123" {
		t.Errorf("clearance = %q", waf.Data.Solution.Clearance)
	}
	if got := solvedUserAgent(waf.Data.Solution.Headers); got != "Mozilla/5.0 UA" {
		t.Errorf("solved user-agent = %q", got)
	}

	const failureBody = `{"success":false,"message":"task timed out"}`
	var failure uncaptchaResponse
	if err := json.Unmarshal([]byte(failureBody), &failure); err != nil {
		t.Fatalf("decode failure response: %v", err)
	}
	if failure.Success || failure.Message != "task timed out" {
		t.Errorf("failure decoded as success %v message %q", failure.Success, failure.Message)
	}
}

func TestSolvedUserAgent(t *testing.T) {
	if got := solvedUserAgent(map[string]string{"User-Agent": "UA"}); got != "UA" {
		t.Errorf("capitalised header lookup = %q, want %q", got, "UA")
	}
	if got := solvedUserAgent(map[string]string{"sec-ch-ua-mobile": "?0"}); got != "" {
		t.Errorf("missing user-agent = %q, want empty", got)
	}
	if got := solvedUserAgent(nil); got != "" {
		t.Errorf("nil headers = %q, want empty", got)
	}
}

// TestUncaptchaRetryable covers which failures are re-sent. uncaptcha signals
// failure through the HTTP status, and retrying a permanent rejection (bad key,
// no balance, blacklisted URL) would only burn the caller's time budget.
func TestUncaptchaRetryable(t *testing.T) {
	retryable := []int{429, 500, 502, 503, 504}
	for _, status := range retryable {
		if !uncaptchaRetryable(status) {
			t.Errorf("status %d should be retryable", status)
		}
	}
	permanent := []int{200, 400, 401, 402, 403, 404}
	for _, status := range permanent {
		if uncaptchaRetryable(status) {
			t.Errorf("status %d should not be retryable", status)
		}
	}
}
