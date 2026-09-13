package readable

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMask(t *testing.T) {
	tests := []struct {
		s          string
		start, end int
		want       string
	}{
		{"1234567890", 2, 2, "12******90"},
		{"secret", 3, 3, "******"},
		{"secret", 10, 0, "******"},
		{"secret", -1, -5, "******"},
		{"", 1, 1, ""},
		{"héllo wörld", 1, 1, "h*********d"},
		{"abc", 0, 0, "***"},
	}
	for _, tt := range tests {
		if got := Mask(tt.s, tt.start, tt.end); got != tt.want {
			t.Errorf("Mask(%q,%d,%d) = %q, want %q", tt.s, tt.start, tt.end, got, tt.want)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	tests := map[string]string{
		"john.doe@gmail.com": "j***@gmail.com",
		"a@b.c":              "***@b.c",
		"ab@b.c":             "***@b.c",
		"abc@b.c":            "a***@b.c",
		"ülrich@x.de":        "ü***@x.de",
		"üü@x.de":            "***@x.de",
		"we@ird@host.com":    "w***@host.com",
		"noat":               "***",
		"@gmail.com":         "***",
		"john@":              "***",
		"":                   "***",
		"\xff@x.com":         "***",
		"\xffab@x.com":       "***",
	}
	for in, want := range tests {
		if got := MaskEmail(in); got != want {
			t.Errorf("MaskEmail(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	tests := map[string]string{
		"+998901234567":      "+998******567",
		"9012345678":         "901*****78",
		"90123456789":        "901******89",
		"123456":             "******",
		"+998 (90) 123-4567": "+998******567",
		"1234567":            "****567",
		"123456789":          "******789",
		"12345":              "*****",
		"+":                  "***",
		"":                   "***",
		"12a45678":           "***",
		"++1234567":          "***",
	}
	for in, want := range tests {
		if got := MaskPhone(in); got != want {
			t.Errorf("MaskPhone(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMaskCard(t *testing.T) {
	tests := map[string]string{
		"8600123412345678":     "**** **** **** 5678",
		"8600 1234-1234 5678":  "**** **** **** 5678",
		"378282246310005":      "**** **** *** 0005",
		"123456789012":         "**** **** 9012",
		"1234567890123456789":  "**** **** **** *** 6789",
		"12345678901":          "****",
		"12345678901234567890": "****",
		"8600x23412345678":     "****",
		"":                     "****",
	}
	for in, want := range tests {
		if got := MaskCard(in); got != want {
			t.Errorf("MaskCard(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMaskToken(t *testing.T) {
	tests := map[string]string{
		"sk_live_abc123456789":             "sk_live_****6789",
		"pk_test_ABCDEFGH1234":             "pk_test_****1234",
		"ghp_ABCDEFGHIJKL1234":             "ghp_****1234",
		"github_pat_ABCDEFGH1234":          "github_pat_****1234",
		"xoxb-123456789012":                "xoxb-****9012",
		"abcdefgh12345678":                 "****5678",
		"abcdefgh1234":                     "****1234",
		"abcdefgh123":                      "****",
		"sk_live_12345678":                 "****",
		"Ab3_x9kLmnopqrstuv":               "****stuv",
		"sk_live_short":                    "****",
		"short":                            "****",
		"":                                 "****",
		"sk_":                              "****",
		"verylongprefix_abcdefgh1234":      "****1234",
		"abcdefg_ABCDEFGHIJKL":             "****IJKL",
		"ab_cd_ef_gh_ij_ABCDEFGHIJKL":      "ab_cd_ef_gh_****IJKL",
		"ключ_абвгдежзий":                  "****жзий",
		"\xff\xfe\xfd\xfc\xfb\xfa\xf9\xf8": "****",
	}
	for in, want := range tests {
		if got := MaskToken(in); got != want {
			t.Errorf("MaskToken(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMaskIP(t *testing.T) {
	tests := map[string]string{
		"192.168.1.42":    "192.168.*.*",
		"2001:db8::1":     "2001:db8:*",
		"fe80::1%eth0":    "fe80:0:*",
		"::ffff:10.0.0.1": "10.0.*.*",
		"999.1.1.1":       "***",
		"":                "***",
		"example.com":     "***",
	}
	for in, want := range tests {
		if got := MaskIP(in); got != want {
			t.Errorf("MaskIP(%q) = %q, want %q", in, got, want)
		}
	}
}

func FuzzMask(f *testing.F) {
	f.Add("hello world", 2, 3)
	f.Add("\xff\xfe", 1, 0)
	f.Fuzz(func(t *testing.T, s string, a, b int) {
		got := Mask(s, a, b)
		if utf8.ValidString(s) && !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		if utf8.RuneCountInString(got) != utf8.RuneCountInString(s) {
			t.Fatalf("rune count changed: %q -> %q", s, got)
		}
	})
}

func FuzzMaskEmail(f *testing.F) {
	f.Add("john.doe@gmail.com")
	f.Add("@@")
	f.Fuzz(func(t *testing.T, s string) {
		got := MaskEmail(s)
		if utf8.ValidString(s) && !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		if got == maskFull {
			return
		}
		local := s[:strings.LastIndexByte(s, '@')]
		gotAt := strings.LastIndexByte(got, '@')
		if gotAt < 0 || got[gotAt:] != s[len(local):] {
			t.Fatalf("domain changed: %q -> %q", s, got)
		}
		gotLocal, ok := strings.CutSuffix(got[:gotAt], maskFull)
		if !ok {
			t.Fatalf("missing placeholder: %q -> %q", s, got)
		}
		n := utf8.RuneCountInString(local)
		if gotLocal == local || utf8.RuneCountInString(gotLocal) > 1 || !strings.HasPrefix(local, gotLocal) ||
			(gotLocal != "" && n < 3) {
			t.Fatalf("local part leak: %q -> %q", s, got)
		}
	})
}

func FuzzMaskPhone(f *testing.F) {
	f.Add("+998901234567")
	f.Fuzz(func(t *testing.T, s string) {
		got := MaskPhone(s)
		if !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		shown, total := idCountDigits(got), idCountDigits(s)
		if shown > 6 || 2*shown > total {
			t.Fatalf("too many digits: %q -> %q", s, got)
		}
		if i := strings.LastIndexByte(got, '*'); i >= 0 && len(got)-i-1 > 3 {
			t.Fatalf("more than 3 trailing digits: %q -> %q", s, got)
		}
	})
}

func FuzzMaskCard(f *testing.F) {
	f.Add("8600123412345678")
	f.Fuzz(func(t *testing.T, s string) {
		got := MaskCard(s)
		if !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		if idCountDigits(got) > 4 {
			t.Fatalf("too many digits: %q -> %q", s, got)
		}
	})
}

func FuzzMaskToken(f *testing.F) {
	f.Add("sk_live_abc123456789")
	f.Fuzz(func(t *testing.T, s string) {
		got := MaskToken(s)
		if !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
		i := strings.Index(got, maskStars4)
		if i < 0 {
			t.Fatalf("missing placeholder: %q -> %q", s, got)
		}
		prefix, tail := got[:i], got[i+len(maskStars4):]
		if prefix != tokenVendorPrefix(s) && prefix != "" {
			t.Fatalf("non-vendor prefix kept: %q -> %q", s, got)
		}
		if !strings.HasPrefix(s, prefix) || !strings.HasSuffix(s, tail) {
			t.Fatalf("output not taken from input: %q -> %q", s, got)
		}
		if utf8.RuneCountInString(tail) > 4 {
			t.Fatalf("tail leak: %q -> %q", s, got)
		}
		if tail != "" && utf8.RuneCountInString(s[len(prefix):]) < 12 {
			t.Fatalf("tail revealed from short secret: %q -> %q", s, got)
		}
	})
}

func FuzzMaskIP(f *testing.F) {
	f.Add("192.168.1.42")
	f.Add("2001:db8::1")
	f.Fuzz(func(t *testing.T, s string) {
		if got := MaskIP(s); !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 %q", got)
		}
	})
}

func idCountDigits(s string) int {
	n := 0
	for i := range len(s) {
		if s[i] >= '0' && s[i] <= '9' {
			n++
		}
	}
	return n
}

func BenchmarkMask(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		Mask("1234567890abcdef", 2, 2)
	}
}

func BenchmarkMaskEmail(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		MaskEmail("john.doe@gmail.com")
	}
}

func BenchmarkMaskCard(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		MaskCard("8600 1234 1234 5678")
	}
}

func BenchmarkMaskToken(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		MaskToken("sk_live_abc123456789")
	}
}

func BenchmarkMaskIPv6(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		MaskIP("2001:db8::1")
	}
}

func BenchmarkMaskIP(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		MaskIP("192.168.1.42")
	}
}
