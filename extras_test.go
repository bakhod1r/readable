package readable

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestETA(t *testing.T) {
	tests := []struct {
		done, total int64
		elapsed     time.Duration
		want        string
	}{
		{25, 100, time.Minute, "~3m left"},
		{50, 100, 90 * time.Minute, "~1h left"},
		{999, 1000, time.Millisecond, "~1s left"},
		{1, math.MaxInt64, time.Hour, "~106751d left"},
		{100, 100, time.Minute, "done"},
		{0, 100, time.Minute, ""},
		{5, 0, time.Minute, ""},
		{5, 10, 0, ""},
	}
	for _, tt := range tests {
		if got := ETA(tt.done, tt.total, tt.elapsed); got != tt.want {
			t.Errorf("ETA(%d, %d, %v) = %q; want %q", tt.done, tt.total, tt.elapsed, got, tt.want)
		}
	}
}

func TestCalendar(t *testing.T) {
	now := time.Date(2026, 9, 13, 9, 0, 0, 0, time.UTC) // Sunday
	at := func(days, h int) time.Time { return time.Date(2026, 9, 13+days, h, 30, 0, 0, time.UTC) }
	tests := []struct {
		t    time.Time
		want string
	}{
		{at(0, 14), "today 14:30"},
		{at(-1, 8), "yesterday 08:30"},
		{at(1, 8), "tomorrow 08:30"},
		{at(3, 8), "Wednesday 08:30"},
		{at(-6, 8), "Monday 08:30"},
		{at(-7, 8), "Sep 6, 08:30"},
		{at(120, 8), "Jan 11, 2027, 08:30"},
	}
	for _, tt := range tests {
		if got := Calendar(now, tt.t); got != tt.want {
			t.Errorf("Calendar(%v) = %q; want %q", tt.t, got, tt.want)
		}
	}
}

func TestTextExtras(t *testing.T) {
	cases := [][2]string{
		{Roman(2026), "MMXXVI"}, {Roman(3999), "MMMCMXCIX"}, {Roman(4), "IV"}, {Roman(0), ""}, {Roman(4000), ""},
		{SI(0.0012, "A"), "1.2 mA"}, {SI(4700, "Ω"), "4.7 kΩ"}, {SI(0, "V"), "0 V"}, {SI(-2.5e9, "W"), "-2.5 GW"},
		{SI(0.999999, "V"), "1 V"}, {SI(1e-12, "F"), "0 nF"}, {SI(3e15, "Hz"), "3000 THz"}, {SI(math.Inf(-1), "V"), "-Inf V"},
		{SI(5e-5, "s"), "50 µs"}, {SI(-1e-12, "F"), "0 nF"}, {SI(1, "m"), "1 m"},
		{Initials("john ronald tolkien"), "JR"}, {Initials("émile"), "É"}, {Initials("  "), ""},
		{Slug("Hello, World! 2026"), "hello-world-2026"}, {Slug("  --Go_Lang--  "), "go-lang"}, {Slug("Привет"), ""},
		{Ellipsis("/usr/local/share/readable/docs", 16), "/usr/loc…le/docs"}, {Ellipsis("short", 5), "short"},
		{Ellipsis("abcdef", 1), "…"}, {Ellipsis("abcdef", 0), ""}, {Ellipsis("héllo wörld", 6), "hél…ld"},
	}
	for i, c := range cases {
		if c[0] != c[1] {
			t.Errorf("case %d: got %q want %q", i, c[0], c[1])
		}
	}
}

func TestFormat(t *testing.T) {
	type inner struct{ A int }
	type user struct {
		Email    string        `readable:"mask=email"`
		Card     string        `readable:"mask=card"`
		Phone    string        `readable:"mask=phone"`
		Token    string        `readable:"mask=token"`
		IP       string        `readable:"mask=ip"`
		Secret   string        `readable:"mask=all"`
		Quota    uint64        `readable:"bytes"`
		Used     int64         `readable:"bytes"`
		Neg      int           `readable:"bytes"`
		Took     time.Duration `readable:"duration"`
		Views    int64         `readable:"number"`
		Share    float64       `readable:"percent"`
		Balance  int64         `readable:"money=USD"`
		Wrong    string        `readable:"bytes"`
		Name     string
		Nested   inner
		Hidden   string `readable:"-"`
		internal string
	}
	u := &user{"john.doe@gmail.com", "4111 1111 1111 1111", "+998 90 123-45-67", "sk_live_abc123456789", "10.1.2.3", "pw",
		1536, 2048, -1, 90 * time.Second, 1234567, 0.25, 150050, "x", "Ann", inner{7}, "h", "i"}
	want := map[string]string{
		"Email": "j***@gmail.com", "Card": "**** **** **** 1111", "Phone": "+998******567", "Token": "sk_live_****6789",
		"IP": "10.1.*.*", "Secret": "**", "Quota": "1.5 KB", "Used": "2 KB", "Neg": "-1", "Took": "1m 30s",
		"Views": "1.23M", "Share": "25%", "Balance": "1,500.50 USD", "Wrong": "x", "Name": "Ann", "Nested": "{7}",
	}
	if got := Format(&u); !reflect.DeepEqual(got, want) {
		t.Errorf("Format = %v\nwant     %v", got, want)
	}
	var nilUser *user
	if Format(nilUser) != nil || Format(42) != nil {
		t.Error("want nil for nil pointer and non-struct")
	}
	_ = u.internal
}

func TestLocale(t *testing.T) {
	now := time.Date(2026, 9, 13, 9, 0, 0, 0, time.UTC)
	ago := func(d time.Duration) time.Time { return now.Add(-d) }
	day := 24 * time.Hour
	cases := [][2]string{
		{Uzbek.RelativeTimeFrom(now, ago(3*time.Minute)), "3 daqiqa oldin"},
		{Uzbek.RelativeTimeFrom(now, now.Add(2*time.Hour)), "2 soatdan keyin"},
		{Uzbek.RelativeTimeFrom(now, now), "hozirgina"},
		{Russian.RelativeTimeFrom(now, ago(time.Minute)), "1 минуту назад"},
		{Russian.RelativeTimeFrom(now, ago(21*time.Second)), "21 секунду назад"},
		{Russian.RelativeTimeFrom(now, ago(3*time.Hour)), "3 часа назад"},
		{Russian.RelativeTimeFrom(now, now.Add(5*day)), "через 5 дней"},
		{Russian.RelativeTimeFrom(now, ago(14*day)), "2 недели назад"},
		{Russian.RelativeTimeFrom(now, ago(800*day)), "2 года назад"},
		{Russian.RelativeTimeFrom(now, now), "только что"},
		{English.RelativeTimeFrom(now, ago(time.Hour)), "1 hour ago"},
		{Locale("xx").RelativeTimeFrom(now, ago(time.Hour)), "1 hour ago"},
		{Uzbek.DurationLong(49*time.Hour + 32*time.Minute), "2 kun, 1 soat, 32 daqiqa"},
		{Russian.DurationLong(49*time.Hour + 32*time.Minute + 5*time.Second), "2 дня, 1 час, 32 минуты, 5 секунд"},
		{Russian.DurationLong(-250 * time.Millisecond), "-250 миллисекунд"},
		{Uzbek.DurationLong(1500 * time.Microsecond), "1 millisoniya"},
		{Uzbek.DurationLong(0), "0 soniya"},
		{Russian.DurationLong(10), "0 секунд"},
		{English.DurationLong(time.Hour), "1 hour"},
		{PluralRU(1, "a", "b", "c"), "a"}, {PluralRU(101, "a", "b", "c"), "a"}, {PluralRU(11, "a", "b", "c"), "c"},
		{PluralRU(22, "a", "b", "c"), "b"}, {PluralRU(12, "a", "b", "c"), "c"}, {PluralRU(0, "a", "b", "c"), "c"},
	}
	for i, c := range cases {
		if c[0] != c[1] {
			t.Errorf("case %d: got %q want %q", i, c[0], c[1])
		}
	}
	if s := Uzbek.RelativeTime(time.Now().Add(-time.Hour - time.Minute)); s != "1 soat oldin" {
		t.Errorf("RelativeTime = %q", s)
	}
}
