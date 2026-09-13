package readable

import (
	"testing"
	"unicode/utf8"
)

func TestID(t *testing.T) {
	tests := map[uint64]string{
		0:                    "0",
		999:                  "999",
		1234:                 "1-234",
		987654321234567:      "987-654-321-234-567",
		18446744073709551615: "18-446-744-073-709-551-615",
	}
	for in, want := range tests {
		if got := ID(in); got != want {
			t.Errorf("ID(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s          string
		head, tail int
		want       string
	}{
		{"a8f91234abcd", 4, 2, "a8f9...cd"},
		{"abcdefghi", 4, 2, "abcdefghi"},
		{"abcdefghij", 4, 2, "abcd...ij"},
		{"abc", -1, -1, "abc"},
		{"абвгдеёжзий", 2, 2, "аб...ий"},
		{"", 0, 0, ""},
	}
	for _, tt := range tests {
		if got := Truncate(tt.s, tt.head, tt.tail); got != tt.want {
			t.Errorf("Truncate(%q,%d,%d) = %q, want %q", tt.s, tt.head, tt.tail, got, tt.want)
		}
	}
}

func TestHash(t *testing.T) {
	if got := Hash("a8f9c0ffee12ab12"); got != "a8f9...ab12" {
		t.Errorf("Hash = %q", got)
	}
	if got := Hash("abcdef"); got != "abcdef" {
		t.Errorf("Hash short = %q", got)
	}
}

func TestShortUUID(t *testing.T) {
	if got := ShortUUID("550e8400-e29b-41d4-a716-446655440000"); got != "550e...0000" {
		t.Errorf("ShortUUID = %q", got)
	}
}

func FuzzTruncate(f *testing.F) {
	f.Add("a8f91234abcd", 4, 2)
	f.Fuzz(func(t *testing.T, s string, h, tl int) {
		got := Truncate(s, h, tl)
		if utf8.ValidString(s) && !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		if len(got) > len(s) {
			t.Fatalf("grew: %q -> %q", s, got)
		}
	})
}

func BenchmarkID(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		ID(987654321234567)
	}
}

func BenchmarkTruncate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Truncate("550e8400-e29b-41d4-a716-446655440000", 4, 4)
	}
}
