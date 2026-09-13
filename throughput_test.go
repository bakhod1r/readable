package readable

import (
	"math"
	"testing"
)

func TestThroughput(t *testing.T) {
	cases := map[uint64]string{0: "0 B/s", 1536: "1.5 KB/s", 125_000_000: "119.21 MB/s"}
	for in, want := range cases {
		if got := Throughput(in); got != want {
			t.Errorf("Throughput(%d) = %q, want %q", in, got, want)
		}
	}
	if got := ThroughputIEC(125_000_000); got != "119.21 MiB/s" {
		t.Errorf("ThroughputIEC = %q", got)
	}
}

func TestBandwidth(t *testing.T) {
	cases := map[uint64]string{
		0: "0 bps", 999: "999 bps", 1500: "1.5 Kbps", 125_000_000: "125 Mbps",
		1e9: "1 Gbps", 2_500_000_000_000: "2.5 Tbps", math.MaxUint64: "18446744.07 Tbps",
	}
	for in, want := range cases {
		if got := Bandwidth(in); got != want {
			t.Errorf("Bandwidth(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestRequestRate(t *testing.T) {
	cases := map[float64]string{15234: "15.23K req/s", 0.5: "0.5 req/s", 0: "0 req/s"}
	for in, want := range cases {
		if got := RequestRate(in); got != want {
			t.Errorf("RequestRate(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestPerMinute(t *testing.T) {
	if got := PerMinute(1200, "req"); got != "1.2K req/min" {
		t.Errorf("got %q", got)
	}
	if got := PerMinute(1200, ""); got != "1.2K/min" {
		t.Errorf("got %q", got)
	}
}

func BenchmarkThroughput(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Throughput(125_000_000)
	}
}

func TestThroughputAllocs(t *testing.T) {
	for name, f := range map[string]func(){
		"Throughput":    func() { allocSink = Throughput(125_000_000) },
		"ThroughputIEC": func() { allocSink = ThroughputIEC(125_000_000) },
		"RequestRate":   func() { allocSink = RequestRate(15234) },
		"PerMinute":     func() { allocSink = PerMinute(1200, "req") },
	} {
		if n := testing.AllocsPerRun(100, f); n != 1 {
			t.Errorf("%s allocs = %v, want 1", name, n)
		}
	}
}
