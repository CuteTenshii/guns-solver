package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// session holds all per-worker request state. Splitting this out of package
// globals lets parallel view workers each egress through their own proxy
// session (a freshly resolved {session} token → a distinct IP) and hold their
// own cf_clearance / guns_clearance cookies, which is required because
// cf_clearance is bound to the IP that solved the Cloudflare challenge.
//
// A session is used by a single goroutine at a time: within runView the HTTP
// steps are sequential, and the concurrent PoW/Turnstile solve touches neither
// the client nor the cookies. It therefore needs no internal locking.
type session struct {
	client   *http.Client
	proxyURL string // resolved proxy this session egresses through ("" = direct)

	gunsClearance string
	cfClearance   string

	// persist mirrors the clearance cookies to disk (the single-run cost-saving
	// cache). Parallel workers use distinct IPs, so their IP-bound cookies must
	// not be shared through the on-disk cache — they run with persist=false.
	persist bool

	// onLog, when set, routes this session's info/warn lines to a sink (the
	// parallel board) instead of the global terminal helpers, so nested
	// subsystem output does not corrupt the live multi-worker display.
	onLog func(level, msg string)
}

// The guns_clearance and cf_clearance cookies are cached between runs in the
// same temp dir as the fetched PoW modules (powCacheRoot).
func clearanceFile() string   { return filepath.Join(powCacheRoot, "clearance.txt") }
func cfClearanceFile() string { return filepath.Join(powCacheRoot, "cf_clearance.txt") }

// newSession builds a session that egresses through rawProxy (empty for a
// direct connection). A {session} placeholder in rawProxy is resolved to a
// fresh random token here, so every call hands out an independent sticky
// session — and thus a distinct exit IP — from a rotating proxy provider. When
// persist is set the clearance cookies are seeded from and written back to the
// on-disk cache.
func newSession(rawProxy string, persist bool) (*session, error) {
	s := &session{
		persist: persist,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	if rawProxy != "" {
		resolved := resolveProxySession(rawProxy)
		u, err := url.Parse(resolved)
		if err != nil {
			return nil, err
		}
		s.client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
		s.proxyURL = resolved
	}

	if persist {
		if data, err := os.ReadFile(clearanceFile()); err == nil {
			s.gunsClearance = strings.TrimSpace(string(data))
		}
		if data, err := os.ReadFile(cfClearanceFile()); err == nil {
			s.cfClearance = strings.TrimSpace(string(data))
		}
	}
	return s, nil
}

// addClearanceCookies attaches the guns.lol and Cloudflare clearance cookies to
// req when they are known. cf_clearance is required by the Cloudflare Managed
// Challenge in front of guns.lol; without it the API returns 403.
func (s *session) addClearanceCookies(req *http.Request) {
	if s.gunsClearance != "" {
		req.AddCookie(&http.Cookie{Name: "guns_clearance", Value: s.gunsClearance})
	}
	if s.cfClearance != "" {
		req.AddCookie(&http.Cookie{Name: "cf_clearance", Value: s.cfClearance})
	}
}

// captureCfClearance stores a cf_clearance cookie if the response sets one, so a
// rotated value survives across this session's requests (and, when persisting,
// across runs).
func (s *session) captureCfClearance(resp *http.Response) {
	for _, c := range resp.Cookies() {
		if c.Name == "cf_clearance" && c.Value != "" && c.Value != s.cfClearance {
			s.cfClearance = c.Value
			s.saveCfClearance()
		}
	}
}

func (s *session) saveGunsClearance() {
	if !s.persist {
		return
	}
	if err := os.MkdirAll(powCacheRoot, 0700); err == nil {
		os.WriteFile(clearanceFile(), []byte(s.gunsClearance), 0600)
	}
}

func (s *session) saveCfClearance() {
	if !s.persist {
		return
	}
	if err := os.MkdirAll(powCacheRoot, 0700); err == nil {
		os.WriteFile(cfClearanceFile(), []byte(s.cfClearance), 0600)
	}
}

// infof / warnf emit a nested status line. In single-run mode they fall through
// to the global terminal helpers; in parallel mode onLog redirects them to the
// board so they print above the live worker rows without corrupting them.
func (s *session) infof(format string, a ...any) {
	if s.onLog != nil {
		s.onLog("info", fmt.Sprintf(format, a...))
		return
	}
	infoln(format, a...)
}

func (s *session) warnf(format string, a ...any) {
	if s.onLog != nil {
		s.onLog("warn", fmt.Sprintf(format, a...))
		return
	}
	warnln(format, a...)
}
