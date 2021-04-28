package main

import (
	"testing"
)

func TestToTenthsPlace(t *testing.T) {
	expected := 2.0
	resp := getRoundedFloat(1.984, 1)
	if resp != expected {
		t.Errorf("roundFloat received %.2f; want %f", resp, expected)
	}
}

func TestToHundredthsPlace(t *testing.T) {
	expected := 1.98
	resp := getRoundedFloat(1.984, 2)
	if resp != expected {
		t.Errorf("roundFloat received %.2f; want %f", resp, expected)
	}
}
