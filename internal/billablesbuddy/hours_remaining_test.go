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
				Actual:           2.5,
				ExpectedSchedule: schedule,
			})
		expected := 3.5

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetEstimatedEndOfDay(t *testing.T) {
	t.Run("Returns a zero time if the remaining hours are less than zero", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		actual := getEstimatedEndOfDay(ts, -1)

		if !actual.IsZero() {
			t.Errorf("Received %s; want zero time", actual)
		}
	})

	t.Run("Returns an estimated EOD that adds the number of remaining hours if less than the workday working hours", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		actual := getEstimatedEndOfDay(ts, 5)
		expected := time.Date(1984, 1, 24, 13, 30, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns an estimated EOD that maxes out at the workday working hours", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		actual := getEstimatedEndOfDay(ts, 30)
		expected := time.Date(1984, 1, 24, 16, 30, 0, 0, time.UTC)

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
			Hour{Actual: 0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{Actual: 0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)
		expected := HoursRemaining{
			EstimatedEOD: time.Date(1984, 1, 24, 16, 30, 0, 0, time.UTC),
		}
		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns a nil time if the day hasn't started", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 30, 0, 0, time.UTC)
		actual := getHoursRemaining(
			ts,
			time.Time{},
			Hour{Actual: 0,
				ExpectedSchedule: scheduleBillables,
			},
			Hour{Actual: 0,
				ExpectedSchedule: scheduleNonbillables,
			},
		)

		if !actual.EstimatedEOD.IsZero() {
			t.Errorf("Received %s; want nil time", actual)
		}
	})
}
