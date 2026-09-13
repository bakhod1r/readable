package readable_test

import (
	"errors"
	"fmt"

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

func ExampleParseBytes() {
	n, err := readable.ParseBytes("1.5 KB")
	fmt.Println(n, err)
	n, _ = readable.ParseBytes("1 MiB")
	fmt.Println(n)
	// Output:
	// 1536 <nil>
	// 1048576
}

func ExampleParseBytesSI() {
	n, _ := readable.ParseBytesSI("1 kB")
	fmt.Println(n)
	n, _ = readable.ParseBytesSI("1 KiB")
	fmt.Println(n)
	// Output:
	// 1000
	// 1024
}

func ExampleParseNumber() {
	n, _ := readable.ParseNumber("12.5K")
	fmt.Println(n)
	n, _ = readable.ParseNumber("1,234,567")
	fmt.Println(n)
	// Output:
	// 12500
	// 1234567
}

func ExampleErrSyntax() {
	_, err := readable.ParseNumber("twelve")
	fmt.Println(errors.Is(err, readable.ErrSyntax))
	_, err = readable.ParseBytes("5 parsecs")
	fmt.Println(errors.Is(err, readable.ErrUnit))
	_, err = readable.ParseBytes("-1 KB")
	fmt.Println(errors.Is(err, readable.ErrRange))
	// Output:
	// true
	// true
	// true
}
