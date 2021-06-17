package forecastclient

import (
	"time"

	"github.com/joefitzgerald/forecast"
)

type Assignment = forecast.Assignment

type Assignments = forecast.Assignments

type Schedule map[time.Time]float64

type TimeEntry struct {
	Total    float64
	Schedule Schedule
}

type assignmentEvaluator func(Assignment) bool

func getFormattedForecastAPIDate(ts time.Time) string {
	tsUTC := getUTCTimeFromLocalTime(ts)
	tsUTCFormatted := getFormattedDate(tsUTC)
	return tsUTCFormatted
}

func getAssignmentAllocationInHours(a Assignment) float64 {
	return float64(a.Allocation) / 3600
}

func getAssignmentDateFromString(d string) time.Time {
	t, _ := time.Parse("2006-01-02", d)
	tF := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return tF
}

func getTotalAssignmentDays(startDate time.Time, a Assignment) int {
	// We'll assume that all employees end their scheduled week on a
	// Friday and are not scheduled to work on weekends
	daysBetweenWeekStartAndEnd := int(time.Friday - startDate.Weekday())
	assignmentEndMax := startDate.AddDate(0, 0, daysBetweenWeekStartAndEnd)
	assignmentS := startDate
	assignmentE := getAssignmentDateFromString(a.EndDate)
	// If the project assignment is after the max date, then the user has
	// been scheduled for a project that spans weeks in Forecast (instead of
	// being blocked off on a week-to-week basis) so we need to cap the week
	// at the end of this current week.
	if assignmentE.After(assignmentEndMax) {
		assignmentE = assignmentEndMax
	}
	totalDays := (assignmentE.Sub(assignmentS).Hours() / 24) + 1
	return int(totalDays)
}

func getTotalAssignmentHoursFromEvaluator(startDate time.Time, a Assignments, evaluator assignmentEvaluator) float64 {
	var hours float64 = 0.0
	for _, assignment := range a {
		if !evaluator(assignment) {
			continue
		}
		assignmentHours := getAssignmentAllocationInHours(assignment)
		assignmentDays := getTotalAssignmentDays(startDate, assignment)
		totalHours := assignmentHours * float64(assignmentDays)
		hours += totalHours
	}

	return hours
}

func getScheduleFromDates(schedule Schedule, startDate time.Time, totalDays int, hoursPerDay float64) Schedule {
	for i := 0; i < totalDays; i++ {
		assignmentDate := startDate.AddDate(0, 0, i)
		// Don't add any schedules on days that are not
		// already represented in the schedule
		_, ok := schedule[assignmentDate]
		if ok {
			schedule[assignmentDate] += hoursPerDay
		}
	}
	return schedule
}

func getScheduledHoursFromEvaluator(startDate time.Time, a Assignments, evaluator assignmentEvaluator) Schedule {
	schedule := Schedule{
		startDate.AddDate(0, 0, 0): 0.0,
		startDate.AddDate(0, 0, 1): 0.0,
		startDate.AddDate(0, 0, 2): 0.0,
		startDate.AddDate(0, 0, 3): 0.0,
		startDate.AddDate(0, 0, 4): 0.0,
		startDate.AddDate(0, 0, 5): 0.0,
		startDate.AddDate(0, 0, 6): 0.0,
	}

	for _, assignment := range a {
		if !evaluator(assignment) {
			continue
		}
		assignmentStartDate := getAssignmentDateFromString(assignment.StartDate)
		if assignmentStartDate.Before(startDate) {
			assignmentStartDate = startDate
		}

		assignmentHoursPerDay := getAssignmentAllocationInHours(assignment)
		assignmentDays := getTotalAssignmentDays(assignmentStartDate, assignment)
		schedule = getScheduleFromDates(schedule, assignmentStartDate, assignmentDays, assignmentHoursPerDay)
	}

	return schedule
}

func getEvaluatedHoursFromAssignments(startDate time.Time, a Assignments, evaluator assignmentEvaluator) TimeEntry {
	return TimeEntry{
		Total:    getTotalAssignmentHoursFromEvaluator(startDate, a, evaluator),
		Schedule: getScheduledHoursFromEvaluator(startDate, a, evaluator),
	}
}

func (c *ForecastClient) getUserAssignments(startDate time.Time, endDate time.Time) Assignments {
	uid := c.getCurrentUserID()
	endDateAdjust := endDate.AddDate(0, 0, -1)
	filter := forecast.AssignmentFilter{
		PersonID:  uid.ID,
		StartDate: getFormattedForecastAPIDate(startDate),
		EndDate:   getFormattedForecastAPIDate(endDateAdjust),
		State:     "active",
	}

	userAssignments, _ := c.Client.AssignmentsWithFilter(filter)
	return userAssignments
}
