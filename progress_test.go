package readable

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestProgress(t *testing.T) {
	tests := []struct {
		done, total int64
		want        string
	}{
		{73, 100, "73%"},
		{999, 1000, "99%"},
		{1000, 1000, "100%"},
		{0, 0, "0%"},
		{5, -1, "0%"},
		{-5, 10, "0%"},
		{20, 10, "100%"},
		{math.MaxInt64 - 1, math.MaxInt64, "99%"},
		{math.MaxInt64 / 2, math.MaxInt64, "49%"},
	}
	for _, tt := range tests {
		if got := Progress(tt.done, tt.total); got != tt.want {
			t.Errorf("Progress(%d, %d) = %q, want %q", tt.done, tt.total, got, tt.want)
		}
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		done, total int64
		width       int
		want        string
	}{
		{73, 100, 20, "██████████████░░░░░░ 73%"},
		{73, 100, 0, "██████████████░░░░░░ 73%"},
		{1, 2, 4, "██░░ 50%"},
		{0, 0, 3, "░░░ 0%"},
		{10, 10, 3, "███ 100%"},
	}
	for _, tt := range tests {
		if got := ProgressBar(tt.done, tt.total, tt.width); got != tt.want {
			t.Errorf("ProgressBar(%d, %d, %d) = %q, want %q", tt.done, tt.total, tt.width, got, tt.want)
		}
	}
	if got := ProgressBar(1, 1, 1000); strings.Count(got, "█") != 200 {
		t.Errorf("width not capped: %d cells", strings.Count(got, "█"))
	}
}

func FuzzProgress(f *testing.F) {
	f.Add(int64(73), int64(100))
	f.Add(int64(math.MaxInt64), int64(math.MaxInt64))
	f.Add(int64(math.MinInt64), int64(-1))
	f.Fuzz(func(t *testing.T, done, total int64) {
		s := Progress(done, total)
		if !strings.HasSuffix(s, "%") {
			t.Fatalf("no %%: %q", s)
		}
		p, err := strconv.Atoi(strings.TrimSuffix(s, "%"))
		if err != nil || p < 0 || p > 100 {
			t.Fatalf("bad percent %q", s)
		}
		if p == 100 && done < total {
			t.Fatalf("100%% before done: %d/%d", done, total)
		}
		_ = ProgressBar(done, total, int(done%300))
	})
}

func BenchmarkProgressBar(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = ProgressBar(73, 100, 20)
	}
}

func TestProgressBarAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { allocSink = ProgressBar(73, 100, 20) }); n != 1 {
		t.Errorf("ProgressBar allocs = %v, want 1", n)
	}
}
