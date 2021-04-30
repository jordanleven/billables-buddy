package main

import (
	"math"
	"strconv"
	"time"
)

const (
	DurationUnitHour           = "Hour"
	DurationUnitMinute         = "Minute"
	DurationAbbreviationHour   = "h"
	DurationAbbreviationMinute = "m"
)

func getDurationAbbreviationFromUnit(duration string) string {
	switch duration {
	case DurationUnitHour:
		return DurationAbbreviationHour
	case DurationUnitMinute:
		return DurationAbbreviationMinute
	}
	return ""
}

func getRoundedFloat(v float64, precision int) float64 {
	f := math.Pow10(precision)
	return math.Round(v*f) / f
}

// getFormattedNumber returns a float without trailing zeros
func getFormattedNumber(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func getFormattedDuration(v float64, precision int) string {
	vR := getRoundedFloat(v, precision)
	vF := getFormattedNumber(vR)
	return vF
}

func getFormattedHour(v float64) (duration string, unit string) {
	d := math.Abs(v)
	u := DurationUnitHour
	p := 2
	dR := getRoundedFloat(d, p)
	if dR > 0 && dR < 1 {
		d = dR * 60
		u = DurationUnitMinute
		p = 0
	}

	vF := getFormattedDuration(d, p)
	return vF, u
}

func getFormattedPercentageFromFloat(v float64) string {
	vP := v * 100
	vR := getRoundedFloat(vP, 1)
	vF := getFormattedNumber(vR)
	return vF
}

func getFormattedTime(t time.Time) string {
	return t.Format("3:04 PM")
}
