package main

import (
	"fmt"
	"math"
	"time"

	bb "billables-buddy/internal/billablesbuddy"

	"github.com/kyokomi/emoji/v2"
)

const (
	emojiStatusOk             = ":white_check_mark:"
	emojiStatusBehind         = ":x:"
	emojiStatusOver           = ":stop_sign:"
	emojiStatusUnknown        = ":warning:"
	harvestUrlWeek            = "https://sparkbox.harvestapp.com/time/week"
	nonUserInterableMenuColor = "#626366"
)

type HoursStatistics = bb.HoursStatistics

func getHoursStatistics() HoursStatistics {
	args := bb.GetHoursStatisticsArgs{
		ForecastAccountID:   forecastAccountId,
		HarvestAccountId:    harvestAccountId,
		HarvestAccountToken: harvestAccountToken,
	}
	return bb.GetTrackedHoursStatistics(args)
}

func printMenuSeperator() {
	fmt.Println("---")
}

func printMenuTitle(s bb.Status) {
	var em string
	switch s {
	case bb.StatusOnTrack, bb.StatusAhead:
		em = emojiStatusOk
	case bb.StatusOver:
		em = emojiStatusOver
	case bb.StatusBehind:
		em = emojiStatusBehind
	default:
		em = emojiStatusUnknown
	}

	mt := emoji.Sprintf("Billables: %s", em)

	fmt.Println(mt)
	fmt.Println("---")
}

func maybeShowEstimatedEndOfDay(h bb.HoursRemaining) {
	estEOD := h.EstimatedEOD
	if estEOD.IsZero() {
		return
	}

	printEstimatedEndOfDay(estEOD)
}

func printEstimatedEndOfDay(estEOD time.Time) {
	estEODF := getFormattedTime(estEOD)

	fmt.Println("Est. EOD: " + estEODF + " | color=" + nonUserInterableMenuColor)
}

func printHoursStatistic(title string, hours bb.HoursStatistic) {
	expectedF, expectedU := getFormattedHour(hours.HoursExpected)
	expectedAbbv := getDurationAbbreviationFromUnit(expectedU)
	actualF, actualU := getFormattedHour(hours.HoursActual)
	actualAbbv := getDurationAbbreviationFromUnit(actualU)

	fmt.Println(title)
	fmt.Println("--Expected: " + expectedF + expectedAbbv + " | color=" + nonUserInterableMenuColor)
	fmt.Println("--Actual: " + actualF + actualAbbv + " | color=" + nonUserInterableMenuColor)
}

func getHoursDifferenceQualifier(hoursDiff float64) string {
	switch h := hoursDiff; {
	case h < 0:
		return "behind"
	case h > 0:
		return "ahead"
	default:
		return "on track"
	}
}

func maybeShowCurrentHoursProgress(s bb.Status, hours bb.HoursStatistic) {
	if hours.HoursExpected <= 0 {
		return
	}

	isOver := s == bb.StatusOver
	printCurrentHoursProgress(isOver, hours)
}

func printCurrentHoursProgress(isOver bool, hours bb.HoursStatistic) {
	hoursDiff := hours.HoursActual - hours.HoursExpected
	hoursDiffPercent := (hours.HoursActual - hours.HoursExpected) / hours.HoursExpected
	hoursDiffPercentAbs := math.Abs(hoursDiffPercent)
	percentF := getFormattedPercentageFromFloat(hoursDiffPercentAbs)
	hoursF, hoursU := getFormattedHour(hoursDiff)
	hoursAbbv := getDurationAbbreviationFromUnit(hoursU)

	var percentQualifier string
	if isOver {
		percentQualifier = "over"
	} else {
		percentQualifier = getHoursDifferenceQualifier(hoursDiff)
	}

	fmt.Println(percentF + "% " + percentQualifier + " (" + hoursF + hoursAbbv + ") | refresh=true | href=" + harvestUrlWeek)
}

func printHourStatistics(s HoursStatistics) {

	printHoursStatistic("Total Hours", s.HoursAll)
	printHoursStatistic("Billable Hours", s.HoursBillable)
	printHoursStatistic("Non-billable Hours", s.HoursNonbillable)
}

func printCurrentBillables(s HoursStatistics) {
	printMenuTitle(s.Status)
	maybeShowEstimatedEndOfDay(s.HoursRemaining)
	maybeShowCurrentHoursProgress(s.Status, s.HoursBillable)
	printMenuSeperator()
	printHourStatistics(s)

}

func main() {
	s := getHoursStatistics()
	printCurrentBillables(s)
}
