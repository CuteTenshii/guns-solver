package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// turnstileSiteKey is the Cloudflare Turnstile sitekey guns.lol renders the
// view widget with.
const turnstileSiteKey = "0x4AAAAAAAgU7T2niLQD-TLm"

// taskPollInterval is the delay between task-result polls. Both services
// document a 2-3s cadence and throttle callers that poll tighter.
const taskPollInterval = 3 * time.Second

// The -captcha-provider values this tool accepts.
const (
	providerCapmonster = "capmonster"
	providerSolverify  = "solverify"
	providerUncaptcha  = "uncaptcha"
)

// supportedCaptchaProviders is the accepted -captcha-provider set, in the order
// it is shown to the user.
var supportedCaptchaProviders = []string{providerCapmonster, providerSolverify, providerUncaptcha}

// turnstileRequest describes the Turnstile widget guarding a profile view.
type turnstileRequest struct {
	PageURL string
	SiteKey string
	// Action and CData are the widget's data-action / data-cdata values, taken
	// from the challenge object on the profile page.
	Action string
	CData  string
}

// clearanceRequest describes a Cloudflare Managed Challenge to clear.
type clearanceRequest struct {
	// PageURL is the site the cookie is for.
	PageURL string
	// ChallengeURL is the exact URL whose response carried the challenge.
	// Providers that browse to the challenge themselves must navigate here
	// rather than to PageURL: Cloudflare only issues cf_clearance when it
	// actually challenges, and the profile page is served unchallenged to a
	// real browser. The cookie is scoped to the whole guns.lol zone, so
	// clearing it on any path covers the rest.
	ChallengeURL string
	// ChallengeMethod is the HTTP method that drew the challenge. The analytics
	// endpoint only answers POST, so a GET navigation reaches a different rule
	// and is never challenged.
	ChallengeMethod string
	// HTMLPage is the "Just a moment" interstitial body captured from the 403.
	// CapMonster replays it; Solverify navigates the challenge itself and
	// ignores it.
	HTMLPage []byte
	// ProxyURL must be the proxy the tool itself egresses through: cf_clearance
	// is bound to the IP that solved the challenge.
	ProxyURL string
	// UserAgent is the browser identity the cookie should be bound to.
	UserAgent string
}

// clearanceResult is a solved Cloudflare Managed Challenge.
type clearanceResult struct {
	// Cookie is the cf_clearance value.
	Cookie string
	// UserAgent is the identity the cookie is bound to, when the provider
	// reports one. CapMonster uses the user-agent it is handed and reports
	// nothing back; Solverify echoes the one its browser used.
	UserAgent string
}

// captchaProvider is a third-party captcha-solving service. Implementations
// hold only an API key and keep no per-request state, so one instance serves
// every parallel view worker.
type captchaProvider interface {
	// Name is the provider's -captcha-provider value, used in messages.
	Name() string
	// SolveTurnstile returns a Turnstile response token for the given widget.
	SolveTurnstile(ctx context.Context, req turnstileRequest) (string, error)
	// SolveCfClearance returns a cf_clearance cookie for req.PageURL.
	SolveCfClearance(ctx context.Context, req clearanceRequest) (clearanceResult, error)
}

// newCaptchaProvider builds the provider selected by name.
func newCaptchaProvider(name, apiKey string) (captchaProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("no API key given for captcha provider %q", name)
	}
	switch name {
	case providerCapmonster:
		return &capmonsterProvider{apiKey: apiKey}, nil
	case providerSolverify:
		return &solverifyProvider{apiKey: apiKey}, nil
	case providerUncaptcha:
		return &uncaptchaProvider{apiKey: apiKey}, nil
	default:
		return nil, fmt.Errorf("unknown captcha provider %q (supported: %s)",
			name, strings.Join(supportedCaptchaProviders, ", "))
	}
}

// captchaFlags is the raw provider selection as it arrived on the command line.
type captchaFlags struct {
	provider string
	key      string
}

// resolve settles the provider name and API key from the flags, applying the
// default provider when none was named. An unknown name is not rejected here —
// newCaptchaProvider reports that, alongside the supported set.
func (f captchaFlags) resolve() (name, key string) {
	name, key = f.provider, f.key
	if name == "" {
		name = providerCapmonster
	}
	return name, key
}

// proxyParts is a proxy URL split into the discrete fields captcha services
// expect, since none of them accept inline credentials in a single URL.
type proxyParts struct {
	Scheme   string
	Host     string
	Port     int
	Login    string
	Password string
}

// parseProxyParts splits a proxy URL (e.g. http://user:pass@host:port) into its
// components.
func parseProxyParts(rawURL string) (proxyParts, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return proxyParts{}, fmt.Errorf("parse proxy url: %w", err)
	}

	p := proxyParts{Scheme: u.Scheme, Host: u.Hostname()}
	if p.Scheme == "" || p.Host == "" {
		return proxyParts{}, fmt.Errorf("proxy url %q must include scheme and host", rawURL)
	}

	port := u.Port()
	if port == "" {
		return proxyParts{}, fmt.Errorf("proxy url %q must include a port", rawURL)
	}
	p.Port, err = strconv.Atoi(port)
	if err != nil {
		return proxyParts{}, fmt.Errorf("invalid proxy port %q: %w", port, err)
	}

	if u.User != nil {
		p.Login = u.User.Username()
		if pw, ok := u.User.Password(); ok {
			p.Password = pw
		}
	}
	return p, nil
}

// postJSON posts body as JSON to endpoint and decodes the reply into out. It
// is the form used by services that authenticate through the request body and
// report failures inside it, so the HTTP status carries no information.
func postJSON(ctx context.Context, endpoint string, body, out any) error {
	_, err := postJSONWithHeaders(ctx, endpoint, nil, body, out)
	return err
}

// postJSONWithHeaders posts body as JSON to endpoint with additional request
// headers, decodes the reply into out, and returns the HTTP status alongside
// any error. Services that authenticate through a header and signal failure
// through the status line need both; the status is also returned with an error
// so a caller can tell a transient failure from a permanent one.
func postJSONWithHeaders(ctx context.Context, endpoint string, headers map[string]string, body, out any) (int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read %s response: %w", endpoint, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return resp.StatusCode, fmt.Errorf("decode %s response (status %d): %w: %s",
			endpoint, resp.StatusCode, err, truncate(string(data), 200))
	}
	return resp.StatusCode, nil
}

// pollTask calls fetch every taskPollInterval until it reports the task
// finished, it fails, or ctx expires. The first poll is delayed too: no service
// resolves a task instantly.
func pollTask(ctx context.Context, fetch func() (done bool, err error)) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(taskPollInterval):
		}

		done, err := fetch()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}
