package billablesbuddy

import (
	"time"

	fc "github.com/jordanleven/billables-buddy/internal/forecastclient"
	hc "github.com/jordanleven/billables-buddy/internal/harvestclient"
)

const (
	workdayWorkingDurationInHours = 8
	workdayBreakDurationInHours   = 1
)

type Schedule = fc.Schedule
type Hour struct {
	Actual           float64
	Expected         float64
	ExpectedSchedule Schedule
	ExpectedTotal    float64
}

type Hours struct {
	HoursAll         Hour
	HoursBillable    Hour
	HoursNonbillable Hour
	TodayStartTime   time.Time
}

func getCurrentWeeklyTrackedHours(a GetHoursStatisticsArgs, startDate time.Time, endDate time.Time) hc.TrackedHours {
	api := hc.GetHarvestAPI(a.HarvestAccountToken, a.HarvestAccountId)
	return api.GetTrackedHoursBetweenDates(startDate, endDate)
}

func getCurrentWeeklyExpectedHours(a GetHoursStatisticsArgs, startDate time.Time, endDate time.Time) fc.ExpectedHours {
	api := fc.GetForecastAPI(a.HarvestAccountToken, a.ForecastAccountID)
	return api.GetExpectedHoursBetweenDates(startDate, endDate)
}

func getExpectedHoursFromPreviousWorkday(ts time.Time, schedule fc.Schedule) float64 {
	hours := 0.0
	tsDate := getDateFromTime(ts)
	tsDateAdjusted := tsDate.Add(time.Second * -1)
	for date, schHours := range schedule {
		if !date.Before(tsDateAdjusted) {
			continue
		}
		hours += schHours
	}
	return hours
}

func getCurrentWorkdayPercentageComplete(ts time.Time, startTime time.Time) float64 {
	totalSpanningHours := workdayWorkingDurationInHours + workdayBreakDurationInHours
	totalSpanningDuration := time.Duration(totalSpanningHours)
	dateEndWorkday := startTime.Add(time.Hour * totalSpanningDuration)
	totalWorkdayMinutes := float64(dateEndWorkday.Sub(startTime))
	totalMinutesSinceStart := float64(ts.Sub(startTime))

	if startTime.IsZero() {
		return 0.0
	}

	if totalMinutesSinceStart >= totalWorkdayMinutes {
		return 1.0
	} else {
		return totalMinutesSinceStart / totalWorkdayMinutes
	}
}

func getExpectedHoursFromCurrentWorkday(ts time.Time, startTime time.Time, schedule fc.Schedule) float64 {
	percentageDayComplete := getCurrentWorkdayPercentageComplete(ts, startTime)
	tsDate := getDateFromTime(ts)
	expectedTotalHoursToday := schedule[tsDate]
	return percentageDayComplete * expectedTotalHoursToday
}

func getExpectedHoursFromSchedule(ts time.Time, todayStartTime time.Time, schedule fc.Schedule) float64 {
	hoursPreviousWorkday := getExpectedHoursFromPreviousWorkday(ts, schedule)
	hoursCurrentWorkday := getExpectedHoursFromCurrentWorkday(ts, todayStartTime, schedule)
	return hoursPreviousWorkday + hoursCurrentWorkday
}

func getHours(ts time.Time, todayStartTime time.Time, actual hc.TimeEntry, expected fc.TimeEntry) Hour {
	return Hour{
		ExpectedTotal:    expected.Total,
		ExpectedSchedule: expected.Schedule,
		Actual:           actual.Total,
		Expected:         getExpectedHoursFromSchedule(ts, todayStartTime, expected.Schedule),
	}
}

func getAdjustedNonbillablesSchedule(billable fc.Schedule, nonbillable fc.Schedule) fc.Schedule {
	scheduleAdjusted := fc.Schedule{}

	for date, hours := range billable {
		if hours == 0 {
			continue
		}
		scheduledHours := hours + nonbillable[date]
		scheduleAdjusted[date] = workdayWorkingDurationInHours - scheduledHours
	}

	return scheduleAdjusted
}

func getAdjustedAllHoursSchedule(billable fc.Schedule, nonbillable fc.Schedule) fc.Schedule {
	scheduleAdjusted := fc.Schedule{}

	for date, hours := range billable {
		scheduleAdjusted[date] = hours + nonbillable[date]
	}

	return scheduleAdjusted
}

func getTotalHoursFromSchedule(schedule fc.Schedule) float64 {
	totalHours := 0.0
	for _, hours := range schedule {
		totalHours += hours
	}

	return totalHours
}

func getAdjustedHoursNonbillables(schedule fc.ExpectedHours) fc.TimeEntry {
	adjustedSchedule := getAdjustedNonbillablesSchedule(schedule.HoursBillable.Schedule, schedule.HoursNonbillable.Schedule)

	adjustedExpectedNonbillables := schedule.HoursNonbillable
	adjustedExpectedNonbillables.Schedule = adjustedSchedule
	adjustedExpectedNonbillables.Total = getTotalHoursFromSchedule(adjustedSchedule)
	return adjustedExpectedNonbillables
}

func getAdjustedHoursAll(billables fc.TimeEntry, nonbillables fc.TimeEntry) fc.TimeEntry {
	return fc.TimeEntry{
		Total:    billables.Total + nonbillables.Total,
		Schedule: getAdjustedAllHoursSchedule(billables.Schedule, nonbillables.Schedule),
	}
}

func getActualAndExpectedHours(a GetHoursStatisticsArgs, s StatisticDates) Hours {
	actualHours := getCurrentWeeklyTrackedHours(a, s.WorkweekBegin, s.WorkweekEnd)
	expectedHours := getCurrentWeeklyExpectedHours(a, s.WorkweekBegin, s.WorkweekEnd)
	expectedHoursNonbillables := getAdjustedHoursNonbillables(expectedHours)
	expectedHoursAll := getAdjustedHoursAll(expectedHours.HoursBillable, expectedHoursNonbillables)

	return Hours{
		HoursBillable:    getHours(s.CurrentTimestamp, actualHours.TodayStartTime, actualHours.Hours.HoursBillable, expectedHours.HoursBillable),
		HoursNonbillable: getHours(s.CurrentTimestamp, actualHours.TodayStartTime, actualHours.Hours.HoursNonbillable, expectedHoursNonbillables),
		HoursAll:         getHours(s.CurrentTimestamp, actualHours.TodayStartTime, actualHours.Hours.HoursAll, expectedHoursAll),
		TodayStartTime:   actualHours.TodayStartTime,
	}
}
