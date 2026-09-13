package readable_test

import (
	"fmt"
	"time"

	"github.com/bakhod1r/readable"
)

func ExampleDurationNatural() {
	fmt.Println(readable.DurationNatural(49*time.Hour + 32*time.Minute))
	fmt.Println(readable.DurationNatural(time.Hour + 5*time.Second))
	fmt.Println(readable.DurationNatural(125 * time.Millisecond))
	// Output:
	// 2 days, 1 hour and 32 minutes
	// 1 hour and 5 seconds
	// 125 milliseconds
}

func ExampleDurationApprox() {
	fmt.Println(readable.DurationApprox(100 * time.Minute))
	fmt.Println(readable.DurationApprox(49 * time.Hour))
	fmt.Println(readable.DurationApprox(30 * time.Second))
	// Output:
	// about 2 hours
	// about 2 days
	// less than a minute
}

func ExampleDurationShort() {
	fmt.Println(readable.DurationShort(90 * time.Second))
	// Output: 1m 30s
}

// Date reads the clock; use DateFrom for deterministic output.
func ExampleDate() {
	fmt.Println(readable.Date(time.Now().Add(-24 * time.Hour))) // Yesterday
}

// Time reads the clock; use TimeFrom for deterministic output.
func ExampleTime() {
	fmt.Println(readable.Time(time.Now().Add(-2 * time.Hour))) // 2 hours ago · Sep 13, 07:20
}

// TimeRange reads the clock; use TimeRangeFrom for deterministic output.
func ExampleTimeRange() {
	start := time.Now().Add(time.Hour)
	fmt.Println(readable.TimeRange(start, start.Add(90*time.Minute))) // 09:30–11:00
}

func ExampleDateFrom() {
	now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
	fmt.Println(readable.DateFrom(now, now.Add(-24*time.Hour)))
	fmt.Println(readable.DateFrom(now, time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC)))
	fmt.Println(readable.DateFrom(now, time.Date(2025, 9, 13, 8, 0, 0, 0, time.UTC)))
	// Output:
	// Yesterday
	// Mar 1
	// Sep 13, 2025
}

func ExampleTimeFrom() {
	now := time.Date(2026, 9, 13, 9, 20, 0, 0, time.UTC)
	fmt.Println(readable.TimeFrom(now, now.Add(-2*time.Hour)))
	// Output: 2 hours ago · Sep 13, 07:20
}

func ExampleTimeRangeFrom() {
	now := time.Date(2026, 9, 13, 8, 0, 0, 0, time.UTC)
	d := func(y int, m time.Month, day, h, min int) time.Time {
		return time.Date(y, m, day, h, min, 0, 0, time.UTC)
	}
	fmt.Println(readable.TimeRangeFrom(now, d(2026, 9, 13, 9, 30), d(2026, 9, 13, 11, 45)))
	fmt.Println(readable.TimeRangeFrom(now, d(2026, 9, 20, 9, 30), d(2026, 9, 20, 11, 45)))
	fmt.Println(readable.TimeRangeFrom(now, d(2026, 9, 13, 23, 30), d(2026, 9, 14, 2, 15)))
	fmt.Println(readable.TimeRangeFrom(now, d(2026, 12, 31, 23, 30), d(2027, 1, 1, 2, 15)))
	// Output:
	// 09:30–11:45
	// Sep 20, 09:30–11:45
	// Sep 13, 23:30 → Sep 14, 02:15
	// Dec 31, 2026, 23:30 → Jan 1, 2027, 02:15
}

func ExampleParseDuration() {
	d, err := readable.ParseDuration("3d 12h")
	fmt.Println(d, err)
	d, _ = readable.ParseDuration("2 days, 1 hour and 32 minutes")
	fmt.Println(d)
	// Output:
	// 84h0m0s <nil>
	// 49h32m0s
}
