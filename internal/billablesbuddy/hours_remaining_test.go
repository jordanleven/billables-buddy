package billablesbuddy

import (
	"testing"
	"time"
)

func TestGetTotalExpectedHoursByEndOfDayFromSchedule(t *testing.T) {
	ts := time.Date(1984, 1, 23, 00, 00, 00, 00, time.UTC)
	schedule := Schedule{
		// Monday
		ts.AddDate(0, 0, 0): 0,
		ts.AddDate(0, 0, 1): 6.2,
		ts.AddDate(0, 0, 2): 4.5,
		ts.AddDate(0, 0, 3): 2.5,
		ts.AddDate(0, 0, 4): 1,
		ts.AddDate(0, 0, 5): 0,
		// Sunday
		ts.AddDate(0, 0, 6): 0,
	}

	t.Run("Returns the correct remaining hours for the day before the week has begun", func(t *testing.T) {
		ts := time.Date(1984, 1, 22, 8, 30, 0, 0, time.UTC)

		actual := getTotalExpectedHoursByEndOfDayFromSchedule(
			ts,
			schedule,
		)
		expected := 0.0

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the correct remaining hours for the day during the week", func(t *testing.T) {
		ts := time.Date(1984, 1, 25, 8, 30, 0, 0, time.UTC)

		actual := getTotalExpectedHoursByEndOfDayFromSchedule(
			ts,
			schedule,
		)
		expected := 10.7

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetRemainingHours(t *testing.T) {
	ts := time.Date(1984, 1, 23, 00, 00, 00, 00, time.UTC)
	schedule := Schedule{
		// Monday
		ts.AddDate(0, 0, 0): 0,
		ts.AddDate(0, 0, 1): 6,
		ts.AddDate(0, 0, 2): 6,
		ts.AddDate(0, 0, 3): 6,
		ts.AddDate(0, 0, 4): 6,
		ts.AddDate(0, 0, 5): 0,
		// Sunday
		ts.AddDate(0, 0, 6): 0,
	}

	t.Run("Returns the difference in hours", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)

		actual := getRemainingHours(
			ts,
			Hour{
				ActualCurrent:    2.5,
				ExpectedSchedule: schedule,
			})
		expected := 3.5

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetEstimatedEndOfDay(t *testing.T) {
	weekStart := time.Date(1984, 1, 23, 00, 00, 00, 00, time.UTC)

	billables := Hour{}
	billables.ActualCurrent = 0
	// Sum is 24
	billables.ExpectedSchedule = Schedule{
		// Monday
		weekStart.AddDate(0, 0, 0): 6.4,
		weekStart.AddDate(0, 0, 1): 6.4,
		weekStart.AddDate(0, 0, 2): 6.4,
		weekStart.AddDate(0, 0, 3): 6.4,
		weekStart.AddDate(0, 0, 4): 6.4,
		weekStart.AddDate(0, 0, 5): 0,
		// Sunday
		weekStart.AddDate(0, 0, 6): 0,
	}

	nonbillables := Hour{}
	nonbillables.ActualCurrent = 0
	// Sum is 8
	nonbillables.ExpectedSchedule = Schedule{
		// Monday
		weekStart.AddDate(0, 0, 0): 1.6,
		weekStart.AddDate(0, 0, 1): 1.6,
		weekStart.AddDate(0, 0, 2): 1.6,
		weekStart.AddDate(0, 0, 3): 1.6,
		weekStart.AddDate(0, 0, 4): 1.6,
		weekStart.AddDate(0, 0, 5): 0,
		// Sunday
		weekStart.AddDate(0, 0, 6): 0,
	}

	t.Run("Returns a zero time if the remaining hours are less than zero", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		billables.ActualCurrent = 40
		nonbillables.ActualCurrent = 40
		actual, _ := getEstimatedEndOfDay(ts, billables, nonbillables)

		if !actual.IsZero() {
			t.Errorf("Received %s; want zero time", actual)
		}
	})

	t.Run("Returns an estimated EOD that adds the number of remaining hours if less than the workday working hours", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		billables.ActualCurrent = 6
		nonbillables.ActualCurrent = 2
		actual, _ := getEstimatedEndOfDay(ts, billables, nonbillables)
		expected := time.Date(1984, 1, 24, 16, 30, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns an estimated EOD that maxes out at the workday working hours", func(t *testing.T) {
		ts := time.Date(1984, 1, 25, 8, 30, 0, 0, time.UTC)
		billables.ActualCurrent = 0
		nonbillables.ActualCurrent = 0
		actual, _ := getEstimatedEndOfDay(ts, billables, nonbillables)
		expected := time.Date(1984, 1, 25, 16, 30, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns a shortened day if the user is ahead on billables but has not fulfilled their nonbillables", func(t *testing.T) {
		ts := time.Date(1984, 1, 27, 8, 30, 0, 0, time.UTC)
		billables.ActualCurrent = 38
		nonbillables.ActualCurrent = 0
		actual, _ := getEstimatedEndOfDay(ts, billables, nonbillables)
		expected := time.Date(1984, 1, 27, 10, 30, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Does not return a shortened day if the user is ahead on nonbillables but has not fulfilled their billables", func(t *testing.T) {
		ts := time.Date(1984, 1, 27, 8, 30, 0, 0, time.UTC)
		billables.ActualCurrent = 24
		nonbillables.ActualCurrent = 15
		actual, _ := getEstimatedEndOfDay(ts, billables, nonbillables)
		expected := time.Date(1984, 1, 27, 16, 30, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}

func TestGetHoursRemaining(t *testing.T) {
	ts := time.Date(1984, 1, 23, 00, 00, 00, 00, time.UTC)
	scheduleBillables := Schedule{
		// Monday
		ts.AddDate(0, 0, 0): 0,
		ts.AddDate(0, 0, 1): 6,
		ts.AddDate(0, 0, 2): 6,
		ts.AddDate(0, 0, 3): 6,
		ts.AddDate(0, 0, 4): 6,
		ts.AddDate(0, 0, 5): 0,
		// Sunday
		ts.AddDate(0, 0, 6): 0,
	}

	scheduleNonbillables := Schedule{
		// Monday
		ts.AddDate(0, 0, 0): 0,
		ts.AddDate(0, 0, 1): 2,
		ts.AddDate(0, 0, 2): 2,
		ts.AddDate(0, 0, 3): 2,
		ts.AddDate(0, 0, 4): 2,
		ts.AddDate(0, 0, 5): 0,
		// Sunday
		ts.AddDate(0, 0, 6): 0,
	}

	t.Run("Returns the expected hours remaining struct", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		startTime := time.Date(1984, 1, 24, 8, 0, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			startTime,
			Hour{
				ActualCurrent:    0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{
				ActualCurrent:    0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)
		expected := HoursRemaining{
			EstimatedEOD: time.Date(1984, 1, 24, 16, 30, 0, 0, time.UTC),
		}
		if actual.EstimatedEOD != expected.EstimatedEOD {
			t.Errorf("Received %s; want %s", actual.EstimatedEOD, expected.EstimatedEOD)
		}
	})

	t.Run("Returns the correct status when the estimated EOD is available", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		startTime := time.Date(1984, 1, 24, 8, 0, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			startTime,
			Hour{
				ActualCurrent:    0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{
				ActualCurrent:    0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)
		expected := HoursRemaining{
			EstimatedEODStatus: EstimatedEODStatusAvailable,
		}

		if actual.EstimatedEODStatus != expected.EstimatedEODStatus {
			t.Errorf("Received %d; want %d", actual.EstimatedEODStatus, expected.EstimatedEODStatus)
		}
	})

	t.Run("Returns the correct status when the user has worked over 8 hours today", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		startTime := time.Date(1984, 1, 24, 8, 0, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			startTime,
			Hour{
				ActualToday:      8.5,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{
				ActualToday:      0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)
		expected := HoursRemaining{
			EstimatedEODStatus: EstimatedEODStatusUnavailableDailyHoursOver,
		}

		if actual.EstimatedEODStatus != expected.EstimatedEODStatus {
			t.Errorf("Received %d; want %d", actual.EstimatedEODStatus, expected.EstimatedEODStatus)
		}
	})

	t.Run("Returns the correct status when the user has worked over their expected hours for the week", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		startTime := time.Date(1984, 1, 24, 8, 0, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			startTime,
			Hour{
				ActualCurrent:    30,
				ActualToday:      0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{
				ActualCurrent:    0,
				ActualToday:      0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)
		expected := HoursRemaining{
			EstimatedEODStatus: EstimatedEODStatusUnavailableWeeklyHoursOver,
		}

		if actual.EstimatedEODStatus != expected.EstimatedEODStatus {
			t.Errorf("Received %d; want %d", actual.EstimatedEODStatus, expected.EstimatedEODStatus)
		}
	})

	t.Run("Returns a nil time if the day hasn't started", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			time.Time{},
			Hour{ActualCurrent: 0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{ActualCurrent: 0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)

		if !actual.EstimatedEOD.IsZero() {
			t.Errorf("Received %s; want nil time", actual.EstimatedEOD)
		}
	})
}
