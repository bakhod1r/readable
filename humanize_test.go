package readable

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestHumanize(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"created_at", "Created at"},
		{"payment_status", "Payment status"},
		{"HTTP_SERVER_ERROR", "HTTP server error"},
		{"userCreatedAt", "User created at"},
		{"user-created-at", "User created at"},
		{"HTTPServer", "HTTP server"},
		{"user_id", "User ID"},
		{"userId", "User ID"},
		{"id", "ID"},
		{"api_url", "API URL"},
		{"parseJSONBody", "Parse JSON body"},
		{"  __a--b  ", "A b"},
		{"___", ""},
		{"file2_name", "File2 name"},
		{"émile_zola", "Émile zola"},
		{"straßeName", "Straße name"},
	}
	for _, tt := range tests {
		if got := Humanize(tt.in); got != tt.want {
			t.Errorf("Humanize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEnum(t *testing.T) {
	tests := []struct{ in, want string }{
		{"PAYMENT_COMPLETED", "Payment completed"},
		{"IN_PROGRESS", "In progress"},
		{"OK", "OK"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Enum(tt.in); got != tt.want {
			t.Errorf("Enum(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func FuzzHumanize(f *testing.F) {
	for _, s := range []string{"", "created_at", "HTTPServer", "userId", "\xff\xfe", "ǅ_ß"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got := Humanize(s)
		if !utf8.ValidString(got) {
			t.Fatalf("invalid UTF-8 output %q for %q", got, s)
		}
		if again := Humanize(s); again != got {
			t.Fatalf("non-deterministic: %q vs %q", got, again)
		}
	})
}

func BenchmarkHumanize(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Humanize("userCreatedAt_HTTPServer_id")
	}
}

// humanizeReference is the original allocation-heavy implementation of
// Humanize, kept as an oracle for FuzzHumanizeReference.
func humanizeReference(s string) string {
	acronyms := map[string]struct{}{
		"API": {}, "HTTP": {}, "HTTPS": {}, "URL": {}, "URI": {}, "ID": {},
		"UUID": {}, "JSON": {}, "XML": {}, "SQL": {}, "IP": {}, "TCP": {},
		"UDP": {}, "DNS": {}, "TLS": {}, "SSL": {}, "JWT": {}, "CPU": {},
		"RAM": {}, "UI": {}, "OK": {},
	}
	s = strings.ToValidUTF8(s, "�")
	var words []string
	runes := []rune(s)
	start := -1
	flush := func(end int) {
		if start >= 0 {
			words = append(words, string(runes[start:end]))
			start = -1
		}
	}
	for i, r := range runes {
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			flush(i)
			continue
		}
		if start >= 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if !unicode.IsUpper(prev) || nextLower {
				flush(i)
			}
		}
		if start < 0 {
			start = i
		}
	}
	flush(len(runes))
	var sb strings.Builder
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		up := strings.ToUpper(w)
		if _, ok := acronyms[up]; ok {
			sb.WriteString(up)
			continue
		}
		lower := strings.ToLower(w)
		if i == 0 {
			r, size := utf8.DecodeRuneInString(lower)
			sb.WriteRune(unicode.ToTitle(r))
			lower = lower[size:]
		}
		sb.WriteString(lower)
	}
	return sb.String()
}

func FuzzHumanizeReference(f *testing.F) {
	for _, s := range []string{
		"", "created_at", "HTTPServer", "userId", "\xff\xfe", "ǅ_ß", "ıd", "Key",
		"a\xff\xfeB", "�\xffX", "ǅǅa", "ȺȾ_x", "parseJSONBody", "  __a--b  ", "HTTPSuRL",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if got, want := Humanize(s), humanizeReference(s); got != want {
			t.Fatalf("Humanize(%q) = %q, reference %q", s, got, want)
		}
	})
}

func TestHumanizeAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { _ = Humanize("userCreatedAt_HTTPServer_id") }); n > 1 {
		t.Errorf("Humanize allocs = %v, want <= 1", n)
	}
}

func BenchmarkHumanizeReference(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = humanizeReference("userCreatedAt_HTTPServer_id")
	}
}
