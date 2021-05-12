package billablesbuddy

import (
	"time"
)

type HoursRemaining struct {
	EstimatedEOD time.Time
}

func getTotalExpectedHoursByEndOfDayFromSchedule(date time.Time, s Schedule) float64 {
	dateAdj := date.Add(1 * time.Second)

	hours := 0.0
	for schDate, schHours := range s {
		if schDate.Before(dateAdj) {
			hours += schHours
		}
	}

	return hours
}

func getRemainingHours(date time.Time, h Hour) float64 {
	eodRemainingHours := getTotalExpectedHoursByEndOfDayFromSchedule(date, h.ExpectedSchedule)
	return eodRemainingHours - h.Actual
}

func getTotalExpectedHoursByEndOfDay(date time.Time, billables Hour, nonbillables Hour) float64 {
	eodRemainingHoursBillables := getTotalExpectedHoursByEndOfDayFromSchedule(date, billables.ExpectedSchedule)
	eodRemainingHoursNonbillables := getTotalExpectedHoursByEndOfDayFromSchedule(date, nonbillables.ExpectedSchedule)
	return eodRemainingHoursBillables + eodRemainingHoursNonbillables
}

func getEstimatedEndOfDay(ts time.Time, billables Hour, nonbillables Hour) time.Time {
	date := getDateFromTime(ts)
	remainingHoursBillable := getRemainingHours(date, billables)
	remainingHoursNonbillable := getRemainingHours(date, nonbillables)

	// Cap our remaining nonbillable hours and don't allow to go negative. If nonbillables
	// go negative, it means the users schedule could possibly be shortend based on a lot of
	// actual nonbillable - even if billables are low
	if remainingHoursNonbillable < 0 {
		remainingHoursNonbillable = 0
	}

	remainingHoursAll := remainingHoursBillable + remainingHoursNonbillable
	actualHoursAll := billables.Actual + nonbillables.Actual
	eodExpectedHours := getTotalExpectedHoursByEndOfDay(date, billables, nonbillables)

	switch {
	case remainingHoursAll < 0, actualHoursAll > eodExpectedHours:
		// EOD is in the past
		return time.Time{}
	case remainingHoursAll >= workdayWorkingDurationInHours:
		// Remaining hours are longer than a single workday
		return ts.Add(workdayWorkingDurationInHours * time.Hour)
	default:
		remainingMinutes := remainingHoursAll * 60
		return ts.Add(time.Duration(remainingMinutes) * time.Minute)
	}
}

func getHoursRemaining(ts time.Time, startTime time.Time, billables Hour, nonbillables Hour) HoursRemaining {
	hoursRemaining := HoursRemaining{}
	// Only get the estimated EOD if the user has started work today
	if !startTime.IsZero() {
		hoursRemaining.EstimatedEOD = getEstimatedEndOfDay(ts, billables, nonbillables)
	}

	return hoursRemaining
}
