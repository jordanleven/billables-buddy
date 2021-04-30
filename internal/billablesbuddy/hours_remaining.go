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

func getEstimatedEndOfDay(ts time.Time, remainingHours float64) time.Time {
	switch {
	case remainingHours < 0:
		// EOD is in the past
		return time.Time{}
	case remainingHours >= workdayWorkingDurationInHours:
		// Remaining hours are longer than a single workday
		return ts.Add(workdayWorkingDurationInHours * time.Hour)
	default:
		remainingMinutes := remainingHours * 60
		return ts.Add(time.Duration(remainingMinutes) * time.Minute)
	}
}

func getHoursRemaining(ts time.Time, startTime time.Time, billables Hour, nonbillables Hour) HoursRemaining {
	date := getDateFromTime(ts)
	remainingHoursBillable := getRemainingHours(date, billables)
	remainingHoursNonbillable := getRemainingHours(date, nonbillables)
	remainingHoursAll := remainingHoursBillable + remainingHoursNonbillable
	hoursRemaining := HoursRemaining{}

	if !startTime.IsZero() {
		hoursRemaining.EstimatedEOD = getEstimatedEndOfDay(ts, remainingHoursAll)
	}

	return hoursRemaining
}
