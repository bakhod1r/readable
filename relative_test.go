package readable

import (
	"testing"
	"time"
)

func TestRelativeTimeFrom(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour
	tests := []struct {
		d         time.Duration
		past, fut string
	}{
		{0, "just now", "just now"},
		{4999 * time.Millisecond, "just now", "just now"},
		{5 * time.Second, "5 seconds ago", "in 5 seconds"},
		{59 * time.Second, "59 seconds ago", "in 59 seconds"},
		{60 * time.Second, "1 minute ago", "in 1 minute"},
		{59 * time.Minute, "59 minutes ago", "in 59 minutes"},
		{60 * time.Minute, "1 hour ago", "in 1 hour"},
		{23*time.Hour + 59*time.Minute, "23 hours ago", "in 23 hours"},
		{day, "1 day ago", "in 1 day"},
		{6 * day, "6 days ago", "in 6 days"},
		{7 * day, "1 week ago", "in 1 week"},
		{29 * day, "4 weeks ago", "in 4 weeks"},
		{30 * day, "1 month ago", "in 1 month"},
		{364 * day, "12 months ago", "in 12 months"},
		{365 * day, "1 year ago", "in 1 year"},
		{730 * day, "2 years ago", "in 2 years"},
	}
	for _, tt := range tests {
		if got := RelativeTimeFrom(now, now.Add(-tt.d)); got != tt.past {
			t.Errorf("past %v = %q, want %q", tt.d, got, tt.past)
		}
		if got := RelativeTimeFrom(now, now.Add(tt.d)); got != tt.fut {
			t.Errorf("future %v = %q, want %q", tt.d, got, tt.fut)
		}
	}
}

func TestRelativeTimeFromZonesAndExtremes(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	tz := time.FixedZone("X", 5*3600)
	if got := RelativeTimeFrom(now, now.Add(-2*time.Hour).In(tz)); got != "2 hours ago" {
		t.Errorf("zone: %q", got)
	}
	if got := RelativeTimeFrom(now, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)); got != "292 years ago" {
		t.Errorf("saturate past: %q", got)
	}
	if got := RelativeTimeFrom(now, time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)); got != "in 292 years" {
		t.Errorf("saturate future: %q", got)
	}
}

func TestRelativeTime(t *testing.T) {
	if got := RelativeTime(time.Now().Add(-3 * time.Hour)); got != "3 hours ago" && got != "2 hours ago" {
		t.Errorf("RelativeTime = %q", got)
	}
}

func TestRelativeTimeFromSubSecond(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	if got := RelativeTimeFrom(now, now.Add(-(time.Minute + 999*time.Millisecond))); got != "1 minute ago" {
		t.Errorf("got %q", got)
	}
	if got := RelativeTimeFrom(now, now.Add(time.Nanosecond)); got != "just now" {
		t.Errorf("got %q", got)
	}
}

func BenchmarkRelativeTimeFrom(b *testing.B) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	t := now.Add(-3 * time.Hour)
	b.ReportAllocs()
	for b.Loop() {
		_ = RelativeTimeFrom(now, t)
	}
}
