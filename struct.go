package readable

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Format formats the exported fields of a struct (or pointer to struct) into
// a map keyed by field name, for admin pages, CLIs and API views. The
// `readable` struct tag picks the formatter:
//
//	bytes      Bytes            (unsigned or signed integers)
//	duration   Duration         (time.Duration or integer nanoseconds)
//	number     Number           (integers)
//	percent    Percent          (floats)
//	money=USD  Money            (integer minor units)
//	mask=K     MaskEmail, MaskPhone, MaskCard, MaskToken, MaskIP or Mask(v, 0, 0)
//	           for K = email, phone, card, token, ip or all (strings)
//	-          field omitted
//
// Untagged fields and fields whose tag does not fit their type use
// fmt.Sprint, except that an unknown mask kind masks everything. A nil
// pointer or non-struct returns nil.
//
//	type User struct {
//		Email string `readable:"mask=email"`
//		Quota uint64 `readable:"bytes"`
//	}
//	Format(User{"john.doe@gmail.com", 1536}) // map[Email:j***@gmail.com Quota:1.5 KB]
func Format(v any) map[string]string {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	out := make(map[string]string, rv.NumField())
	for i := range rv.NumField() {
		f := rv.Type().Field(i)
		tag := f.Tag.Get("readable")
		if !f.IsExported() || tag == "-" {
			continue
		}
		out[f.Name] = formatField(rv.Field(i), tag)
	}
	return out
}

func formatField(v reflect.Value, tag string) string {
	kind, arg, _ := strings.Cut(tag, "=")
	isInt := v.CanInt()
	switch {
	case kind == "bytes" && v.CanUint():
		return Bytes(v.Uint())
	case kind == "bytes" && isInt && v.Int() >= 0:
		return Bytes(uint64(v.Int()))
	case kind == "duration" && isInt:
		return Duration(time.Duration(v.Int()))
	case kind == "number" && isInt:
		return Number(v.Int())
	case kind == "percent" && v.CanFloat():
		return Percent(v.Float())
	case kind == "money" && isInt:
		return Money(v.Int(), arg)
	case kind == "mask" && v.Kind() == reflect.String:
		return maskKind(arg, v.String())
	}
	return fmt.Sprint(v.Interface())
}

func maskKind(kind, s string) string {
	switch kind {
	case "email":
		return MaskEmail(s)
	case "phone":
		return MaskPhone(s)
	case "card":
		return MaskCard(s)
	case "token":
		return MaskToken(s)
	case "ip":
		return MaskIP(s)
	}
	return Mask(s, 0, 0)
}
