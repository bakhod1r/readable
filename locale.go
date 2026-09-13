package readable

import (
	"strconv"
	"strings"
	"time"
)

// Locale selects the language of the localised formatters. Unknown locales
// fall back to English.
type Locale string

// Supported locales.
const (
	English Locale = "en"
	Uzbek   Locale = "uz"
	Russian Locale = "ru"
)

// RelativeTime is RelativeTime in locale l. It calls time.Now.
func (l Locale) RelativeTime(t time.Time) string {
	return l.RelativeTimeFrom(time.Now(), t)
}

// RelativeTimeFrom is RelativeTimeFrom in locale l, with the same buckets.
//
//	Uzbek.RelativeTimeFrom(now, now.Add(-3*time.Minute))  // "3 daqiqa oldin"
//	Uzbek.RelativeTimeFrom(now, now.Add(2*time.Hour))     // "2 soatdan keyin"
//	Russian.RelativeTimeFrom(now, now.Add(-1*time.Minute)) // "1 минуту назад"
//	Russian.RelativeTimeFrom(now, now.Add(5*24*time.Hour)) // "через 5 дней"
func (l Locale) RelativeTimeFrom(now, t time.Time) string {
	n, unit, future, justNow := relativeParts(now, t)
	switch l {
	case Uzbek:
		if justNow {
			return "hozirgina"
		}
		s := strconv.FormatUint(n, 10) + " " + uzUnits[unit]
		if future {
			return s + "dan keyin"
		}
		return s + " oldin"
	case Russian:
		if justNow {
			return "только что"
		}
		forms := ruUnits[unit]
		if unit == relSecond || unit == relMinute || unit == relWeek {
			forms[0] = forms[3] // accusative: "1 минуту назад"
		}
		s := strconv.FormatUint(n, 10) + " " + PluralRU(n, forms[0], forms[1], forms[2])
		if future {
			return "через " + s
		}
		return s + " назад"
	}
	return RelativeTimeFrom(now, t)
}

// DurationLong is DurationLong in locale l: days, hours, minutes and seconds
// joined by ", "; shorter durations use milliseconds.
//
//	Uzbek.DurationLong(49*time.Hour + 32*time.Minute)   // "2 kun, 1 soat, 32 daqiqa"
//	Russian.DurationLong(49*time.Hour + 32*time.Minute) // "2 дня, 1 час, 32 минуты"
func (l Locale) DurationLong(d time.Duration) string {
	if l != Uzbek && l != Russian {
		return DurationLong(d)
	}
	neg, abs := absInt64(int64(d))
	type part struct {
		size uint64
		unit int
	}
	parts := []part{{uint64(24 * time.Hour), relDay}, {uint64(time.Hour), relHour},
		{uint64(time.Minute), relMinute}, {uint64(time.Second), relSecond}}
	if abs < uint64(time.Second) {
		parts = []part{{uint64(time.Millisecond), localeMillisecond}}
	}
	var out []string
	for _, p := range parts {
		v := abs / p.size
		abs %= p.size
		if v == 0 {
			continue
		}
		if l == Uzbek {
			out = append(out, strconv.FormatUint(v, 10)+" "+uzUnits[p.unit])
		} else {
			f := ruUnits[p.unit]
			out = append(out, strconv.FormatUint(v, 10)+" "+PluralRU(v, f[0], f[1], f[2]))
		}
	}
	if len(out) == 0 {
		if l == Uzbek {
			return "0 soniya"
		}
		return "0 секунд"
	}
	s := strings.Join(out, ", ")
	if neg {
		s = "-" + s
	}
	return s
}

// localeMillisecond extends the relSecond...relYear unit indexes.
const localeMillisecond = relYear + 1

var uzUnits = [...]string{"soniya", "daqiqa", "soat", "kun", "hafta", "oy", "yil", "millisoniya"}

// ruUnits holds one, few, many and accusative-one forms.
var ruUnits = [...][4]string{
	{"секунда", "секунды", "секунд", "секунду"},
	{"минута", "минуты", "минут", "минуту"},
	{"час", "часа", "часов", "час"},
	{"день", "дня", "дней", "день"},
	{"неделя", "недели", "недель", "неделю"},
	{"месяц", "месяца", "месяцев", "месяц"},
	{"год", "года", "лет", "год"},
	{"миллисекунда", "миллисекунды", "миллисекунд", "миллисекунду"},
}

// PluralRU picks the Russian plural form for n: one (1, 21, 101...), few
// (2–4, 22–24...) or many (0, 5–20, 25–30...).
//
//	PluralRU(21, "файл", "файла", "файлов") // "файл"
//	PluralRU(3, "файл", "файла", "файлов")  // "файла"
//	PluralRU(11, "файл", "файла", "файлов") // "файлов"
func PluralRU(n uint64, one, few, many string) string {
	switch n10, n100 := n%10, n%100; {
	case n10 == 1 && n100 != 11:
		return one
	case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
		return few
	}
	return many
}
