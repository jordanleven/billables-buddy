package billablesbuddy

import (
	"time"
)

type EstimatedEODStatus int

const (
	EstimatedEODStatusAvailable EstimatedEODStatus = iota
	EstimatedEODStatusUnavailableDailyHoursOver
	EstimatedEODStatusUnavailableWeeklyHoursOver
)

type HoursRemaining struct {
	EstimatedEOD time.Time
	EstimatedEODStatus
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
	return eodRemainingHours - h.ActualCurrent
}

func getTotalExpectedHoursByEndOfDay(date time.Time, billables Hour, nonbillables Hour) float64 {
	eodRemainingHoursBillables := getTotalExpectedHoursByEndOfDayFromSchedule(date, billables.ExpectedSchedule)
	eodRemainingHoursNonbillables := getTotalExpectedHoursByEndOfDayFromSchedule(date, nonbillables.ExpectedSchedule)
	return eodRemainingHoursBillables + eodRemainingHoursNonbillables
}

func getEstimatedEndOfDay(ts time.Time, billables Hour, nonbillables Hour) (time.Time, EstimatedEODStatus) {
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
	actualHoursAll := billables.ActualCurrent + nonbillables.ActualCurrent
	eodExpectedHours := getTotalExpectedHoursByEndOfDay(date, billables, nonbillables)
	hoursToday := billables.ActualToday + nonbillables.ActualToday
	switch {
	// User has worked more than the specified working hours per day today
	case hoursToday >= workdayWorkingDurationInHours:
		return time.Time{}, EstimatedEODStatusUnavailableDailyHoursOver
	// User has already worked all the expected hours this week
	case actualHoursAll > eodExpectedHours:
		return time.Time{}, EstimatedEODStatusUnavailableWeeklyHoursOver
	// The remaining hours for the week are more than the regular number of working hours per day
	case remainingHoursAll >= workdayWorkingDurationInHours:
		return ts.Add(workdayWorkingDurationInHours * time.Hour), EstimatedEODStatusAvailable
	default:
		remainingMinutes := remainingHoursAll * 60
		return ts.Add(time.Duration(remainingMinutes) * time.Minute), EstimatedEODStatusAvailable
	}
}

func getHoursRemaining(ts time.Time, startTime time.Time, billables Hour, nonbillables Hour) HoursRemaining {
	hoursRemaining := HoursRemaining{}
	// Only get the estimated EOD if the user has started work today
	if !startTime.IsZero() {
		estimatedEOD, estimatedEODStatus := getEstimatedEndOfDay(ts, billables, nonbillables)
		hoursRemaining.EstimatedEOD = estimatedEOD
		hoursRemaining.EstimatedEODStatus = estimatedEODStatus
	}

	return hoursRemaining
}
