package readable

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// Locale selects the language of the localised formatters. Unknown locales
// fall back to English. Use ParseLocale to map a BCP 47 tag such as "ru-RU".
type Locale string

// Supported locales.
const (
	English       Locale = "en"
	Uzbek         Locale = "uz" // Latin script
	UzbekCyrillic Locale = "uz-Cyrl"
	Russian       Locale = "ru"
	Kazakh        Locale = "kk"
	Turkish       Locale = "tr"
	German        Locale = "de"
	French        Locale = "fr"
	Spanish       Locale = "es"
)

// ParseLocale maps a BCP 47 language tag to a supported Locale, matching
// case-insensitively and accepting '_' for '-'. Region subtags are ignored;
// "uz-Cyrl" selects UzbekCyrillic. Unsupported tags return English, false.
//
//	ParseLocale("ru-RU")   // Russian, true
//	ParseLocale("uz_Cyrl") // UzbekCyrillic, true
//	ParseLocale("xx")      // English, false
func ParseLocale(tag string) (Locale, bool) {
	tag = strings.ToLower(strings.ReplaceAll(tag, "_", "-"))
	lang, rest, _ := strings.Cut(tag, "-")
	if lang == "uz" && (rest == "cyrl" || strings.HasPrefix(rest, "cyrl-")) {
		return UzbekCyrillic, true
	}
	if l := Locale(lang); l == English || locales[l] != nil {
		return l, true
	}
	return English, false
}

// Plural category indexes into a unit's word forms.
const (
	pluralOne  = 0
	pluralFew  = 1
	pluralMany = 2
)

// pluralRule returns the plural category of n for a locale.
type pluralRule func(n uint64) int

func pluralNone(uint64) int { return pluralOne }

func pluralOneOther(n uint64) int {
	if n == 1 {
		return pluralOne
	}
	return pluralFew
}

func pluralFrench(n uint64) int {
	if n <= 1 {
		return pluralOne
	}
	return pluralFew
}

func pluralRussian(n uint64) int {
	switch n10, n100 := n%10, n%100; {
	case n10 == 1 && n100 != 11:
		return pluralOne
	case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
		return pluralFew
	}
	return pluralMany
}

// unitWords holds a time unit's forms by plural category: dur for durations
// ("2 Tage"), past for relative times ("vor 2 Tagen") and future where it
// differs again ("2 soatdan keyin"). Nil future reuses past; nil past reuses
// dur.
type unitWords struct {
	dur, past, future []string
}

// localeData is everything a non-English locale needs.
type localeData struct {
	plural pluralRule
	// units is indexed by relSecond...relYear, then localeMillisecond.
	units [localeMillisecond + 1]unitWords
	// Relative time is pastPrefix + "N unit" + pastSuffix, likewise future.
	pastPrefix, pastSuffix     string
	futurePrefix, futureSuffix string
	justNow                    string
	and                        string // list conjunction
	group, decimal             string // number separators
	percentPrefix              bool   // "%15" rather than "15%"
	percentSpace               string // between number and a trailing '%'
}

// localeMillisecond extends the relSecond...relYear unit indexes.
const localeMillisecond = relYear + 1

const (
	nbsp       = "\u00a0"
	narrowNbsp = "\u202f"
)

func words(forms ...string) unitWords { return unitWords{dur: forms} }

// suffixed builds Turkic unit words whose future form adds an ablative
// suffix ("soat" -> "soatdan").
func suffixed(word, ablative string) unitWords {
	return unitWords{dur: []string{word}, future: []string{word + ablative}}
}

var locales = map[Locale]*localeData{
	Uzbek: {
		plural: pluralNone,
		units: [...]unitWords{
			suffixed("soniya", "dan"), suffixed("daqiqa", "dan"), suffixed("soat", "dan"),
			suffixed("kun", "dan"), suffixed("hafta", "dan"), suffixed("oy", "dan"),
			suffixed("yil", "dan"), words("millisoniya"),
		},
		pastSuffix: " oldin", futureSuffix: " keyin", justNow: "hozirgina",
		and: "va", group: nbsp, decimal: ",",
	},
	UzbekCyrillic: {
		plural: pluralNone,
		units: [...]unitWords{
			suffixed("сония", "дан"), suffixed("дақиқа", "дан"), suffixed("соат", "дан"),
			suffixed("кун", "дан"), suffixed("ҳафта", "дан"), suffixed("ой", "дан"),
			suffixed("йил", "дан"), words("миллисония"),
		},
		pastSuffix: " олдин", futureSuffix: " кейин", justNow: "ҳозиргина",
		and: "ва", group: nbsp, decimal: ",",
	},
	Kazakh: {
		plural: pluralNone,
		units: [...]unitWords{
			suffixed("секунд", "тан"), suffixed("минут", "тан"), suffixed("сағат", "тан"),
			suffixed("күн", "нен"), suffixed("апта", "дан"), suffixed("ай", "дан"),
			suffixed("жыл", "дан"), words("миллисекунд"),
		},
		pastSuffix: " бұрын", futureSuffix: " кейін", justNow: "дәл қазір",
		and: "және", group: nbsp, decimal: ",",
	},
	Turkish: {
		plural: pluralNone,
		units: [...]unitWords{
			words("saniye"), words("dakika"), words("saat"), words("gün"),
			words("hafta"), words("ay"), words("yıl"), words("milisaniye"),
		},
		pastSuffix: " önce", futureSuffix: " sonra", justNow: "az önce",
		and: "ve", group: ".", decimal: ",", percentPrefix: true,
	},
	Russian: {
		plural: pluralRussian,
		units: [...]unitWords{
			{dur: []string{"секунда", "секунды", "секунд"}, past: []string{"секунду", "секунды", "секунд"}},
			{dur: []string{"минута", "минуты", "минут"}, past: []string{"минуту", "минуты", "минут"}},
			words("час", "часа", "часов"),
			words("день", "дня", "дней"),
			{dur: []string{"неделя", "недели", "недель"}, past: []string{"неделю", "недели", "недель"}},
			words("месяц", "месяца", "месяцев"),
			words("год", "года", "лет"),
			words("миллисекунда", "миллисекунды", "миллисекунд"),
		},
		pastSuffix: " назад", futurePrefix: "через ", justNow: "только что",
		and: "и", group: nbsp, decimal: ",", percentSpace: nbsp,
	},
	German: {
		plural: pluralOneOther,
		units: [...]unitWords{
			words("Sekunde", "Sekunden"), words("Minute", "Minuten"), words("Stunde", "Stunden"),
			{dur: []string{"Tag", "Tage"}, past: []string{"Tag", "Tagen"}},
			words("Woche", "Wochen"),
			{dur: []string{"Monat", "Monate"}, past: []string{"Monat", "Monaten"}},
			{dur: []string{"Jahr", "Jahre"}, past: []string{"Jahr", "Jahren"}},
			words("Millisekunde", "Millisekunden"),
		},
		pastPrefix: "vor ", futurePrefix: "in ", justNow: "gerade eben",
		and: "und", group: ".", decimal: ",", percentSpace: nbsp,
	},
	French: {
		plural: pluralFrench,
		units: [...]unitWords{
			words("seconde", "secondes"), words("minute", "minutes"), words("heure", "heures"),
			words("jour", "jours"), words("semaine", "semaines"), words("mois", "mois"),
			words("an", "ans"), words("milliseconde", "millisecondes"),
		},
		pastPrefix: "il y a ", futurePrefix: "dans ", justNow: "à l’instant",
		and: "et", group: narrowNbsp, decimal: ",", percentSpace: narrowNbsp,
	},
	Spanish: {
		plural: pluralOneOther,
		units: [...]unitWords{
			words("segundo", "segundos"), words("minuto", "minutos"), words("hora", "horas"),
			words("día", "días"), words("semana", "semanas"), words("mes", "meses"),
			words("año", "años"), words("milisegundo", "milisegundos"),
		},
		pastPrefix: "hace ", futurePrefix: "dentro de ", justNow: "ahora mismo",
		and: "y", group: ".", decimal: ",", percentSpace: nbsp,
	},
}

// pick returns the form of forms for plural category c, falling back to the
// last form when a locale lists fewer, or "" when there are none.
func pick(forms []string, c int) string {
	if len(forms) == 0 {
		return ""
	}
	return forms[min(c, len(forms)-1)]
}

// RelativeTime is RelativeTime in locale l. It calls time.Now.
func (l Locale) RelativeTime(t time.Time) string {
	return l.RelativeTimeFrom(time.Now(), t)
}

// RelativeTimeFrom is RelativeTimeFrom in locale l, with the same buckets.
//
//	Uzbek.RelativeTimeFrom(now, now.Add(-3*time.Minute))  // "3 daqiqa oldin"
//	Uzbek.RelativeTimeFrom(now, now.Add(2*time.Hour))     // "2 soatdan keyin"
//	Russian.RelativeTimeFrom(now, now.Add(-1*time.Minute)) // "1 минуту назад"
//	German.RelativeTimeFrom(now, now.Add(-72*time.Hour))   // "vor 3 Tagen"
//	French.RelativeTimeFrom(now, now.Add(time.Hour))       // "dans 1 heure"
func (l Locale) RelativeTimeFrom(now, t time.Time) string {
	d := locales[l]
	if d == nil {
		return RelativeTimeFrom(now, t)
	}
	n, unit, future, justNow := relativeParts(now, t)
	if justNow {
		return d.justNow
	}
	w := d.units[unit]
	prefix, suffix, forms := d.pastPrefix, d.pastSuffix, w.past
	if future {
		prefix, suffix = d.futurePrefix, d.futureSuffix
		if w.future != nil {
			forms = w.future
		}
	}
	if forms == nil {
		forms = w.dur
	}
	return prefix + strconv.FormatUint(n, 10) + " " + pick(forms, d.plural(n)) + suffix
}

// DurationLong is DurationLong in locale l: days, hours, minutes and seconds
// joined by ", "; shorter durations use milliseconds.
//
//	Uzbek.DurationLong(49*time.Hour + 32*time.Minute)   // "2 kun, 1 soat, 32 daqiqa"
//	Russian.DurationLong(49*time.Hour + 32*time.Minute) // "2 дня, 1 час, 32 минуты"
//	German.DurationLong(49*time.Hour + 32*time.Minute)  // "2 Tage, 1 Stunde, 32 Minuten"
func (l Locale) DurationLong(d time.Duration) string {
	data := locales[l]
	if data == nil {
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
		out = append(out, strconv.FormatUint(v, 10)+" "+pick(data.units[p.unit].dur, data.plural(v)))
	}
	if len(out) == 0 {
		return "0 " + pick(data.units[relSecond].dur, data.plural(0))
	}
	s := strings.Join(out, ", ")
	if neg {
		s = "-" + s
	}
	return s
}

// List joins items as a list in locale l: English uses List (with the Oxford
// comma); other locales use ", " and the locale's conjunction before the last
// item.
//
//	Russian.List([]string{"Go", "Redis", "Kafka"}) // "Go, Redis и Kafka"
//	Uzbek.List([]string{"Go", "Redis"})            // "Go va Redis"
func (l Locale) List(items []string) string {
	d := locales[l]
	if d == nil {
		return List(items)
	}
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	}
	last := len(items) - 1
	return strings.Join(items[:last], ", ") + " " + d.and + " " + items[last]
}

// Int formats n with the thousands separator of locale l: "," for English,
// "." for German, Spanish and Turkish, a no-break space (U+00A0) for Russian,
// Uzbek and Kazakh, and a narrow no-break space (U+202F) for French.
//
//	English.Int(1234567) // "1,234,567"
//	German.Int(1234567)  // "1.234.567"
func (l Locale) Int(n int64) string {
	neg, abs := absInt64(n)
	digits := strconv.FormatUint(abs, 10)
	if neg {
		return "-" + groupDigits(digits, l.group())
	}
	return groupDigits(digits, l.group())
}

// Percent is Percent in locale l, using its decimal and thousands separators
// and its placement of the percent sign.
//
//	Russian.Percent(0.1534) // "15,34 %" (U+00A0 before '%')
//	Turkish.Percent(0.1534) // "%15,34"
func (l Locale) Percent(ratio float64) string {
	d := locales[l]
	x := ratio * 100
	if d == nil || math.IsNaN(x) || math.IsInf(x, 0) {
		return Percent(ratio)
	}
	s := string(appendFloat(nil, x, DefaultPrecision))
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	intPart, frac, hasFrac := strings.Cut(s, ".")
	num := groupDigits(intPart, d.group)
	if hasFrac {
		num += d.decimal + frac
	}
	if d.percentPrefix {
		return sign + "%" + num
	}
	return sign + num + d.percentSpace + "%"
}

// Plural returns the form of forms matching n under the plural rules of
// locale l. forms are ordered one, other for English, German, Spanish and
// French (where 0 is also "one"); one, few, many for Russian; Uzbek, Kazakh
// and Turkish use a single form after numerals. Missing forms fall back to the
// last one given; no forms yield "".
//
//	Russian.Plural(22, "файл", "файла", "файлов") // "файла"
//	French.Plural(0, "fichier", "fichiers")        // "fichier"
func (l Locale) Plural(n uint64, forms ...string) string {
	rule := pluralOneOther
	if d := locales[l]; d != nil {
		rule = d.plural
	}
	return pick(forms, rule(n))
}

// group returns the thousands separator of l.
func (l Locale) group() string {
	if d := locales[l]; d != nil {
		return d.group
	}
	return ","
}

// groupDigits inserts sep between groups of three in a string of digits.
func groupDigits(digits, sep string) string {
	if len(digits) <= 3 {
		return digits
	}
	var b strings.Builder
	b.Grow(len(digits) + len(digits)/3*len(sep))
	for i := range len(digits) {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteString(sep)
		}
		b.WriteByte(digits[i])
	}
	return b.String()
}

// PluralRU picks the Russian plural form for n: one (1, 21, 101...), few
// (2–4, 22–24...) or many (0, 5–20, 25–30...).
//
//	PluralRU(21, "файл", "файла", "файлов") // "файл"
//	PluralRU(3, "файл", "файла", "файлов")  // "файла"
//	PluralRU(11, "файл", "файла", "файлов") // "файлов"
func PluralRU(n uint64, one, few, many string) string {
	return Russian.Plural(n, one, few, many)
}
