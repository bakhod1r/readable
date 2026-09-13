package readable

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Humanize turns an identifier into a sentence-case label. Words are split on
// '_', '-', Unicode white space and camelCase boundaries: an upper-case rune
// starts a new word when the previous rune is not upper case or the next rune
// is lower case, so "HTTPServer" splits as "HTTP" and "Server". The first
// word is capitalised and the rest lower-cased, except the acronyms API, CPU,
// DNS, HTTP, HTTPS, ID, IP, JSON, JWT, OK, RAM, SQL, SSL, TCP, TLS, UDP, UI,
// URI, URL, UUID and XML, which are upper-cased wherever they appear
// (matched case-insensitively). Each run of invalid UTF-8 bytes becomes a
// single U+FFFD. An input with no words returns "".
//
//	Humanize("created_at")        // "Created at"
//	Humanize("HTTP_SERVER_ERROR") // "HTTP server error"
//	Humanize("userCreatedAt")     // "User created at"
//	Humanize("user_id")           // "User ID"
func Humanize(s string) string {
	var sb strings.Builder
	start := -1 // byte offset of the current word, or -1 between words
	var prev rune
	for i := 0; i < len(s); {
		r, size := humanizeDecode(s, i)
		if r == '_' || r == '-' || unicode.IsSpace(r) {
			if start >= 0 {
				humanizeWord(&sb, s[start:i], len(s))
				start = -1
			}
			i += size
			continue
		}
		if start >= 0 && unicode.IsUpper(r) {
			next, _ := humanizeDecode(s, i+size)
			if !unicode.IsUpper(prev) || unicode.IsLower(next) {
				humanizeWord(&sb, s[start:i], len(s))
				start = i
			}
		}
		if start < 0 {
			start = i
		}
		prev = r
		i += size
	}
	if start >= 0 {
		humanizeWord(&sb, s[start:], len(s))
	}
	return sb.String()
}

// Enum turns a SCREAMING_SNAKE_CASE enum value into a sentence-case label.
// It is equivalent to Humanize.
//
//	Enum("PAYMENT_COMPLETED") // "Payment completed"
//	Enum("IN_PROGRESS")       // "In progress"
func Enum(s string) string {
	return Humanize(s)
}

// humanizeDecode decodes the rune at s[i:]. A run of invalid bytes decodes as
// one utf8.RuneError spanning the whole run, matching strings.ToValidUTF8.
// At the end of s it returns (utf8.RuneError, 0).
func humanizeDecode(s string, i int) (rune, int) {
	if i < len(s) && s[i] < utf8.RuneSelf {
		return rune(s[i]), 1
	}
	r, size := utf8.DecodeRuneInString(s[i:])
	if r != utf8.RuneError || size != 1 {
		return r, size
	}
	j := i + 1
	for j < len(s) {
		if r2, s2 := utf8.DecodeRuneInString(s[j:]); r2 != utf8.RuneError || s2 != 1 {
			break
		}
		j++
	}
	return utf8.RuneError, j - i
}

// humanizeWord appends word w to sb, preceded by a space unless it is the
// first word. inputLen sizes the single buffer allocation.
func humanizeWord(sb *strings.Builder, w string, inputLen int) {
	first := sb.Len() == 0
	if first {
		sb.Grow(inputLen + inputLen/2)
	} else {
		sb.WriteByte(' ')
	}
	var key [humanizeAcronymMaxLen]byte
	if n, ok := humanizeAcronym(&key, w); ok {
		sb.Write(key[:n])
		return
	}
	for i := 0; i < len(w); {
		if c := w[i]; c < utf8.RuneSelf && !first {
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}
			sb.WriteByte(c)
			i++
			continue
		}
		r, size := humanizeDecode(w, i)
		r = unicode.ToLower(r)
		if first && i == 0 {
			r = unicode.ToTitle(r)
		}
		sb.WriteRune(r)
		i += size
	}
}

// humanizeAcronymMaxLen is the rune length of the longest known acronym.
const humanizeAcronymMaxLen = 5

// humanizeAcronym reports whether the upper-case form of w (as produced by
// strings.ToUpper) is a known acronym. On success key[:n] holds that form.
func humanizeAcronym(key *[humanizeAcronymMaxLen]byte, w string) (n int, ok bool) {
	for _, r := range w {
		u := unicode.ToUpper(r)
		if n == len(key) || u < 'A' || u > 'Z' {
			return 0, false
		}
		key[n] = byte(u)
		n++
	}
	switch string(key[:n]) {
	case "API", "HTTP", "HTTPS", "URL", "URI", "ID", "UUID", "JSON", "XML",
		"SQL", "IP", "TCP", "UDP", "DNS", "TLS", "SSL", "JWT", "CPU", "RAM",
		"UI", "OK":
		return n, true
	}
	return 0, false
}
