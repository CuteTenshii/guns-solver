package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// board is a live, multi-slot progress display for the parallel view runner. It
// paints one row per concurrent worker slot plus a running summary, and emits a
// persistent ✓/✗ line for each finished view *above* the live region so results
// accumulate into the scrollback while the slots keep animating below.
//
// When color is disabled (piped/redirected output) it prints only the
// per-completion lines, one per view, keeping logs readable and stable.
//
// All mutating methods take b.mu; render/logLine assume it is already held.
type board struct {
	mu      sync.Mutex
	slots   []slotState
	total   int
	done    int
	failed  int
	frame   int
	printed int // number of live-region lines currently on screen
	live    bool
	stop    chan struct{}
	stopped chan struct{}
}

// slotState is the current work of one worker slot.
type slotState struct {
	active bool
	item   int // 1-based view index this slot is processing
	phase  string
	start  time.Time
}

// startBoard creates a board for total views across concurrency worker slots
// and, on a terminal, starts its render loop.
func startBoard(concurrency, total int) *board {
	b := &board{
		slots:   make([]slotState, concurrency),
		total:   total,
		live:    colorEnabled,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	if b.live {
		go b.loop()
	} else {
		close(b.stopped)
		infoln("recording %d views, %d at a time", total, concurrency)
	}
	return b
}

func (b *board) loop() {
	defer close(b.stopped)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-b.stop:
			b.mu.Lock()
			b.render()
			b.mu.Unlock()
			return
		case <-ticker.C:
			b.mu.Lock()
			b.frame++
			b.render()
			b.mu.Unlock()
		}
	}
}

// claim marks a slot as newly working on item.
func (b *board) claim(slot, item int) {
	b.mu.Lock()
	b.slots[slot] = slotState{active: true, item: item, phase: "starting", start: time.Now()}
	b.mu.Unlock()
}

// setPhase updates the phase label shown for a slot.
func (b *board) setPhase(slot int, phase string) {
	b.mu.Lock()
	if b.slots[slot].active {
		b.slots[slot].phase = phase
	}
	b.mu.Unlock()
}

// finish records a completed view, emits its result line, and idles the slot.
func (b *board) finish(slot, item int, err error, elapsed time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.failed++
	} else {
		b.done++
	}
	b.slots[slot] = slotState{}
	b.logLine(completionLine(item, err, elapsed))
}

// log routes a session's nested status line (view index + message) above the
// live region.
func (b *board) log(slot, item int, level, msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	marker := dim("·")
	if level == "warn" {
		marker = yellow("!")
	}
	b.logLine(fmt.Sprintf("  %s %s %s", marker, dim(fmt.Sprintf("view %d", item)), dim(msg)))
}

// stats returns the completed and failed counts.
func (b *board) stats() (done, failed int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.done, b.failed
}

// stopBoard halts the render loop and clears the live region, leaving only the
// accumulated per-view result lines on screen.
func (b *board) stopBoard() {
	if !b.live {
		return
	}
	close(b.stop)
	<-b.stopped
	b.mu.Lock()
	if b.printed > 0 {
		fmt.Printf("\033[%dA\r\033[J", b.printed)
		b.printed = 0
	}
	b.mu.Unlock()
}

// logLine prints a persistent line above the live region. In live mode it wipes
// the region, prints the line, and lets the next render repaint the region
// beneath it. Assumes b.mu is held.
func (b *board) logLine(line string) {
	if !b.live {
		fmt.Println(line)
		return
	}
	if b.printed > 0 {
		fmt.Printf("\033[%dA\r\033[J", b.printed)
	} else {
		fmt.Print("\r\033[J")
	}
	fmt.Println(line)
	b.printed = 0
	b.render()
}

// render repaints the live region in place. Assumes b.mu is held.
func (b *board) render() {
	if !b.live {
		return
	}
	var sb strings.Builder
	if b.printed > 0 {
		fmt.Fprintf(&sb, "\033[%dA", b.printed)
	}

	lines := 0
	running := 0
	for i := range b.slots {
		s := b.slots[i]
		sb.WriteString("\r\033[K")
		slotLabel := dim(fmt.Sprintf("slot %d", i+1))
		if s.active {
			running++
			marker := cyan(spinnerFrames[b.frame%len(spinnerFrames)])
			elapsed := dim(time.Since(s.start).Round(time.Millisecond).String())
			fmt.Fprintf(&sb, "  %s %s  %s  %s  %s", marker, slotLabel, cyan(fmt.Sprintf("view %d", s.item)), s.phase, elapsed)
		} else {
			fmt.Fprintf(&sb, "  %s %s  %s", dim("·"), slotLabel, dim("idle"))
		}
		sb.WriteString("\n")
		lines++
	}

	sb.WriteString("\r\033[K")
	summary := fmt.Sprintf("%d/%d done", b.done, b.total)
	fmt.Fprintf(&sb, "  %s %s", bold("›"), bold(summary))
	if b.failed > 0 {
		sb.WriteString(dim(" · ") + red(fmt.Sprintf("%d failed", b.failed)))
	}
	sb.WriteString(dim(fmt.Sprintf(" · %d running", running)))
	sb.WriteString("\n")
	lines++

	fmt.Print(sb.String())
	b.printed = lines
}

// completionLine formats the persistent result line for a finished view.
func completionLine(item int, err error, elapsed time.Duration) string {
	tag := dim(fmt.Sprintf("view %d", item))
	el := dim(elapsed.Round(time.Millisecond).String())
	if err != nil {
		return fmt.Sprintf("  %s %s  %s  %s", red("✗"), tag, err.Error(), el)
	}
	return fmt.Sprintf("  %s %s  %s  %s", green("✓"), tag, "recorded", el)
}
