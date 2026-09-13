package readable_test

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bakhod1r/readable"
)

// The package formats raw machine values for people: counts, sizes,
// durations, money in minor units, percentages and lists.
func Example() {
	fmt.Println(readable.Number(1_500_000))
	fmt.Println(readable.Bytes(1_048_576))
	fmt.Println(readable.Duration(3*time.Hour + 25*time.Minute))
	fmt.Println(readable.Money(150050, "USD"))
	fmt.Println(readable.Percent(0.1534))
	fmt.Println(readable.List([]string{"Go", "Redis", "Kafka"}))
	// Output:
	// 1.5M
	// 1 MB
	// 3h 25m
	// 1,500.50 USD
	// 15.34%
	// Go, Redis, and Kafka
}

func ExampleNumber() {
	fmt.Println(readable.Number(999))
	fmt.Println(readable.Number(12500))
	fmt.Println(readable.Number(1234567))
	fmt.Println(readable.Number(-1500))
	// 999.999K rounds up and is promoted to the next unit.
	fmt.Println(readable.Number(999_999))
	// The full int64 range is safe, including MinInt64.
	fmt.Println(readable.Number(math.MinInt64))
	// Output:
	// 999
	// 12.5K
	// 1.23M
	// -1.5K
	// 1M
	// -9223372.04T
}

func ExampleNumberWithPrecision() {
	fmt.Println(readable.NumberWithPrecision(1234567, 1))
	fmt.Println(readable.NumberWithPrecision(1234567, 3))
	// Precision is a maximum: trailing zeros are trimmed.
	fmt.Println(readable.NumberWithPrecision(1_000_000, 3))
	// Values above MaxPrecision are clamped.
	fmt.Println(readable.NumberWithPrecision(1234567, 20))
	// Output:
	// 1.2M
	// 1.235M
	// 1M
	// 1.234567M
}

func ExampleNumberWithOptions() {
	fixed := readable.NumberOptions{Precision: 2, FixedPrecision: true}
	fmt.Println(readable.NumberWithOptions(1_000_000, fixed))
	fmt.Println(readable.NumberWithOptions(1_500_000, readable.NumberOptions{Precision: 2, FixedPrecision: true, Decimal: ","}))
	// Output:
	// 1.00M
	// 1,50M
}

func ExampleNumberFloat() {
	fmt.Println(readable.NumberFloat(0.5))
	fmt.Println(readable.NumberFloat(1500.25))
	// Output:
	// 0.5
	// 1.5K
}

func ExampleBytes() {
	fmt.Println(readable.Bytes(512))
	fmt.Println(readable.Bytes(1536))
	fmt.Println(readable.Bytes(1048576))
	// Output:
	// 512 B
	// 1.5 KB
	// 1 MB
}

func ExampleBytesIEC() {
	fmt.Println(readable.BytesIEC(1024))
	// Output: 1 KiB
}

func ExampleBytesSI() {
	fmt.Println(readable.BytesSI(1000))
	// Output: 1 kB
}

func ExampleFileSize() {
	fmt.Println(readable.FileSize(5242880))
	fmt.Println(readable.FileSize(1590000))
	// Output:
	// 5 MB
	// 1.5 MB
}

func ExampleDuration() {
	fmt.Println(readable.Duration(3*time.Hour + 25*time.Minute))
	fmt.Println(readable.Duration(125 * time.Millisecond))
	// Output:
	// 3h 25m
	// 125ms
}

func ExampleDurationWithPrecision() {
	fmt.Println(readable.DurationWithPrecision(3*24*time.Hour+12*time.Hour+32*time.Minute, 2))
	// Output: 3d 12h
}

func ExampleDurationLong() {
	fmt.Println(readable.DurationLong(3*time.Hour + 25*time.Minute + 12*time.Second))
	// Output: 3 hours, 25 minutes, 12 seconds
}

func ExampleDurationLongWithPrecision() {
	fmt.Println(readable.DurationLongWithPrecision(3*time.Hour+25*time.Minute+12*time.Second, 2))
	// Output: 3 hours, 25 minutes
}

func ExampleLatency() {
	fmt.Println(readable.Latency(350 * time.Microsecond))
	fmt.Println(readable.Latency(1400 * time.Millisecond))
	// Output:
	// 350µs
	// 1.4s
}

func ExampleMoney() {
	fmt.Println(readable.Money(150000000, "UZS"))
	fmt.Println(readable.Money(150050, "USD"))
	// Output:
	// 150,000,000 UZS
	// 1,500.50 USD
}

func ExampleMoneySymbol() {
	fmt.Println(readable.MoneySymbol(150050, "USD"))
	fmt.Println(readable.MoneySymbol(1000, "UZS"))
	// Output:
	// $1,500.50
	// 1,000 UZS
}

func ExampleMoneyWithPrecision() {
	fmt.Println(readable.MoneyWithPrecision(1234567, "USDT", 4))
	// Output: 123.4567 USDT
}

func ExamplePercent() {
	fmt.Println(readable.Percent(0.1534))
	// Output: 15.34%
}

func ExamplePercentWithPrecision() {
	fmt.Println(readable.PercentWithPrecision(0.123456, 1))
	// Output: 12.3%
}

func ExamplePercentChange() {
	s, err := readable.PercentChange(100, 120)
	fmt.Println(s, err)
	s, _ = readable.PercentChange(120, 100)
	fmt.Println(s)
	// Output:
	// 20% <nil>
	// -16.67%
}

func ExampleErrUndefined() {
	// Change from zero has no defined percentage.
	_, err := readable.PercentChange(0, 5)
	if errors.Is(err, readable.ErrUndefined) {
		fmt.Println("undefined:", err)
	}
	// Output: undefined: readable: undefined result
}

// RelativeTime reads the clock, so its output is not shown here; use
// RelativeTimeFrom in tests.
func ExampleRelativeTime() {
	createdAt := time.Now().Add(-2 * time.Hour)
	fmt.Println(readable.RelativeTime(createdAt)) // 2 hours ago
}

func ExampleRelativeTimeFrom() {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	fmt.Println(readable.RelativeTimeFrom(now, now.Add(-2*time.Hour)))
	fmt.Println(readable.RelativeTimeFrom(now, now.Add(3*24*time.Hour)))
	// Output:
	// 2 hours ago
	// in 3 days
}

func ExampleOrdinal() {
	fmt.Println(readable.Ordinal(1), readable.Ordinal(12), readable.Ordinal(23), readable.Ordinal(-2))
	// Output: 1st 12th 23rd -2nd
}

func ExampleRate() {
	fmt.Println(readable.Rate(125000, time.Second))
	// Output: 125K/s
}

func ExampleRateWithLabel() {
	fmt.Println(readable.RateWithLabel(1200, time.Second, "req"))
	// Output: 1.2K req/s
}

func ExampleBytesRate() {
	fmt.Println(readable.BytesRate(5662310, time.Second))
	// Output: 5.4 MB/s
}
