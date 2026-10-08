package health

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func names(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("component-%02d", i)
	}
	return out
}

func TestCheckAll_ReturnsEveryResult(t *testing.T) {
	broken := errors.New("not ready")
	results := CheckAll(context.Background(), names(10), 3, func(_ context.Context, name string) error {
		if name == "component-04" || name == "component-07" {
			return broken
		}
		return nil
	})

	if len(results) != 10 {
		t.Fatalf("got %d results, want 10", len(results))
	}
	got := Unhealthy(results)
	slices.Sort(got)
	if want := []string{"component-04", "component-07"}; !slices.Equal(got, want) {
		t.Errorf("Unhealthy() = %v, want %v", got, want)
	}
}

func TestCheckAll_RespectsLimit(t *testing.T) {
	var running, peak atomic.Int32

	CheckAll(context.Background(), names(20), 4, func(context.Context, string) error {
		now := running.Add(1)
		for {
			old := peak.Load()
			if now <= old || peak.CompareAndSwap(old, now) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		running.Add(-1)
		return nil
	})

	if p := peak.Load(); p > 4 {
		t.Errorf("up to %d checks ran at once, limit is 4", p)
	}
	if p := peak.Load(); p < 2 {
		t.Errorf("only %d check ran at a time; checks must run in parallel", p)
	}
}

func TestCheckAll_IsFasterThanSerial(t *testing.T) {
	start := time.Now()
	CheckAll(context.Background(), names(10), 10, func(context.Context, string) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})
	// Serial would take 500ms.
	if d := time.Since(start); d > 250*time.Millisecond {
		t.Errorf("took %v; checks do not seem to run in parallel", d)
	}
}

func TestCheckAll_ZeroLimitStillWorks(t *testing.T) {
	results := CheckAll(context.Background(), names(3), 0, func(context.Context, string) error { return nil })
	if len(results) != 3 {
		t.Errorf("got %d results, want 3 (treat a limit below 1 as 1)", len(results))
	}
}
