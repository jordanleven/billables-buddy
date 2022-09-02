package billablesbuddy

import (
	"math"
	"time"

	fc "billables-buddy/internal/forecastclient"

	hc "billables-buddy/internal/harvestclient"
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

type ProjectHours struct {
	Name  string
	Hours Hour
}

type HoursConsolidated struct {
	HoursAll         Hour
	HoursBillable    Hour
	HoursNonbillable Hour
}

type Hours struct {
	TodayStartTime    time.Time
	HoursConsolidated HoursConsolidated
	HoursByProject    []ProjectHours
}

type getHoursOpt struct {
	ts             time.Time
	todayStartTime time.Time
	actual         hc.TimeEntry
	expected       fc.TimeEntry
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

func getHours(opts getHoursOpt) Hour {
	return Hour{
		ExpectedTotal:    opts.expected.Total,
		ExpectedSchedule: opts.expected.Schedule,
		Actual:           opts.actual.Total,
		Expected:         getExpectedHoursFromSchedule(opts.ts, opts.todayStartTime, opts.expected.Schedule),
	}
}

func getAdjustedNonbillablesSchedule(billable fc.Schedule, nonbillable fc.Schedule, timeoff fc.Schedule) fc.Schedule {
	scheduleAdjusted := fc.Schedule{}

	for date, hours := range billable {
		scheduledNonbillables := nonbillable[date]
		scheduledTimeOff := timeoff[date]
		nonScheduledNonBillables := 0.0

		if hours > 0 || scheduledNonbillables > 0 || scheduledTimeOff > 0 {
			nonScheduledNonBillables = math.Max(workdayWorkingDurationInHours-hours, scheduledNonbillables) - scheduledTimeOff
		}

		scheduleAdjusted[date] = nonScheduledNonBillables
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

func getAdjustedHoursNonbillables(schedule fc.ExpectedHoursConsolidated) fc.TimeEntry {
	adjustedSchedule := getAdjustedNonbillablesSchedule(schedule.HoursBillable.Schedule, schedule.HoursNonbillable.Schedule, schedule.HoursTimeOff.Schedule)
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

func getHoursConsolidated(s StatisticDates, actualHours hc.TrackedHours, expectedHours fc.ExpectedHours) HoursConsolidated {
	expectedHoursConsolidated := expectedHours.HoursConsolidated
	expectedHoursNonbillables := getAdjustedHoursNonbillables(expectedHoursConsolidated)
	expectedHoursAll := getAdjustedHoursAll(expectedHoursConsolidated.HoursBillable, expectedHoursNonbillables)

	actualHoursConsolidated := actualHours.HoursConsolidated

	return HoursConsolidated{
		HoursBillable: getHours(
			getHoursOpt{
				ts:             s.CurrentTimestamp,
				todayStartTime: actualHours.TodayStartTime,
				actual:         actualHoursConsolidated.HoursBillable,
				expected:       expectedHoursConsolidated.HoursBillable,
			},
		),
		HoursNonbillable: getHours(
			getHoursOpt{
				ts:             s.CurrentTimestamp,
				todayStartTime: actualHours.TodayStartTime,
				actual:         actualHoursConsolidated.HoursNonbillable,
				expected:       expectedHoursNonbillables,
			},
		),
		HoursAll: getHours(
			getHoursOpt{
				ts:             s.CurrentTimestamp,
				todayStartTime: actualHours.TodayStartTime,
				actual:         actualHoursConsolidated.HoursAll,
				expected:       expectedHoursAll,
			},
		),
	}
}

func getHoursByProjectForecasted(s StatisticDates, actualHours hc.TrackedHours, expectedHours fc.ExpectedHours) ([]ProjectHours, []int) {
	var hoursByProjectForecast []ProjectHours
	var harvestIds []int

	for _, assignment := range expectedHours.HoursByProject {
		harvestID := assignment.HarvestID
		expectedProjectHours := assignment.Hours
		actualProjectHours := actualHours.HoursByProject[harvestID].Hours

		project := ProjectHours{
			Name: assignment.ProjectName,
			Hours: getHours(
				getHoursOpt{
					ts:             s.CurrentTimestamp,
					todayStartTime: actualHours.TodayStartTime,
					actual:         actualProjectHours,
					expected:       expectedProjectHours,
				},
			),
		}

		hoursByProjectForecast = append(hoursByProjectForecast, project)
		harvestIds = append(harvestIds, harvestID)
	}

	return hoursByProjectForecast, harvestIds
}

func arrayDoesContain(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getHoursByProjectActual(s StatisticDates, actualHours hc.TrackedHours, harvestIds []int) []ProjectHours {
	var hoursByProjectActual []ProjectHours

	for harvestId, harvestEntry := range actualHours.HoursByProject {
		if arrayDoesContain(harvestIds, harvestId) {
			continue
		}

		actualProjectHours := harvestEntry.Hours

		project := ProjectHours{
			Name: harvestEntry.ProjectName,
			Hours: getHours(
				getHoursOpt{
					ts:             s.CurrentTimestamp,
					todayStartTime: actualHours.TodayStartTime,
					actual:         actualProjectHours,
				},
			),
		}
		hoursByProjectActual = append(hoursByProjectActual, project)
	}

	return hoursByProjectActual
}

func getHoursByProject(s StatisticDates, actualHours hc.TrackedHours, expectedHours fc.ExpectedHours) []ProjectHours {

	hoursByProjectForecasted, harvestIds := getHoursByProjectForecasted(s, actualHours, expectedHours)
	hoursByProjectActual := getHoursByProjectActual(s, actualHours, harvestIds)

	hoursByProject := append(hoursByProjectForecasted, hoursByProjectActual...)
	return hoursByProject
}

func getActualAndExpectedHours(a GetHoursStatisticsArgs, s StatisticDates) Hours {
	actualHours := getCurrentWeeklyTrackedHours(a, s.WorkweekBegin, s.WorkweekEnd)
	expectedHours := getCurrentWeeklyExpectedHours(a, s.WorkweekBegin, s.WorkweekEnd)

	hoursConsolidated := getHoursConsolidated(s, actualHours, expectedHours)
	hoursByProject := getHoursByProject(s, actualHours, expectedHours)

	return Hours{
		TodayStartTime:    actualHours.TodayStartTime,
		HoursConsolidated: hoursConsolidated,
		HoursByProject:    hoursByProject,
	}
}
