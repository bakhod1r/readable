package readable

// Bool returns "Yes" for true and "No" for false.
//
//	Bool(true) // "Yes"
func Bool(b bool) string {
	return BoolLabel(b, "Yes", "No")
}

// BoolLabel returns trueLabel if b is true, otherwise falseLabel.
//
//	BoolLabel(false, "Enabled", "Disabled") // "Disabled"
func BoolLabel(b bool, trueLabel, falseLabel string) string {
	if b {
		return trueLabel
	}
	return falseLabel
}
