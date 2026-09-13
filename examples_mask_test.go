package readable_test

import (
	"fmt"

	"github.com/bakhod1r/readable"
)

func ExampleMask() {
	fmt.Println(readable.Mask("1234567890", 2, 2))
	// Output: 12******90
}

func ExampleMaskEmail() {
	fmt.Println(readable.MaskEmail("john.doe@gmail.com"))
	fmt.Println(readable.MaskEmail("jo@gmail.com"))
	// Output:
	// j***@gmail.com
	// ***@gmail.com
}

func ExampleMaskPhone() {
	fmt.Println(readable.MaskPhone("+998 90 123-45-67"))
	fmt.Println(readable.MaskPhone("(555) 123-4567"))
	// Output:
	// +998******567
	// 555*****67
}

func ExampleMaskCard() {
	fmt.Println(readable.MaskCard("8600 1234 1234 5678"))
	// Output: **** **** **** 5678
}

func ExampleMaskToken() {
	fmt.Println(readable.MaskToken("sk_live_abc123456789"))
	fmt.Println(readable.MaskToken("Ab3_x9kLmnopqrstuv"))
	fmt.Println(readable.MaskToken("sk_live_1234"))
	// Output:
	// sk_live_****6789
	// ****stuv
	// ****
}

func ExampleMaskIP() {
	fmt.Println(readable.MaskIP("192.168.1.42"))
	fmt.Println(readable.MaskIP("2001:db8::1"))
	// Output:
	// 192.168.*.*
	// 2001:db8:*
}

func ExampleID() {
	fmt.Println(readable.ID(987654321234567))
	// Output: 987-654-321-234-567
}

func ExampleTruncate() {
	fmt.Println(readable.Truncate("a8f91234abcd", 4, 2))
	// Output: a8f9...cd
}

func ExampleHash() {
	fmt.Println(readable.Hash("a8f9c0ffee12ab12"))
	// Output: a8f9...ab12
}

func ExampleShortUUID() {
	fmt.Println(readable.ShortUUID("550e8400-e29b-41d4-a716-446655440000"))
	// Output: 550e...0000
}
