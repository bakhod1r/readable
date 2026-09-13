package readable

import (
	"testing"
	"time"
)

var dateTestZone = time.FixedZone("UTC+5", 5*3600)

func TestDateFrom(t *testing.T) {
	now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		now, t time.Time
		want   string
	}{
		{"today start", now, time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC), "Today"},
		{"today end", now, time.Date(2026, 9, 13, 23, 59, 59, 0, time.UTC), "Today"},
		{"yesterday end", now, time.Date(2026, 9, 12, 23, 59, 59, 0, time.UTC), "Yesterday"},
		{"tomorrow start", now, time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC), "Tomorrow"},
		{"same year", now, time.Date(2026, 9, 11, 10, 0, 0, 0, time.UTC), "Sep 11"},
		{"other year", now, time.Date(2025, 9, 13, 10, 0, 0, 0, time.UTC), "Sep 13, 2025"},
		{"year boundary yesterday", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "Yesterday"},
		{"year boundary tomorrow", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), "Tomorrow"},
		{"leap day tomorrow", time.Date(2028, 2, 28, 9, 0, 0, 0, time.UTC), time.Date(2028, 2, 29, 9, 0, 0, 0, time.UTC), "Tomorrow"},
		{"leap day yesterday", time.Date(2028, 3, 1, 9, 0, 0, 0, time.UTC), time.Date(2028, 2, 29, 9, 0, 0, 0, time.UTC), "Yesterday"},
		{"leap day same year", now.AddDate(2, 0, 0), time.Date(2028, 2, 29, 9, 0, 0, 0, time.UTC), "Feb 29"},
		// now 20:00 UTC is 01:00 on Sep 14 in UTC+5, so t on Sep 14 there is today.
		{"zone of t", time.Date(2026, 9, 13, 20, 0, 0, 0, time.UTC), time.Date(2026, 9, 14, 0, 30, 0, 0, dateTestZone), "Today"},
		{"year 1 yesterday", time.Date(1, 1, 2, 0, 0, 0, 0, time.UTC), time.Date(1, 1, 1, 23, 59, 59, 0, time.UTC), "Yesterday"},
		{"year 1 other year", now, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), "Jan 1, 0001"},
		{"before 1970 today", time.Date(1969, 12, 31, 0, 0, 1, 0, time.UTC), time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC), "Today"},
		{"epoch boundary tomorrow", time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "Tomorrow"},
		{"epoch boundary yesterday", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC), "Yesterday"},
		{"year 9999 tomorrow", time.Date(9999, 12, 30, 12, 0, 0, 0, time.UTC), time.Date(9999, 12, 31, 12, 0, 0, 0, time.UTC), "Tomorrow"},
		{"year 9999 same year", time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(9999, 12, 31, 12, 0, 0, 0, time.UTC), "Dec 31"},
		{"far future tomorrow", time.Date(123456, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(123456, 3, 2, 0, 0, 0, 0, time.UTC), "Tomorrow"},
		{"zone of t yesterday", time.Date(2026, 9, 13, 20, 0, 0, 0, time.UTC), time.Date(2026, 9, 13, 23, 0, 0, 0, dateTestZone), "Yesterday"},
	}
	for _, tt := range tests {
		if got := DateFrom(tt.now, tt.t); got != tt.want {
			t.Errorf("%s: DateFrom = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestDateCivilDay(t *testing.T) {
	tests := []struct {
		t    time.Time
		want int64
	}{
		{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(1970, 1, 1, 23, 59, 59, 0, time.UTC), 0},
		{time.Date(1969, 12, 31, 0, 0, 0, 0, time.UTC), -1},
		{time.Date(1969, 12, 31, 23, 59, 59, 999999999, time.UTC), -1},
		{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), -719162},
		{time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC), 2932896},
		{time.Date(1970, 1, 1, 0, 30, 0, 0, dateTestZone), 0},
	}
	for _, tt := range tests {
		if got := dateCivilDay(tt.t); got != tt.want {
			t.Errorf("dateCivilDay(%v) = %d, want %d", tt.t, got, tt.want)
		}
	}
}

func TestDate(t *testing.T) {
	if got := Date(time.Now()); got != "Today" {
		t.Errorf("Date(time.Now()) = %q, want Today", got)
	}
}

func TestTimeFrom(t *testing.T) {
	now := time.Date(2026, 9, 13, 9, 20, 0, 0, time.UTC)
	tests := []struct {
		t    time.Time
		want string
	}{
		{time.Date(2026, 9, 13, 7, 20, 0, 0, time.UTC), "2 hours ago · Sep 13, 07:20"},
		{time.Date(2026, 9, 13, 12, 20, 0, 0, dateTestZone), "2 hours ago · Sep 13, 12:20"},
		{time.Date(2025, 9, 13, 7, 20, 0, 0, time.UTC), "1 year ago · Sep 13, 2025, 07:20"},
		{time.Date(2026, 9, 14, 9, 20, 0, 0, time.UTC), "in 1 day · Sep 14, 09:20"},
	}
	for _, tt := range tests {
		if got := TimeFrom(now, tt.t); got != tt.want {
			t.Errorf("TimeFrom(%v) = %q, want %q", tt.t, got, tt.want)
		}
	}
}

func TestTime(t *testing.T) {
	ts := time.Now().Add(-3 * time.Hour)
	if got, want := Time(ts), TimeFrom(time.Now(), ts); got[:len("3 hours ago")] != want[:len("3 hours ago")] {
		t.Errorf("Time = %q, want prefix of %q", got, want)
	}
}

func TestTimeRangeFrom(t *testing.T) {
	now := time.Date(2026, 9, 13, 8, 0, 0, 0, time.UTC)
	d := func(y int, m time.Month, day, h, min int) time.Time {
		return time.Date(y, m, day, h, min, 0, 0, time.UTC)
	}
	tests := []struct {
		name       string
		now        time.Time
		start, end time.Time
		want       string
	}{
		{"today", now, d(2026, 9, 13, 9, 30), d(2026, 9, 13, 11, 45), "09:30–11:45"},
		{"same day not today", now, d(2026, 9, 20, 9, 30), d(2026, 9, 20, 11, 45), "Sep 20, 09:30–11:45"},
		{"same day other year", now, d(2025, 9, 20, 9, 30), d(2025, 9, 20, 11, 45), "Sep 20, 2025, 09:30–11:45"},
		{"different days", now, d(2026, 9, 13, 23, 30), d(2026, 9, 14, 2, 15), "Sep 13, 23:30 → Sep 14, 02:15"},
		{"midnight boundary", now, d(2026, 9, 13, 23, 59), d(2026, 9, 14, 0, 0), "Sep 13, 23:59 → Sep 14, 00:00"},
		{"years", now, d(2026, 12, 31, 23, 30), d(2027, 1, 1, 2, 15), "Dec 31, 2026, 23:30 → Jan 1, 2027, 02:15"},
		{"swapped", now, d(2026, 9, 13, 11, 45), d(2026, 9, 13, 9, 30), "09:30–11:45"},
		{"equal", now, d(2026, 9, 13, 9, 30), d(2026, 9, 13, 9, 30), "09:30–09:30"},
		{"leap day", time.Date(2028, 2, 29, 1, 0, 0, 0, time.UTC), d(2028, 2, 28, 22, 0), d(2028, 2, 29, 1, 0), "Feb 28, 22:00 → Feb 29, 01:00"},
		// end is shown in start's location: 04:00 UTC == 09:00 UTC+5.
		{"end zone", now, time.Date(2026, 9, 13, 8, 0, 0, 0, dateTestZone), d(2026, 9, 13, 4, 0), "08:00–09:00"},
		// start 01:00 UTC+5 on Sep 14 is 20:00 UTC Sep 13; today is judged in start's zone.
		{"start zone not today", now, time.Date(2026, 9, 14, 1, 0, 0, 0, dateTestZone), time.Date(2026, 9, 14, 2, 0, 0, 0, dateTestZone), "Sep 14, 01:00–02:00"},
	}
	for _, tt := range tests {
		if got := TimeRangeFrom(tt.now, tt.start, tt.end); got != tt.want {
			t.Errorf("%s: TimeRangeFrom = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestTimeRange(t *testing.T) {
	s := time.Date(2001, 9, 13, 9, 30, 0, 0, time.UTC)
	if got, want := TimeRange(s, s.Add(time.Hour)), "Sep 13, 2001, 09:30–10:30"; got != want {
		t.Errorf("TimeRange = %q, want %q", got, want)
	}
}

func BenchmarkDateFrom(b *testing.B) {
	b.ReportAllocs()
	now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
	ts := time.Date(2026, 9, 11, 10, 0, 0, 0, time.UTC)
	for b.Loop() {
		_ = DateFrom(now, ts)
	}
}

func BenchmarkTimeRange(b *testing.B) {
	b.ReportAllocs()
	now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
	s := time.Date(2026, 9, 13, 23, 30, 0, 0, time.UTC)
	for b.Loop() {
		_ = TimeRangeFrom(now, s, s.Add(3*time.Hour))
	}
}
