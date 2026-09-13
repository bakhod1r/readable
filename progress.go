package readable

import (
	"math/bits"
	"strconv"
	"strings"
)

const (
	progressDefaultWidth = 20
	progressMaxWidth     = 200
	percentScale         = 100
)

// progressScale returns floor(done*scale/total) with done clamped to
// [0, total], computed in 128 bits. It returns 0 when total <= 0.
func progressScale(done, total int64, scale uint64) uint64 {
	if total <= 0 || done <= 0 {
		return 0
	}
	if done >= total {
		return scale
	}
	hi, lo := bits.Mul64(uint64(done), scale)
	q, _ := bits.Div64(hi, lo, uint64(total))
	return q
}

// Progress formats done out of total as a whole percentage, rounded down so
// "100%" appears only when done reaches total. done is clamped to
// [0, total]; a non-positive total yields "0%".
//
//	Progress(73, 100)   // "73%"
//	Progress(999, 1000) // "99%"
func Progress(done, total int64) string {
	return strconv.FormatUint(progressScale(done, total, percentScale), 10) + "%"
}

// ProgressBar renders a text progress bar of width cells followed by
// Progress. Filled cells (U+2588) are floor(width*done/total); the rest are
// U+2591. A non-positive width means 20; widths above 200 are capped.
//
//	ProgressBar(73, 100, 20) // "██████████████░░░░░░ 73%"
func ProgressBar(done, total int64, width int) string {
	if width <= 0 {
		width = progressDefaultWidth
	}
	width = min(width, progressMaxWidth)
	filled := int(progressScale(done, total, uint64(width)))
	var b strings.Builder
	b.Grow(width*len("█") + len(" 100%"))
	for i := range width {
		if i < filled {
			b.WriteString("█")
		} else {
			b.WriteString("░")
		}
	}
	var arr [len(" 100%")]byte
	pct := append(arr[:0], ' ')
	pct = strconv.AppendUint(pct, progressScale(done, total, percentScale), 10)
	b.Write(append(pct, '%'))
	return b.String()
}
