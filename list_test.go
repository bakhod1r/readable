package readable

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"Go"}, "Go"},
		{[]string{"Go", "PostgreSQL"}, "Go and PostgreSQL"},
		{[]string{"Go", "PostgreSQL", "Redis"}, "Go, PostgreSQL, and Redis"},
		{[]string{"a", "b", "c", "d"}, "a, b, c, and d"},
	}
	for _, tt := range tests {
		if got := List(tt.in); got != tt.want {
			t.Errorf("List(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestListWithOptions(t *testing.T) {
	five := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	tests := []struct {
		name string
		in   []string
		opts ListOptions
		want string
	}{
		{"limit", five, ListOptions{Limit: 3}, "Go, Redis, Kafka, and 2 more"},
		{"limit one", five, ListOptions{Limit: 1}, "Go and 4 more"},
		{"limit equal", five[:3], ListOptions{Limit: 3}, "Go, Redis, and Kafka"},
		{"limit over", five[:2], ListOptions{Limit: 3}, "Go and Redis"},
		{"negative limit", five[:3], ListOptions{Limit: -1}, "Go, Redis, and Kafka"},
		{"min int limit", five, ListOptions{Limit: math.MinInt}, "Go, Redis, Kafka, NATS, and Postgres"},
		{"negative limit single", five[:1], ListOptions{Limit: -5}, "Go"},
		{"negative limit two", five[:2], ListOptions{Limit: -5, Conjunction: "or"}, "Go or Redis"},
		{"limit single item", five[:1], ListOptions{Limit: 1}, "Go"},
		{"limit two of three", five[:3], ListOptions{Limit: 2}, "Go, Redis, and 1 more"},
		{"empty items", []string{"", ""}, ListOptions{}, " and "},
		{"empty ignored when non-empty", five[:1], ListOptions{Empty: "none"}, "Go"},
		{"custom conjunction", five[:3], ListOptions{Conjunction: "&"}, "Go, Redis, & Kafka"},
		{"or", []string{"Go", "Redis"}, ListOptions{Conjunction: "or"}, "Go or Redis"},
		{"or three", five[:3], ListOptions{Conjunction: "or"}, "Go, Redis, or Kafka"},
		{"empty default", nil, ListOptions{}, ""},
		{"empty label", nil, ListOptions{Empty: "none"}, "none"},
		{"limit or", five, ListOptions{Limit: 2, Conjunction: "or"}, "Go, Redis, or 3 more"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListWithOptions(tt.in, tt.opts); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkList(b *testing.B) {
	items := []string{"Go", "PostgreSQL", "Redis", "Kafka"}
	b.ReportAllocs()
	for b.Loop() {
		_ = List(items)
	}
}

func TestListWithOptionsAllocs(t *testing.T) {
	five := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	many := make([]string, 1500)
	for i := range many {
		many[i] = "x"
	}
	cases := map[string]struct {
		in   []string
		opts ListOptions
	}{
		"two":        {five[:2], ListOptions{}},
		"five":       {five, ListOptions{}},
		"limit one":  {five, ListOptions{Limit: 1}},
		"limit":      {five, ListOptions{Limit: 3, Conjunction: "or"}},
		"many limit": {many, ListOptions{Limit: 2}},
	}
	for name, c := range cases {
		want := ListWithOptions(c.in, c.opts)
		if n := testing.AllocsPerRun(100, func() { _ = ListWithOptions(c.in, c.opts) }); n > 1 {
			t.Errorf("%s: allocs = %v, want <= 1", name, n)
		}
		if name == "many limit" && want != "x, x, and 1498 more" {
			t.Errorf("%s: got %q", name, want)
		}
	}
	if n := testing.AllocsPerRun(100, func() { _ = List(five[:1]) }); n != 0 {
		t.Errorf("single item allocs = %v, want 0", n)
	}
}

func TestListWithOptionsMatchesNaive(t *testing.T) {
	// Compare every limit/length combination with an independent construction.
	five := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	for limit := -1; limit <= len(five)+1; limit++ {
		for n := 1; n <= len(five); n++ {
			got := ListWithOptions(five[:n], ListOptions{Limit: limit})
			shown := five[:n]
			var parts []string
			if limit > 0 && limit < n {
				shown = five[:limit]
				parts = append(append(parts, shown...), strconv.Itoa(n-limit)+" more")
			} else {
				parts = append(parts, shown...)
			}
			var want string
			switch len(parts) {
			case 1:
				want = parts[0]
			case 2:
				want = parts[0] + " and " + parts[1]
			default:
				want = strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
			}
			if got != want {
				t.Errorf("n=%d limit=%d: got %q, want %q", n, limit, got, want)
			}
		}
	}
}

func BenchmarkListWithOptionsLimit(b *testing.B) {
	items := []string{"Go", "PostgreSQL", "Redis", "Kafka", "NATS"}
	opts := ListOptions{Limit: 2, Conjunction: "or"}
	b.ReportAllocs()
	for b.Loop() {
		_ = ListWithOptions(items, opts)
	}
}

func BenchmarkListTwo(b *testing.B) {
	items := []string{"Go", "PostgreSQL"}
	b.ReportAllocs()
	for b.Loop() {
		_ = List(items)
	}
}

func TestListDocExamples(t *testing.T) {
	five := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	two := five[:2]
	tests := []struct{ got, want string }{
		{List([]string{"Go"}), "Go"},
		{List([]string{"Go", "PostgreSQL"}), "Go and PostgreSQL"},
		{List([]string{"Go", "PostgreSQL", "Redis"}), "Go, PostgreSQL, and Redis"},
		{ListWithOptions(five, ListOptions{Limit: 3}), "Go, Redis, Kafka, and 2 more"},
		{ListWithOptions(two, ListOptions{Conjunction: "or"}), "Go or Redis"},
	}
	for i, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("example %d: got %q, want %q", i, tt.got, tt.want)
		}
	}
}
