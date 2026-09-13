package readable_test

import (
	"fmt"

	"github.com/bakhod1r/readable"
)

func ExampleList() {
	fmt.Println(readable.List([]string{"Go", "PostgreSQL", "Redis"}))
	// Output: Go, PostgreSQL, and Redis
}

func ExampleListWithOptions() {
	items := []string{"Go", "Redis", "Kafka", "NATS", "Postgres"}
	fmt.Println(readable.ListWithOptions(items, readable.ListOptions{Limit: 3}))
	fmt.Println(readable.ListWithOptions(items[:2], readable.ListOptions{Conjunction: "or"}))
	fmt.Println(readable.ListWithOptions(nil, readable.ListOptions{Empty: "none"}))
	// Output:
	// Go, Redis, Kafka, and 2 more
	// Go or Redis
	// none
}

func ExampleCount() {
	fmt.Println(readable.Count(1, "file"))
	fmt.Println(readable.Count(1500, "user"))
	fmt.Println(readable.Count(3, "match"))
	// Output:
	// 1 file
	// 1500 users
	// 3 matches
}

func ExampleCountPlural() {
	fmt.Println(readable.CountPlural(2, "person", "people"))
	// Output: 2 people
}

func ExamplePlural() {
	fmt.Println(readable.Plural(2, "city"))
	// Output: cities
}

func ExampleHumanize() {
	fmt.Println(readable.Humanize("created_at"))
	fmt.Println(readable.Humanize("HTTPServer"))
	fmt.Println(readable.Humanize("user_id"))
	// Output:
	// Created at
	// HTTP server
	// User ID
}

func ExampleEnum() {
	fmt.Println(readable.Enum("IN_PROGRESS"))
	// Output: In progress
}

func ExampleBool() {
	fmt.Println(readable.Bool(true))
	// Output: Yes
}

func ExampleBoolLabel() {
	fmt.Println(readable.BoolLabel(false, "Enabled", "Disabled"))
	// Output: Disabled
}
