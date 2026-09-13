package readable

import (
	"math"
	"testing"
	"time"
)

func TestLocaleRelativeTimeWide(t *testing.T) {
	now := time.Date(2026, 9, 13, 9, 0, 0, 0, time.UTC)
	ago := func(d time.Duration) time.Time { return now.Add(-d) }
	later := now.Add
	day := 24 * time.Hour
	tests := []struct {
		l    Locale
		t    time.Time
		want string
	}{
		{UzbekCyrillic, ago(3 * time.Minute), "3 дақиқа олдин"},
		{UzbekCyrillic, later(2 * time.Hour), "2 соатдан кейин"},
		{UzbekCyrillic, now, "ҳозиргина"},
		{Kazakh, ago(3 * time.Minute), "3 минут бұрын"},
		{Kazakh, later(2 * day), "2 күннен кейін"},
		{Kazakh, later(2 * time.Hour), "2 сағаттан кейін"},
		{Kazakh, now, "дәл қазір"},
		{Turkish, ago(3 * time.Minute), "3 dakika önce"},
		{Turkish, later(2 * time.Hour), "2 saat sonra"},
		{Turkish, now, "az önce"},
		{German, ago(3 * day), "vor 3 Tagen"},
		{German, ago(day), "vor 1 Tag"},
		{German, later(60 * day), "in 2 Monaten"},
		{German, now, "gerade eben"},
		{French, ago(3 * time.Minute), "il y a 3 minutes"},
		{French, later(time.Hour), "dans 1 heure"},
		{French, ago(800 * day), "il y a 2 ans"},
		{French, now, "à l’instant"},
		{Spanish, ago(3 * time.Minute), "hace 3 minutos"},
		{Spanish, later(60 * day), "dentro de 2 meses"},
		{Spanish, ago(day), "hace 1 día"},
		{Spanish, now, "ahora mismo"},
		{Uzbek, later(14 * day), "2 haftadan keyin"},
		{Russian, later(time.Minute), "через 1 минуту"},
		{English, later(2 * time.Hour), RelativeTimeFrom(now, later(2*time.Hour))},
	}
	for _, tt := range tests {
		if got := tt.l.RelativeTimeFrom(now, tt.t); got != tt.want {
			t.Errorf("%s.RelativeTimeFrom(%v) = %q, want %q", tt.l, tt.t.Sub(now), got, tt.want)
		}
	}
}

func TestLocaleDurationLongWide(t *testing.T) {
	d := 49*time.Hour + 32*time.Minute
	tests := []struct {
		l    Locale
		d    time.Duration
		want string
	}{
		{UzbekCyrillic, d, "2 кун, 1 соат, 32 дақиқа"},
		{Kazakh, d, "2 күн, 1 сағат, 32 минут"},
		{Turkish, d, "2 gün, 1 saat, 32 dakika"},
		{German, d, "2 Tage, 1 Stunde, 32 Minuten"},
		{French, d, "2 jours, 1 heure, 32 minutes"},
		{Spanish, d, "2 días, 1 hora, 32 minutos"},
		{French, 0, "0 seconde"},
		{German, 0, "0 Sekunden"},
		{Spanish, -250 * time.Millisecond, "-250 milisegundos"},
		{Turkish, time.Millisecond, "1 milisaniye"},
	}
	for _, tt := range tests {
		if got := tt.l.DurationLong(tt.d); got != tt.want {
			t.Errorf("%s.DurationLong(%v) = %q, want %q", tt.l, tt.d, got, tt.want)
		}
	}
}

func TestLocaleList(t *testing.T) {
	three := []string{"Go", "Redis", "Kafka"}
	tests := []struct {
		l     Locale
		items []string
		want  string
	}{
		{English, three, "Go, Redis, and Kafka"},
		{Uzbek, three, "Go, Redis va Kafka"},
		{UzbekCyrillic, three, "Go, Redis ва Kafka"},
		{Russian, three, "Go, Redis и Kafka"},
		{Kazakh, three, "Go, Redis және Kafka"},
		{Turkish, three, "Go, Redis ve Kafka"},
		{German, three, "Go, Redis und Kafka"},
		{French, three, "Go, Redis et Kafka"},
		{Spanish, three, "Go, Redis y Kafka"},
		{Russian, three[:2], "Go и Redis"},
		{Russian, three[:1], "Go"},
		{Russian, nil, ""},
	}
	for _, tt := range tests {
		if got := tt.l.List(tt.items); got != tt.want {
			t.Errorf("%s.List(%q) = %q, want %q", tt.l, tt.items, got, tt.want)
		}
	}
}

func TestLocaleInt(t *testing.T) {
	tests := []struct {
		l    Locale
		n    int64
		want string
	}{
		{English, 1234567, "1,234,567"},
		{Russian, 1234567, "1 234 567"},
		{Uzbek, -1234567, "-1 234 567"},
		{French, 1234567, "1 234 567"},
		{German, 1234567, "1.234.567"},
		{Turkish, 999, "999"},
		{Spanish, -9223372036854775808, "-9.223.372.036.854.775.808"},
	}
	for _, tt := range tests {
		if got := tt.l.Int(tt.n); got != tt.want {
			t.Errorf("%s.Int(%d) = %q, want %q", tt.l, tt.n, got, tt.want)
		}
	}
}

func TestLocalePercent(t *testing.T) {
	tests := []struct {
		l     Locale
		ratio float64
		want  string
	}{
		{English, 0.1534, "15.34%"},
		{Russian, 0.1534, "15,34 %"},
		{German, 0.1534, "15,34 %"},
		{French, 0.1534, "15,34 %"},
		{Spanish, 1, "100 %"},
		{Turkish, 0.1534, "%15,34"},
		{Turkish, -0.5, "-%50"},
		{Uzbek, 0.1534, "15,34%"},
		{Kazakh, 0.1534, "15,34%"},
		{German, 12345.6, "1.234.560 %"},
		{Russian, math.NaN(), "NaN%"},
	}
	for _, tt := range tests {
		if got := tt.l.Percent(tt.ratio); got != tt.want {
			t.Errorf("%s.Percent(%v) = %q, want %q", tt.l, tt.ratio, got, tt.want)
		}
	}
}

func TestLocalePlural(t *testing.T) {
	tests := []struct {
		l     Locale
		n     uint64
		forms []string
		want  string
	}{
		{Russian, 22, []string{"файл", "файла", "файлов"}, "файла"},
		{English, 1, []string{"file", "files"}, "file"},
		{English, 0, []string{"file", "files"}, "files"},
		{French, 0, []string{"fichier", "fichiers"}, "fichier"},
		{French, 2, []string{"fichier", "fichiers"}, "fichiers"},
		{German, 2, []string{"Datei", "Dateien"}, "Dateien"},
		{Uzbek, 5, []string{"fayl"}, "fayl"},
		{Turkish, 5, []string{"dosya", "dosyalar"}, "dosya"},
		{Russian, 5, []string{"файл"}, "файл"},
		{Russian, 5, nil, ""},
	}
	for _, tt := range tests {
		if got := tt.l.Plural(tt.n, tt.forms...); got != tt.want {
			t.Errorf("%s.Plural(%d, %q) = %q, want %q", tt.l, tt.n, tt.forms, got, tt.want)
		}
	}
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		tag  string
		want Locale
		ok   bool
	}{
		{"en", English, true}, {"en-US", English, true}, {"ru_RU", Russian, true},
		{"UZ", Uzbek, true}, {"uz-Latn-UZ", Uzbek, true}, {"uz-Cyrl", UzbekCyrillic, true},
		{"uz_cyrl_UZ", UzbekCyrillic, true}, {"kk-KZ", Kazakh, true}, {"tr", Turkish, true},
		{"de-AT", German, true}, {"fr-CA", French, true}, {"es-419", Spanish, true},
		{"", English, false}, {"xx", English, false}, {"english", English, false},
	}
	for _, tt := range tests {
		got, ok := ParseLocale(tt.tag)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ParseLocale(%q) = %q, %v, want %q, %v", tt.tag, got, ok, tt.want, tt.ok)
		}
	}
}

func TestLocaleDataComplete(t *testing.T) {
	for l, d := range locales {
		for u, w := range d.units {
			if len(w.dur) == 0 {
				t.Errorf("%s unit %d has no duration forms", l, u)
			}
		}
		if d.and == "" || d.justNow == "" {
			t.Errorf("%s missing words", l)
		}
	}
}
