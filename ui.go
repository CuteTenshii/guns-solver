package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// colorEnabled reports whether ANSI styling should be emitted. It is disabled
// when stdout is not a terminal, when NO_COLOR is set (https://no-color.org),
// or for dumb terminals, so redirected output stays clean and machine-readable.
var colorEnabled = detectColor()

func detectColor() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// ANSI SGR codes, applied only when colorEnabled.
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiCyan    = "\033[36m"
	ansiMagenta = "\033[35m"
)

func style(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + ansiReset
}

func bold(s string) string    { return style(ansiBold, s) }
func dim(s string) string     { return style(ansiDim, s) }
func red(s string) string     { return style(ansiRed, s) }
func green(s string) string   { return style(ansiGreen, s) }
func yellow(s string) string  { return style(ansiYellow, s) }
func cyan(s string) string    { return style(ansiCyan, s) }
func magenta(s string) string { return style(ansiMagenta, s) }

// banner prints the tool header.
func banner() {
	fmt.Println()
	fmt.Printf("  %s %s\n", magenta(bold("guns")), dim("· lol view solver"))
	fmt.Println()
}

// interruptSpinner clears a running spinner's current frame so a log line can be
// printed above it without corrupting the animation; the spinner redraws itself
// on the next tick.
func interruptSpinner() {
	if colorEnabled && activeSpinner != nil {
		fmt.Print("\r\033[K")
	}
}

// infoln prints a secondary, dimmed informational line (nested detail under a
// step). Used by subsystems like the clearance and cf_clearance flows.
func infoln(format string, a ...any) {
	interruptSpinner()
	fmt.Printf("  %s %s\n", dim("·"), dim(fmt.Sprintf(format, a...)))
}

// warnln prints a highlighted notice that is not a hard failure.
func warnln(format string, a ...any) {
	interruptSpinner()
	fmt.Printf("  %s %s\n", yellow("!"), fmt.Sprintf(format, a...))
}

// doneln prints a final success line, used once the whole flow completes.
func doneln(format string, a ...any) {
	fmt.Printf("\n  %s %s\n\n", green(bold("✓")), bold(fmt.Sprintf(format, a...)))
}

// spinner animates a single in-progress step on one line and resolves it to a
// success (✓) or failure (✗) marker. On a non-terminal it degrades to a plain
// static line so piped/logged output stays readable.
type spinner struct {
	label string
	stop  chan struct{}
	done  chan struct{}
	mu    sync.Mutex
	start time.Time
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// startSpinner begins animating a step labelled with the given message.
func startSpinner(label string) *spinner {
	s := &spinner{
		label: label,
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
		start: time.Now(),
	}
	activeSpinner = s
	if !colorEnabled {
		fmt.Printf("  %s %s\n", dim("→"), label)
		close(s.done)
		return s
	}
	go s.loop()
	return s
}

func (s *spinner) loop() {
	defer close(s.done)
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()
	frame := 0
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.mu.Lock()
			label := s.label
			s.mu.Unlock()
			fmt.Printf("\r  %s %s", cyan(spinnerFrames[frame%len(spinnerFrames)]), label)
			frame++
		}
	}
}

// update changes the label of a running spinner in place.
func (s *spinner) update(label string) {
	s.mu.Lock()
	s.label = label
	s.mu.Unlock()
}

func (s *spinner) clearLine() {
	if colorEnabled {
		fmt.Print("\r\033[K")
	}
}

// resolve stops the animation and emits a final marker line.
func (s *spinner) resolve(marker, label, detail string) {
	if colorEnabled {
		select {
		case <-s.stop:
		default:
			close(s.stop)
		}
		<-s.done
		s.clearLine()
	}
	if activeSpinner == s {
		activeSpinner = nil
	}
	line := fmt.Sprintf("  %s %s", marker, label)
	if detail != "" {
		line += "  " + dim(detail)
	}
	fmt.Println(line)
}

// succeed marks the step done, appending the elapsed time.
func (s *spinner) succeed(label string) {
	s.resolve(green("✓"), label, s.elapsed())
}

// succeedDetail marks the step done with a custom trailing detail.
func (s *spinner) succeedDetail(label, detail string) {
	s.resolve(green("✓"), label, detail)
}

// fail marks the step failed using its current label; the caller typically
// exits afterwards.
func (s *spinner) fail() {
	s.mu.Lock()
	label := s.label
	s.mu.Unlock()
	s.resolve(red("✗"), label, "")
}

func (s *spinner) elapsed() string {
	return time.Since(s.start).Round(time.Millisecond).String()
}

// activeSpinner tracks the spinner currently animating, so fatalf can tear it
// down cleanly before printing an error instead of leaving a dangling frame.
var activeSpinner *spinner

// truncateMiddle shortens long tokens for display, keeping both ends.
func truncateMiddle(s string, max int) string {
	if len(s) <= max || max < 5 {
		return s
	}
	half := (max - 1) / 2
	return s[:half] + "…" + s[len(s)-half:]
}
