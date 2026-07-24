package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// userAgent is the User-Agent sent on every guns.lol request, and the UA passed
// to CapMonster when minting a cf_clearance cookie. cf_clearance is bound to the
// UA that solved the Cloudflare challenge, so the same value must be used
// everywhere for the whole flow to stay consistent.
var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"

// perViewTimeout bounds a single record-a-view attempt (challenge fetch through
// submission, including the CapMonster round-trips).
const perViewTimeout = 180 * time.Second

// defaultConcurrency caps how many views run at once when -concurrency is not
// given, so a large -count does not open an unbounded number of proxy sessions.
const defaultConcurrency = 5

func fatalf(format string, a ...any) {
	if activeSpinner != nil {
		activeSpinner.fail()
	}
	fmt.Fprintf(os.Stderr, "\n  %s %s\n\n", red(bold("✗")), fmt.Sprintf(format, a...))
	os.Exit(1)
}

func main() {
	username := flag.String("username", "", "Profile username")
	capmonsterKey := flag.String("capmonster-key", "", "CapMonster API key for solving Turnstile and minting cf_clearance")
	proxy := flag.String("proxy", "", "Proxy URL for guns.lol requests (e.g. http://user:pass@host:port). A {session} placeholder is replaced with a fresh token per worker for rotating sticky-session proxies")
	linkID := flag.String("link-id", "", "Link UUID to record a click event instead of a profile view")
	count := flag.Int("count", 1, "Number of profile views to record")
	concurrency := flag.Int("concurrency", 0, "Maximum views to record simultaneously (default: min(count, 5))")
	flag.Parse()

	banner()

	if *linkID != "" {
		runLinkClick(*username, *linkID, *proxy)
		return
	}

	if *username == "" {
		fatalf("Missing required flag %s (example: guns-solver -username <username>)", bold("-username"))
	}
	if *capmonsterKey == "" {
		fatalf("Missing required flag %s (needed to solve Cloudflare Turnstile)", bold("-capmonster-key"))
	}
	if *count < 1 {
		fatalf("Invalid -count %d: must be at least 1", *count)
	}

	conc := *concurrency
	if conc < 1 {
		conc = defaultConcurrency
	}
	if conc > *count {
		conc = *count
	}

	cfg := viewConfig{Username: *username, CapmonsterKey: *capmonsterKey, DeviceNum: 0, Referrer: ""}

	if *count == 1 {
		runSingle(cfg, *proxy)
		return
	}
	runParallel(cfg, *proxy, *count, conc)
}

// runLinkClick records a single link-click event.
func runLinkClick(username, linkID, proxy string) {
	if username == "" {
		fatalf("Missing required flag %s (needed to record a link click)", bold("-username"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), perViewTimeout)
	defer cancel()

	sess, err := newSession(proxy, true)
	if err != nil {
		fatalf("Invalid proxy URL: %s", err)
	}

	sp := startSpinner("Acquiring Cloudflare clearance")
	if _, err := sess.FetchWorkerData(ctx, username); err != nil {
		fatalf("Error acquiring Cloudflare clearance: %s", err)
	}
	sp.succeed("Acquired Cloudflare clearance")

	sp = startSpinner("Submitting link click")
	if err := sess.SubmitLinkClick(username, linkID); err != nil {
		fatalf("Error submitting link click: %s", err)
	}
	sp.succeed("Submitted link click")

	doneln("Link click recorded for %s", cyan(linkID))
}

// runSingle records one view with the animated per-step display.
func runSingle(cfg viewConfig, proxy string) {
	ctx, cancel := context.WithTimeout(context.Background(), perViewTimeout)
	defer cancel()

	sess, err := newSession(proxy, true)
	if err != nil {
		fatalf("Invalid proxy URL: %s", err)
	}

	infoln("target %s", cyan("guns.lol/"+cfg.Username))
	if err := sess.runView(ctx, cfg, &spinnerReporter{}); err != nil {
		fatalf("%s", err)
	}
	doneln("View recorded for %s", cyan("guns.lol/"+cfg.Username))
}

// runParallel records count views across conc worker slots, each on its own
// proxy session (and thus its own IP and cf_clearance), rendering a live board.
func runParallel(cfg viewConfig, proxy string, count, conc int) {
	if proxy == "" || !strings.Contains(proxy, proxySessionPlaceholder) {
		warnln("workers will share one egress IP; use a {session} proxy so each view gets a distinct IP")
	}

	// Ctrl-C cancels in-flight work; workers stop claiming new views and the
	// board reports what completed.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	b := startBoard(conc, count)

	var next int32
	var wg sync.WaitGroup
	for slot := 0; slot < conc; slot++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				item := int(atomic.AddInt32(&next, 1))
				if item > count {
					return
				}
				b.claim(slot, item)
				start := time.Now()
				err := recordOneView(ctx, cfg, proxy, b, slot, item)
				b.finish(slot, item, err, time.Since(start))
			}
		}(slot)
	}
	wg.Wait()
	b.stopBoard()

	done, failed := b.stats()
	target := cyan("guns.lol/" + cfg.Username)
	if failed == 0 {
		doneln("Recorded %d views for %s", done, target)
		return
	}
	fmt.Printf("\n  %s %s\n\n", yellow(bold("!")),
		bold(fmt.Sprintf("Recorded %d/%d views for %s — %d failed", done, count, cfg.Username, failed)))
}

// recordOneView runs a single view on a fresh session whose {session} proxy
// token (if any) resolves to a new IP, routing nested logs to the board.
func recordOneView(ctx context.Context, cfg viewConfig, proxy string, b *board, slot, item int) error {
	vctx, cancel := context.WithTimeout(ctx, perViewTimeout)
	defer cancel()

	sess, err := newSession(proxy, false)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}
	sess.onLog = func(level, msg string) { b.log(slot, item, level, msg) }

	return sess.runView(vctx, cfg, &rowReporter{b: b, slot: slot})
}
