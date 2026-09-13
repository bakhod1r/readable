package readable

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Count formats n, a space and the singular or regular plural form of the
// noun, without digit grouping. The singular is used when |n| == 1.
//
//	Count(1, "file")    // "1 file"
//	Count(0, "file")    // "0 files"
//	Count(1500, "user") // "1500 users"
//
// See Plural for the pluralisation rules.
func Count(n int64, singular string) string {
	var arr [64]byte
	buf := strconv.AppendInt(arr[:0], n, 10)
	buf = append(buf, ' ')
	if n == 1 || n == -1 {
		buf = append(buf, singular...)
	} else {
		buf = pluralAppend(buf, singular)
	}
	return string(buf)
}

// CountPlural is like Count with an explicit plural form, used verbatim when
// |n| != 1.
//
//	CountPlural(2, "person", "people") // "2 people"
func CountPlural(n int64, singular, plural string) string {
	word := plural
	if n == 1 || n == -1 {
		word = singular
	}
	var arr [64]byte
	buf := strconv.AppendInt(arr[:0], n, 10)
	buf = append(buf, ' ')
	buf = append(buf, word...)
	return string(buf)
}

// Plural returns singular when |n| == 1 and its regular English plural
// otherwise. Rules, with the ending matched case-insensitively: a word ending
// in s, x, z, ch or sh gets "es"; a y preceded by a consonant becomes "ies";
// anything else gets "s". The suffix is upper case when the word contains an
// upper-case letter and no lower-case one, so "CITY" becomes "CITIES". The
// empty string stays empty. Irregular nouns are not handled; use CountPlural
// for those.
//
//	Plural(1, "user")  // "user"
//	Plural(2, "match") // "matches"
//	Plural(2, "city")  // "cities"
//	Plural(2, "BOX")   // "BOXES"
func Plural(n int64, singular string) string {
	if n == 1 || n == -1 || singular == "" {
		return singular
	}
	var arr [64]byte
	return string(pluralAppend(arr[:0], singular))
}

// pluralAppend appends the regular plural of s to buf.
func pluralAppend(buf []byte, s string) []byte {
	if s == "" {
		return buf
	}
	upper := pluralIsUpper(s)
	suffix := func(lower, up string) string {
		if upper {
			return up
		}
		return lower
	}
	last := pluralLower(s[len(s)-1])
	var prev byte
	if len(s) >= 2 {
		prev = pluralLower(s[len(s)-2])
	}
	switch {
	case last == 's', last == 'x', last == 'z', (last == 'h' && (prev == 'c' || prev == 's')):
		buf = append(buf, s...)
		return append(buf, suffix("es", "ES")...)
	case last == 'y' && len(s) >= 2 && !pluralIsVowel(s[:len(s)-1]):
		buf = append(buf, s[:len(s)-1]...)
		return append(buf, suffix("ies", "IES")...)
	}
	buf = append(buf, s...)
	return append(buf, suffix("s", "S")...)
}

// pluralLower lower-cases an ASCII letter; other bytes are returned unchanged.
func pluralLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

// pluralIsVowel reports whether the last rune of s is a vowel a, e, i, o or u
// in either case.
func pluralIsVowel(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// pluralIsUpper reports whether s has an upper-case letter and no lower-case
// letter.
func pluralIsUpper(s string) bool {
	hasUpper := false
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			return false
		case unicode.IsUpper(r):
			hasUpper = true
		}
	}
	return hasUpper
}
