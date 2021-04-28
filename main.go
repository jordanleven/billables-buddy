package main

import (
	"fmt"

	bb "github.com/jordanleven/billables-buddy/internal/billablesbuddy"
)

const (
	emojiStatusOk             = ":white_check_mark:"
	emojiStatusBehind         = ":x:"
	emojiStatusOver           = ":stop_sign:"
	emojiStatusUnknown        = ":warning:"
	nonUserInterableMenuColor = "#626366"
)

type HoursStatistics = bb.HoursStatistics

func getHoursStatistics() HoursStatistics {
	args := bb.GetHoursStatisticsArgs{
		HarvestAccountId:    harvestAccountId,
		HarvestAccountToken: harvestAccountToken,
	}
	return bb.GetTrackedHoursStatistics(args)
}

func printMenuTitle() {
	fmt.Println("Billables Buddy")
	fmt.Println("---")
}

func printHoursStatistic(title string, actual float64) {
	actualF, actualU := getFormattedHour(actual)
	actualAbbv := getDurationAbbreviationFromUnit(actualU)

	fmt.Println(title)
	fmt.Println("--Actual: " + actualF + actualAbbv + " | color=" + nonUserInterableMenuColor)
}

func printHourStatistics(s HoursStatistics) {
	billableActual := s.HoursBillable.HoursActual
	nonBillableActual := s.HoursNonbillable.HoursActual
	totalHoursActual := s.HoursTotal.HoursActual

	printHoursStatistic("Total Hours", totalHoursActual)
	printHoursStatistic("Billable Hours", billableActual)
	printHoursStatistic("Non-billable Hours", nonBillableActual)
}

func printCurrentBillables(s HoursStatistics) {
	printMenuTitle()
	printHourStatistics(s)

}

func main() {
	s := getHoursStatistics()
	printCurrentBillables(s)
}
