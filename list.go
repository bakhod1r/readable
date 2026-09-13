package readable

import (
	"strconv"
	"strings"
)

// ListOptions configures ListWithOptions.
type ListOptions struct {
	// Limit is the maximum number of items shown; the rest are summarised
	// as "N more". Zero or negative means unlimited.
	Limit int
	// Conjunction joins the last item. Empty means "and".
	Conjunction string
	// Empty is returned for an empty list. Default "".
	Empty string
}

// List joins items as an English list using the Oxford comma. An empty
// list returns "".
//
//	List([]string{"Go"})                        // "Go"
//	List([]string{"Go", "PostgreSQL"})          // "Go and PostgreSQL"
//	List([]string{"Go", "PostgreSQL", "Redis"}) // "Go, PostgreSQL, and Redis"
func List(items []string) string {
	return ListWithOptions(items, ListOptions{})
}

// ListWithOptions is like List with explicit options. When Limit is positive
// and smaller than len(items), only the first Limit items are shown and the
// remainder is summarised as a final "N more" element, which takes part in
// the Oxford-comma layout like any other item. Two elements are joined as
// "A conj B" without a comma. The result is built with at most one
// allocation; a single item is returned as is.
//
//	ListWithOptions(five, ListOptions{Limit: 3})         // "Go, Redis, Kafka, and 2 more"
//	ListWithOptions(two, ListOptions{Conjunction: "or"}) // "Go or Redis"
func ListWithOptions(items []string, opts ListOptions) string {
	if len(items) == 0 {
		return opts.Empty
	}
	conj := opts.Conjunction
	if conj == "" {
		conj = "and"
	}
	hidden := 0
	if opts.Limit > 0 && opts.Limit < len(items) {
		hidden = len(items) - opts.Limit
		items = items[:opts.Limit]
	}
	if len(items) == 1 && hidden == 0 {
		return items[0]
	}

	var numArr [20]byte
	var more []byte
	if hidden > 0 {
		more = strconv.AppendInt(numArr[:0], int64(hidden), 10)
		more = append(more, " more"...)
	}
	n := len(items)
	if more != nil {
		n++
	}

	// Layout: n-1 elements each followed by ", " (or a single " " when
	// n == 2), then conj, " " and the last element.
	sep := ", "
	if n == 2 {
		sep = " "
	}
	size := len(more) + len(conj) + (n-1)*len(sep) + 1
	for _, s := range items {
		size += len(s)
	}
	var sb strings.Builder
	sb.Grow(size)
	for _, s := range items[:n-1] {
		sb.WriteString(s)
		sb.WriteString(sep)
	}
	sb.WriteString(conj)
	sb.WriteByte(' ')
	if more != nil {
		sb.Write(more)
	} else {
		sb.WriteString(items[n-1])
	}
	return sb.String()
}
