package readable

import "testing"

func TestBool(t *testing.T) {
	if Bool(true) != "Yes" || Bool(false) != "No" {
		t.Error("Bool mismatch")
	}
	if BoolLabel(true, "On", "Off") != "On" || BoolLabel(false, "On", "Off") != "Off" {
		t.Error("BoolLabel mismatch")
	}
}
