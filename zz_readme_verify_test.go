package readable_test

import (
	"errors"
	"math"
	"slices"
	"testing"
	"time"

	r "github.com/bakhod1r/readable"
)

func TestZZReadme(t *testing.T) {
	d := 49*time.Hour + 32*time.Minute
	now := time.Date(2026, 9, 13, 9, 20, 0, 0, time.UTC)
	at := func(h, m int) time.Time { return time.Date(2026, 9, 13, h, m, 0, 0, time.UTC) }
	items := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	pc1, e1 := r.PercentChange(100, 120)
	pc2, e2 := r.PercentChange(0, 5)
	cases := [][2]string{
		{r.Number(1_500_000), "1.5M"}, {r.Bytes(1_048_576), "1 MB"}, {r.Duration(3*time.Hour + 25*time.Minute), "3h 25m"},
		{r.Money(150050, "USD"), "1,500.50 USD"}, {r.Percent(0.1534), "15.34%"}, {r.List([]string{"Go", "Redis", "Kafka"}), "Go, Redis, and Kafka"},
		{r.MaskEmail("john.doe@gmail.com"), "j***@gmail.com"},
		{r.Number(12500), "12.5K"}, {r.Number(999_999), "1M"}, {r.Number(math.MinInt64), "-9223372.04T"},
		{r.NumberWithPrecision(1234567, 3), "1.235M"}, {r.NumberWithOptions(1_000_000, r.NumberOptions{Precision: 2, FixedPrecision: true}), "1.00M"},
		{r.NumberFloat(1500.25), "1.5K"}, {r.NumberWords(1234567), "1.23 million"}, {r.Ordinal(23), "23rd"}, {r.Ordinal(12), "12th"},
		{r.Count(3, "match"), "3 matches"}, {r.CountPlural(2, "person", "people"), "2 people"},
		{r.Bytes(1536), "1.5 KB"}, {r.BytesIEC(1024), "1 KiB"}, {r.BytesSI(1000), "1 kB"}, {r.FileSize(5242880), "5 MB"},
		{r.Throughput(125_000_000), "119.21 MB/s"}, {r.ThroughputIEC(125_000_000), "119.21 MiB/s"}, {r.Bandwidth(125_000_000), "125 Mbps"},
		{r.BytesRate(5662310, time.Second), "5.4 MB/s"}, {r.Rate(125000, time.Second), "125K/s"}, {r.RateWithLabel(1200, time.Second, "req"), "1.2K req/s"},
		{r.RequestRate(15234), "15.23K req/s"}, {r.PerMinute(1200, "req"), "1.2K req/min"},
		{r.Duration(d), "2d 1h 32m"}, {r.DurationWithOptions(d, r.DurationOptions{Units: 2}), "2d 1h"}, {r.DurationLong(d), "2 days, 1 hour, 32 minutes"},
		{r.DurationNatural(d), "2 days, 1 hour and 32 minutes"}, {r.DurationApprox(d), "about 2 days"}, {r.Latency(1500 * time.Microsecond), "1.5ms"},
		{r.RelativeTimeFrom(now, now.Add(-2*time.Hour)), "2 hours ago"}, {r.RelativeTimeFrom(now, now.Add(72*time.Hour)), "in 3 days"},
		{r.DateFrom(now, now.Add(-24*time.Hour)), "Yesterday"}, {r.TimeFrom(now, now.Add(-2*time.Hour)), "2 hours ago · Sep 13, 07:20"},
		{r.TimeRangeFrom(at(8, 0), at(9, 30), at(11, 45)), "09:30–11:45"},
		{r.Money(150_000_000, "UZS"), "150,000,000 UZS"}, {r.MoneySymbol(150050, "USD"), "$1,500.50"}, {r.MoneyWithPrecision(1234567, "USDT", 4), "123.4567 USDT"},
		{r.MoneyCompact(150_000_000, "USD"), "$1.5M"}, {r.MoneyAccounting(-1_500_000, "UZS"), "(1,500,000 UZS)"}, {r.MoneyChange(-1_500_000, "UZS"), "−1,500,000 UZS"},
		{r.PercentWithPrecision(0.123456, 1), "12.3%"}, {pc1, "20%"}, {pc2, ""}, {r.Progress(999, 1000), "99%"}, {r.ProgressBar(73, 100, 20), "██████████████░░░░░░ 73%"},
		{r.List([]string{"Go", "PostgreSQL", "Redis"}), "Go, PostgreSQL, and Redis"}, {r.ListWithOptions(items, r.ListOptions{Limit: 3}), "Go, Redis, Kafka, and 2 more"},
		{r.Plural(2, "city"), "cities"}, {r.Humanize("created_at"), "Created at"}, {r.Humanize("HTTPServer"), "HTTP server"}, {r.Enum("IN_PROGRESS"), "In progress"},
		{r.Bool(true), "Yes"}, {r.BoolLabel(false, "Enabled", "Disabled"), "Disabled"},
		{r.MaskEmail("jo@gmail.com"), "***@gmail.com"}, {r.MaskPhone("+998 90 123-45-67"), "+998******567"}, {r.MaskCard("8600 1234 1234 5678"), "**** **** **** 5678"},
		{r.MaskToken("sk_live_abc123456789"), "sk_live_****6789"}, {r.MaskToken("sk_live_1234"), "****"}, {r.MaskIP("192.168.1.42"), "192.168.*.*"},
		{r.Mask("1234567890", 2, 2), "12******90"}, {r.ID(987654321234567), "987-654-321-234-567"},
		{r.ShortUUID("550e8400-e29b-41d4-a716-446655440000"), "550e...0000"}, {r.Truncate("a8f91234abcd", 4, 2), "a8f9...cd"},
		// design notes / guarantees
		{r.NumberWithPrecision(1234567, 20), "1.234567M"}, {r.Money(100, "XYZ"), "100 XYZ"}, {r.MaskEmail("not-an-email"), "***"},
		{r.MaskPhone("12345"), "*****"}, {r.MaskPhone("12a45678"), "***"}, {r.MaskCard("1234"), "****"}, {r.MaskIP("nope"), "***"}, {r.Mask("secret", 3, 3), "******"},
		{r.NumberFloat(math.NaN()), "NaN"},
	}
	if e1 != nil || !errors.Is(e2, r.ErrUndefined) || r.DefaultPrecision != 2 || r.MaxPrecision != 9 {
		t.Errorf("errs/consts: %v %v", e1, e2)
	}
	pb, eb := r.ParseBytes("1.5 KB")
	pd, ed := r.ParseDuration("2 days, 1 hour and 32 minutes")
	_, en := r.ParseNumber("twelve")
	if pb != 1536 || pd != d || eb != nil || ed != nil || !errors.Is(en, r.ErrSyntax) {
		t.Errorf("parse: %d %v %v %v %v", pb, pd, eb, ed, en)
	}
	extra := [][2]string{
		{r.Redact("login john.doe@gmail.com password=hunter2"), "login j***@gmail.com password=****"},
		{r.Uzbek.RelativeTimeFrom(now, now.Add(-3*time.Minute)), "3 daqiqa oldin"},
		{r.Russian.DurationLong(49*time.Hour + 32*time.Minute), "2 дня, 1 час, 32 минуты"},
		{r.ETA(25, 100, time.Minute), "~3m left"}, {r.Roman(2026), "MMXXVI"},
		{r.German.RelativeTimeFrom(now, now.Add(-72*time.Hour)), "vor 3 Tagen"},
		{r.Russian.List([]string{"Go", "Redis", "Kafka"}), "Go, Redis и Kafka"},
		{r.Turkish.Percent(0.1534), "%15,34"},
	}
	for i, c := range slices.Concat(cases, extra) {
		if c[0] != c[1] {
			t.Errorf("case %d: got %q want %q", i, c[0], c[1])
		}
	}
}
