package readable

import (
	"strconv"
	"unicode/utf8"
)

// ID formats n in groups of three digits separated by '-', for reading
// long numeric identifiers aloud.
//
//	ID(987654321234567) // "987-654-321-234-567"
//	ID(1234)            // "1-234"
//	ID(0)               // "0"
func ID(n uint64) string {
	s := strconv.FormatUint(n, 10)
	if len(s) <= 3 {
		return s
	}
	out := make([]byte, 0, len(s)+(len(s)-1)/3)
	lead := len(s) % 3
	if lead == 0 {
		lead = 3
	}
	out = append(out, s[:lead]...)
	for i := lead; i < len(s); i += 3 {
		out = append(out, '-')
		out = append(out, s[i:i+3]...)
	}
	return string(out)
}

// Truncate shortens s to head leading and tail trailing runes joined by
// "...". Negative counts are treated as 0. If the result would not be shorter
// than s (head+tail+3 >= rune count), s is returned unchanged.
//
// Truncate is for readability, not secrecy: it reveals head+tail runes.
//
//	Truncate("a8f91234abcd", 4, 2) // "a8f9...cd"
func Truncate(s string, head, tail int) string {
	head, tail = max(head, 0), max(tail, 0)
	n := utf8.RuneCountInString(s)
	if head >= n || tail >= n || head+tail+3 >= n {
		return s
	}
	start := 0
	for range head {
		_, size := utf8.DecodeRuneInString(s[start:])
		start += size
	}
	end := len(s)
	for range tail {
		_, size := utf8.DecodeLastRuneInString(s[:end])
		end -= size
	}
	return s[:start] + "..." + s[end:]
}

// Hash shortens a hex digest or commit hash to its first and last 4 runes.
//
//	Hash("a8f9c0ffee12ab12") // "a8f9...ab12"
func Hash(h string) string {
	return Truncate(h, 4, 4)
}

// ShortUUID shortens a UUID to its first and last 4 runes. The input is not
// validated; strings of 11 runes or fewer are returned unchanged.
//
//	ShortUUID("550e8400-e29b-41d4-a716-446655440000") // "550e...0000"
func ShortUUID(u string) string {
	return Truncate(u, 4, 4)
}
