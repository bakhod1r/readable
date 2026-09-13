package readable

import (
	"math"
	"testing"
)

func TestPlural(t *testing.T) {
	tests := []struct {
		n    int64
		in   string
		want string
	}{
		{1, "user", "user"},
		{-1, "user", "user"},
		{0, "user", "users"},
		{2, "user", "users"},
		{2, "bus", "buses"},
		{2, "box", "boxes"},
		{2, "quiz", "quizzes"},
		{2, "FEZ", "FEZZES"},
		{2, "waltz", "waltzes"},
		{2, "topaz", "topazes"},
		{2, "oz", "ozes"},
		{2, "match", "matches"},
		{2, "dish", "dishes"},
		{2, "city", "cities"},
		{2, "day", "days"},
		{2, "y", "ys"},
		{2, "", ""},
		{math.MinInt64, "file", "files"},
		{2, "CITY", "CITIES"},
		{2, "BOX", "BOXES"},
		{2, "MATCH", "MATCHES"},
		{2, "FILE", "FILES"},
		{2, "File", "Files"},
		{2, "Box", "Boxes"},
		{2, "DAY", "DAYS"},
		{2, "BUS_ID", "BUS_IDS"},
		{2, "file2", "file2s"},
		{2, "42", "42s"},
		{2, "ÉTÉ", "ÉTÉS"},
		{2, "Ay", "Ays"},
		{2, "ÜY", "ÜIES"},
		{2, "Иy", "Иies"},
		{2, "Кey", "Кeys"},
		{2, "watCH", "watCHes"},
	}
	for _, tt := range tests {
		if got := Plural(tt.n, tt.in); got != tt.want {
			t.Errorf("Plural(%d, %q) = %q, want %q", tt.n, tt.in, got, tt.want)
		}
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		n    int64
		in   string
		want string
	}{
		{1, "file", "1 file"},
		{5, "file", "5 files"},
		{0, "file", "0 files"},
		{-1, "file", "-1 file"},
		{1500, "user", "1500 users"},
		{math.MinInt64, "file", "-9223372036854775808 files"},
		{math.MaxInt64, "file", "9223372036854775807 files"},
		{2, "CITY", "2 CITIES"},
		{3, "", "3 "},
		{2, "a-very-long-noun-that-does-not-fit-in-the-stack-buffer-at-all-ok", "2 a-very-long-noun-that-does-not-fit-in-the-stack-buffer-at-all-oks"},
	}
	for _, tt := range tests {
		if got := Count(tt.n, tt.in); got != tt.want {
			t.Errorf("Count(%d, %q) = %q, want %q", tt.n, tt.in, got, tt.want)
		}
	}
}

func TestCountPlural(t *testing.T) {
	if got := CountPlural(2, "person", "people"); got != "2 people" {
		t.Errorf("got %q", got)
	}
	if got := CountPlural(1, "person", "people"); got != "1 person" {
		t.Errorf("got %q", got)
	}
	if got := CountPlural(-3, "person", "people"); got != "-3 people" {
		t.Errorf("got %q", got)
	}
}

func TestCountAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { _ = Count(1500, "match") }); n > 1 {
		t.Errorf("Count allocs = %v, want <= 1", n)
	}
	if n := testing.AllocsPerRun(100, func() { _ = CountPlural(2, "person", "people") }); n > 1 {
		t.Errorf("CountPlural allocs = %v, want <= 1", n)
	}
}

func BenchmarkPlural(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Plural(2, "city")
	}
}

func BenchmarkCountPlural(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = CountPlural(2, "person", "people")
	}
}

func BenchmarkCount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Count(1500, "match")
	}
}

func TestPluralDocExamples(t *testing.T) {
	tests := []struct{ got, want string }{
		{Count(1, "file"), "1 file"},
		{Count(0, "file"), "0 files"},
		{Count(1500, "user"), "1500 users"},
		{CountPlural(2, "person", "people"), "2 people"},
		{Plural(1, "user"), "user"},
		{Plural(2, "match"), "matches"},
		{Plural(2, "city"), "cities"},
		{Plural(2, "BOX"), "BOXES"},
	}
	for i, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("example %d: got %q, want %q", i, tt.got, tt.want)
		}
	}
}
