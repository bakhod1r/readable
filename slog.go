package readable

import (
	"log/slog"
	"strings"
)

// RedactAttr is a slog.HandlerOptions.ReplaceAttr function that masks
// sensitive attributes before they are written:
//
//   - keys containing password, passwd, secret, token, apikey, authorization,
//     cookie or privatekey (case-insensitive, ignoring "_" and "-"): "****"
//   - keys email, phone, card or pan, and ip: MaskEmail, MaskPhone, MaskCard
//     and MaskIP
//   - every other string value: Redact
//
// Non-string values under sensitive keys are replaced too; other non-string
// values are kept.
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout,
//		&slog.HandlerOptions{ReplaceAttr: readable.RedactAttr}))
func RedactAttr(_ []string, a slog.Attr) slog.Attr {
	key := strings.NewReplacer("_", "", "-", "").Replace(strings.ToLower(a.Key))
	for _, s := range []string{"password", "passwd", "secret", "token", "apikey", "authorization", "cookie", "privatekey"} {
		if strings.Contains(key, s) {
			return slog.String(a.Key, maskStars4)
		}
	}
	if a.Value.Kind() != slog.KindString {
		return a
	}
	v := a.Value.String()
	switch key {
	case "email":
		v = MaskEmail(v)
	case "phone":
		v = MaskPhone(v)
	case "card", "pan":
		v = MaskCard(v)
	case "ip":
		v = MaskIP(v)
	default:
		v = Redact(v)
	}
	return slog.String(a.Key, v)
}
