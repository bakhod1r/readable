package readable

var binaryUnits = []unit{
	{1, "B"},
	{1 << 10, "KB"},
	{1 << 20, "MB"},
	{1 << 30, "GB"},
	{1 << 40, "TB"},
	{1 << 50, "PB"},
	{1 << 60, "EB"},
}

var iecUnits = []unit{
	{1, "B"},
	{1 << 10, "KiB"},
	{1 << 20, "MiB"},
	{1 << 30, "GiB"},
	{1 << 40, "TiB"},
	{1 << 50, "PiB"},
	{1 << 60, "EiB"},
}

var siUnits = []unit{
	{1, "B"},
	{1e3, "kB"},
	{1e6, "MB"},
	{1e9, "GB"},
	{1e12, "TB"},
	{1e15, "PB"},
	{1e18, "EB"},
}

// Bytes formats n bytes with binary (1024-based) multiples labelled with the
// conventional KB, MB, GB, TB, PB and EB, with up to DefaultPrecision
// fraction digits, rounded half up and trimmed of trailing zeros.
//
//	Bytes(512)     // "512 B"
//	Bytes(1536)    // "1.5 KB"
//	Bytes(1048576) // "1 MB"
//
// Use BytesIEC for unambiguous KiB/MiB labels and BytesSI for 1000-based units.
func Bytes(n uint64) string {
	return formatUnits(false, n, binaryUnits, DefaultPrecision, false, ".", " ")
}

// BytesIEC formats n bytes with 1024-based IEC units: B, KiB, MiB, GiB, TiB,
// PiB, EiB.
//
//	BytesIEC(1024) // "1 KiB"
func BytesIEC(n uint64) string {
	return formatUnits(false, n, iecUnits, DefaultPrecision, false, ".", " ")
}

// BytesSI formats n bytes with 1000-based SI units: B, kB, MB, GB, TB, PB, EB.
//
//	BytesSI(1000) // "1 kB"
func BytesSI(n uint64) string {
	return formatUnits(false, n, siUnits, DefaultPrecision, false, ".", " ")
}

// FileSize formats a file size (as returned by os.FileInfo.Size) for display
// to end users: binary units like Bytes, but at most one fraction digit.
// Negative sizes keep their sign.
//
//	FileSize(5242880) // "5 MB"
//	FileSize(1590000) // "1.5 MB"
func FileSize(n int64) string {
	neg, abs := absInt64(n)
	return formatUnits(neg, abs, binaryUnits, 1, false, ".", " ")
}
