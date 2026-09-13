package readable

var numberWordUnits = []unit{
	{1, ""},
	{1e3, " thousand"},
	{1e6, " million"},
	{1e9, " billion"},
	{1e12, " trillion"},
}

// NumberWords is like Number but spells the scale as a short-scale English
// word: thousand, million, billion, trillion. At most DefaultPrecision
// fraction digits are kept, rounded half up and trimmed.
//
//	NumberWords(999)     // "999"
//	NumberWords(12500)   // "12.5 thousand"
//	NumberWords(1234567) // "1.23 million"
func NumberWords(n int64) string {
	neg, abs := absInt64(n)
	return formatUnits(neg, abs, numberWordUnits, DefaultPrecision, false, ".", "")
}
