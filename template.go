package readable

// FuncMap returns the package's formatters by lowercase-first name for use
// with text/template and html/template:
//
//	tmpl := template.New("").Funcs(readable.FuncMap())
//	// {{ bytes .Size }} · {{ duration .Elapsed }} · {{ money .Amount "USD" }}
//
// The map is typed map[string]any so it converts to either package's
// FuncMap without importing them. Each call returns a new map. Only
// single-result formatters are included; PercentChange and the Parse
// functions are not.
func FuncMap() map[string]any {
	return map[string]any{
		"bytes": Bytes, "bytesIEC": BytesIEC, "bytesSI": BytesSI, "fileSize": FileSize,
		"bytesRate": BytesRate, "throughput": Throughput, "throughputIEC": ThroughputIEC, "bandwidth": Bandwidth,
		"number": Number, "numberFloat": NumberFloat, "numberWithPrecision": NumberWithPrecision,
		"numberWords": NumberWords, "ordinal": Ordinal, "count": Count, "countPlural": CountPlural, "plural": Plural,
		"duration": Duration, "durationShort": DurationShort, "durationWithOptions": DurationWithOptions,
		"durationLong": DurationLong, "durationNatural": DurationNatural, "durationApprox": DurationApprox,
		"latency": Latency, "relativeTime": RelativeTime, "relativeTimeFrom": RelativeTimeFrom,
		"date": Date, "dateFrom": DateFrom, "time": Time, "timeFrom": TimeFrom,
		"timeRange": TimeRange, "timeRangeFrom": TimeRangeFrom,
		"money": Money, "moneySymbol": MoneySymbol, "moneyWithPrecision": MoneyWithPrecision,
		"moneyCompact": MoneyCompact, "moneyAccounting": MoneyAccounting, "moneyChange": MoneyChange,
		"percent": Percent, "percentWithPrecision": PercentWithPrecision,
		"progress": Progress, "progressBar": ProgressBar,
		"rate": Rate, "rateWithLabel": RateWithLabel, "perMinute": PerMinute, "requestRate": RequestRate,
		"humanize": Humanize, "enum": Enum, "list": List, "bool": Bool, "boolLabel": BoolLabel,
		"maskEmail": MaskEmail, "maskPhone": MaskPhone, "maskCard": MaskCard, "maskToken": MaskToken,
		"maskIP": MaskIP, "mask": Mask, "id": ID, "hash": Hash, "shortUUID": ShortUUID, "truncate": Truncate,
		"redact": Redact, "eta": ETA, "calendar": Calendar, "roman": Roman, "si": SI,
		"initials": Initials, "slug": Slug, "ellipsis": Ellipsis, "pluralRU": PluralRU,
	}
}
