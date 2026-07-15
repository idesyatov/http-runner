package generator

import (
	"testing"
	"time"
)

// TestHistogram covers bucketing: empty input, an all-equal degenerate range,
// and a normal spread across equal-width buckets.
func TestHistogram(t *testing.T) {
	// Empty input yields no buckets.
	if got := histogram(nil, 10); got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}

	// All samples equal → a single bucket holding everything.
	equal := []time.Duration{5 * time.Millisecond, 5 * time.Millisecond, 5 * time.Millisecond}
	single := histogram(equal, 10)
	if len(single) != 1 || single[0].Count != 3 {
		t.Fatalf("expected one bucket with count 3, got %v", single)
	}

	// A 0..9ms spread into 10 buckets (width 0.9ms): all counts sum to the input
	// size and the min/max land in the first/last buckets.
	var spread []time.Duration
	for i := 0; i < 10; i++ {
		spread = append(spread, time.Duration(i)*time.Millisecond)
	}
	buckets := histogram(spread, 10)
	if len(buckets) != 10 {
		t.Fatalf("expected 10 buckets, got %d", len(buckets))
	}
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != len(spread) {
		t.Errorf("expected bucket counts to sum to %d, got %d", len(spread), total)
	}
	if buckets[0].Count == 0 {
		t.Errorf("expected the minimum sample in the first bucket")
	}
	if buckets[len(buckets)-1].Count == 0 {
		t.Errorf("expected the maximum sample in the last bucket")
	}
}
