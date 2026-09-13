package readable

var bandwidthUnits = []unit{
	{1, "bps"},
	{1e3, "Kbps"},
	{1e6, "Mbps"},
	{1e9, "Gbps"},
	{1e12, "Tbps"},
}

// Throughput formats a data rate in bytes per second like Bytes followed by
// "/s".
//
//	Throughput(125000000) // "119.21 MB/s"
func Throughput(bytesPerSecond uint64) string {
	return unitsPerSecond(bytesPerSecond, binaryUnits)
}

// ThroughputIEC is like Throughput with IEC units (KiB, MiB, ...).
//
//	ThroughputIEC(125000000) // "119.21 MiB/s"
func ThroughputIEC(bytesPerSecond uint64) string {
	return unitsPerSecond(bytesPerSecond, iecUnits)
}

func unitsPerSecond(n uint64, units []unit) string {
	var arr [64]byte
	buf := appendUnits(arr[:0], false, n, units, DefaultPrecision, false, ".", " ")
	return string(append(buf, "/s"...))
}

// Bandwidth formats a network rate in bits per second with 1000-based units
// bps, Kbps, Mbps, Gbps and Tbps and up to DefaultPrecision trimmed fraction
// digits, rounded half up. Rates of a quadrillion bps or more stay in Tbps.
//
//	Bandwidth(125000000) // "125 Mbps"
func Bandwidth(bitsPerSecond uint64) string {
	return formatUnits(false, bitsPerSecond, bandwidthUnits, DefaultPrecision, false, ".", " ")
}

// RequestRate formats a request rate like NumberFloat followed by " req/s".
//
//	RequestRate(15234) // "15.23K req/s"
func RequestRate(perSecond float64) string {
	return perUnit(perSecond, "req", "/s")
}

// PerMinute formats count like NumberFloat followed by an optional label and
// "/min".
//
//	PerMinute(1200, "req") // "1.2K req/min"
//	PerMinute(1200, "")    // "1.2K/min"
func PerMinute(count float64, label string) string {
	return perUnit(count, label, "/min")
}
