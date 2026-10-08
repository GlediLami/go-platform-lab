package wait

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestUntil_DoneAfterSomeTries(t *testing.T) {
	calls := 0
	err := Until(context.Background(), time.Millisecond, func(context.Context) (bool, error) {
		calls++
		return calls == 3, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("condition called %d times, want 3", calls)
	}
}

func TestUntil_ChecksImmediately(t *testing.T) {
	start := time.Now()
	err := Until(context.Background(), time.Hour, func(context.Context) (bool, error) {
		return true, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Errorf("Until waited before the first check; it must check right away")
	}
}

func TestUntil_StopsOnError(t *testing.T) {
	boom := errors.New("boom")
	calls := 0
	err := Until(context.Background(), time.Millisecond, func(context.Context) (bool, error) {
		calls++
		return false, boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want it to wrap the condition error", err)
	}
	if calls != 1 {
		t.Errorf("condition called %d times, want 1 (stop on first error)", calls)
	}
}

func TestUntil_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := Until(ctx, 5*time.Millisecond, func(context.Context) (bool, error) {
		return false, nil
	})
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("err = %v, want it to wrap ErrTimeout", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want it to also wrap context.DeadlineExceeded", err)
	}
	if time.Since(start) > time.Second {
		t.Errorf("Until kept going long after the context ended")
	}
}

func TestUntil_PassesContextToCondition(t *testing.T) {
	type key struct{}
	ctx := context.WithValue(context.Background(), key{}, "hello")
	err := Until(ctx, time.Millisecond, func(c context.Context) (bool, error) {
		if c.Value(key{}) != "hello" {
			t.Errorf("condition did not get the caller's context")
		}
		return true, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
