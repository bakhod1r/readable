package readable_test

import (
	"fmt"
	"time"

	"github.com/bakhod1r/readable"
)

func ExampleMoneyCompact() {
	fmt.Println(readable.MoneyCompact(150_000_000, "USD"))
	fmt.Println(readable.MoneyCompact(1_500_000, "UZS"))
	// Output:
	// $1.5M
	// 1.5M UZS
}

func ExampleMoneyAccounting() {
	fmt.Println(readable.MoneyAccounting(-1_500_000, "UZS"))
	// Output: (1,500,000 UZS)
}

func ExampleMoneyChange() {
	fmt.Println(readable.MoneyChange(1_500_000, "UZS"))
	fmt.Println(readable.MoneyChange(-1_500_000, "UZS"))
	// Output:
	// +1,500,000 UZS
	// −1,500,000 UZS
}

func ExampleThroughput() {
	fmt.Println(readable.Throughput(125_000_000))
	// Output: 119.21 MB/s
}

func ExampleThroughputIEC() {
	fmt.Println(readable.ThroughputIEC(125_000_000))
	// Output: 119.21 MiB/s
}

func ExampleBandwidth() {
	fmt.Println(readable.Bandwidth(125_000_000))
	// Output: 125 Mbps
}

func ExampleRequestRate() {
	fmt.Println(readable.RequestRate(15234))
	// Output: 15.23K req/s
}

func ExamplePerMinute() {
	fmt.Println(readable.PerMinute(1200, "req"))
	// Output: 1.2K req/min
}

func ExampleLogDuration() {
	fmt.Println(readable.LogDuration(1532 * time.Millisecond))
	// Output: 1.532s
}

func ExampleLogBytes() {
	fmt.Println(readable.LogBytes(12345678))
	// Output: 11.77MiB
}

func ExampleLogNumber() {
	fmt.Println(readable.LogNumber(1234567))
	// Output: 1234567
}

func ExampleProgress() {
	fmt.Println(readable.Progress(999, 1000))
	// Output: 99%
}

func ExampleProgressBar() {
	fmt.Println(readable.ProgressBar(73, 100, 20))
	// Output: ██████████████░░░░░░ 73%
}

func ExampleNumberWords() {
	fmt.Println(readable.NumberWords(1234567))
	// Output: 1.23 million
}
