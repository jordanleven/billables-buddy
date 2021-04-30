package main

import (
	"testing"
	"time"
)

func TestGetDurationAbbreviationFromUnit(t *testing.T) {
	t.Run("Returns the expected abbreviation for hours", func(t *testing.T) {
		expected := "h"
		actual := getDurationAbbreviationFromUnit("Hour")
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the expected abbreviation for minutes", func(t *testing.T) {
		expected := "m"
		actual := getDurationAbbreviationFromUnit("Minute")
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}

func TestGetRoundedFloat(t *testing.T) {
	t.Run("Returns rounded float to tenths place", func(t *testing.T) {
		expected := 2.0
		actual := getRoundedFloat(1.984, 1)
		if actual != expected {
			t.Errorf("received %.2f; want %f", actual, expected)
		}
	})

	t.Run("Returns rounded float to hundredths place", func(t *testing.T) {
		expected := 1.98
		actual := getRoundedFloat(1.984, 2)
		if actual != expected {
			t.Errorf("received %.2f; want %f", actual, expected)
		}
	})
}

func TestGetFormattedNumber(t *testing.T) {
	t.Run("Returns a string without trailing zeros", func(t *testing.T) {
		expected := "2"
		actual := getFormattedNumber(2.000)
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns a string with the only required trailing zeros", func(t *testing.T) {
		expected := "2.02"
		actual := getFormattedNumber(2.0200)
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}

func TestGetFormattedDuration(t *testing.T) {
	t.Run("Returns the correctly formatted duration when not rounding", func(t *testing.T) {
		expected := "2.02"
		actual := getFormattedDuration(2.0240, 2)
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correctly formatted duration when rounding", func(t *testing.T) {
		expected := "2.03"
		actual := getFormattedDuration(2.0290, 2)
		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}

func TestGetFormattedHour(t *testing.T) {
	t.Run("Returns the correctly formatted duration when not rounding", func(t *testing.T) {
		expected := "2.02"
		actual, _ := getFormattedHour(2.0240)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correctly formatted duration when rounding", func(t *testing.T) {
		expected := "2.03"
		actual, _ := getFormattedHour(2.0250)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correct unit when not rounding", func(t *testing.T) {
		expected := "Minute"
		_, actual := getFormattedHour(0.490)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correct unit when rounding from minutes to hours", func(t *testing.T) {
		expected := "Hour"
		_, actual := getFormattedHour(0.9950)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}

func TestGetFormattedPercentageFromFloat(t *testing.T) {
	t.Run("Returns the correctly formatted float when not rounding", func(t *testing.T) {
		expected := "1.2"
		actual := getFormattedPercentageFromFloat(0.012)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correctly formatted float when rounding to the tenths place", func(t *testing.T) {
		expected := "1.3"
		actual := getFormattedPercentageFromFloat(0.0125)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correctly formatted float when rounding to whole number", func(t *testing.T) {
		expected := "2"
		actual := getFormattedPercentageFromFloat(0.01984)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}

func TestGetFormattedTime(t *testing.T) {
	t.Run("Returns the correctly formatted time", func(t *testing.T) {
		ts := time.Date(1984, 01, 24, 8, 0, 0, 0, time.UTC)
		expected := "8:00 AM"
		actual := getFormattedTime(ts)

		if actual != expected {
			t.Errorf("received %s; want %s", actual, expected)
		}
	})
}
