package billablesbuddy

import (
	"math"
	"testing"
	"time"

	fc "github.com/jordanleven/billables-buddy/internal/forecastclient"
)

func schedulesAreEqual(actual fc.Schedule, expected fc.Schedule) (bool, float64, float64) {
	for index, actualValue := range actual {
		expectedValue := expected[index]
		if math.Abs(actualValue-expectedValue) > 1e-9 {
			return false, actualValue, expectedValue
		}
	}

	return true, 0, 0
}

func TestGetExpectedHoursFromPreviousWorkday(t *testing.T) {
	ts := time.Date(1984, 1, 35, 00, 00, 00, 00, time.UTC)
	schedule := fc.Schedule{
		// Monday
		ts.AddDate(0, 0, -2): 1.0,
		ts.AddDate(0, 0, -1): 2.0,
		ts.AddDate(0, 0, 1):  4.0,
		ts.AddDate(0, 0, 2):  6.0,
		ts.AddDate(0, 0, 3):  8.0,
		ts.AddDate(0, 0, 4):  10.0,
		// Sunday
		ts.AddDate(0, 0, 5): 7.0,
	}

	t.Run("Returns the correct number of hours in the past", func(t *testing.T) {
		actual := getExpectedHoursFromPreviousWorkday(ts, schedule)
		expected := 3.0

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetCurrentWorkdayPercentageComplete(t *testing.T) {
	ts := time.Date(1984, 1, 24, 12, 00, 00, 00, time.UTC)

	t.Run("Returns the correct workday percentage when the workday hasn't begun", func(t *testing.T) {
		var startTime time.Time

		actual := getCurrentWorkdayPercentageComplete(ts, startTime)
		expected := 0.0

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the correct workday percentage when the workday has begun", func(t *testing.T) {
		startTime := time.Date(1984, 1, 24, 9, 0, 0, 0, time.UTC)

		actual := getCurrentWorkdayPercentageComplete(ts, startTime)
		expected := 1.0 / 3

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the correct workday percentage when the workday has ended", func(t *testing.T) {
		startTime := time.Date(1984, 1, 24, 3, 0, 0, 0, time.UTC)

		actual := getCurrentWorkdayPercentageComplete(ts, startTime)
		expected := 1.0

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetAdjustedNonbillablesSchedule(t *testing.T) {
	ts := time.Date(1984, 1, 35, 00, 00, 00, 00, time.UTC)

	t.Run("Supplement nonbillables when no scheduled nonbillables account for any nonbillable time", func(t *testing.T) {
		billableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 6.4,
			ts.AddDate(0, 0, -1): 6.4,
			ts.AddDate(0, 0, 1):  6.4,
			ts.AddDate(0, 0, 2):  6.4,
			ts.AddDate(0, 0, 3):  6.4,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		nonbillableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 0,
			ts.AddDate(0, 0, -1): 0,
			ts.AddDate(0, 0, 1):  0,
			ts.AddDate(0, 0, 2):  0,
			ts.AddDate(0, 0, 3):  0,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		timeOffSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 0,
			ts.AddDate(0, 0, -1): 0,
			ts.AddDate(0, 0, 1):  0,
			ts.AddDate(0, 0, 2):  0,
			ts.AddDate(0, 0, 3):  0,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		actual := getAdjustedNonbillablesSchedule(billableSchedule, nonbillableSchedule, timeOffSchedule)
		expected := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 1.6,
			ts.AddDate(0, 0, -1): 1.6,
			ts.AddDate(0, 0, 1):  1.6,
			ts.AddDate(0, 0, 2):  1.6,
			ts.AddDate(0, 0, 3):  1.6,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		areEqual, actualValue, expectedValue := schedulesAreEqual(actual, expected)
		if !areEqual {
			t.Errorf("Received %f; want %f", actualValue, expectedValue)
		}
	})

	t.Run("Supplement nonbillables when scheduled nonbillables account for some nonbillable time", func(t *testing.T) {
		billableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 6,
			ts.AddDate(0, 0, -1): 6,
			ts.AddDate(0, 0, 1):  6,
			ts.AddDate(0, 0, 2):  6,
			ts.AddDate(0, 0, 3):  6,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		nonbillableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 1.5,
			ts.AddDate(0, 0, -1): 1.5,
			ts.AddDate(0, 0, 1):  1.5,
			ts.AddDate(0, 0, 2):  1.5,
			ts.AddDate(0, 0, 3):  1.5,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		timeOffSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 0,
			ts.AddDate(0, 0, -1): 0,
			ts.AddDate(0, 0, 1):  0,
			ts.AddDate(0, 0, 2):  0,
			ts.AddDate(0, 0, 3):  0,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		actual := getAdjustedNonbillablesSchedule(billableSchedule, nonbillableSchedule, timeOffSchedule)
		expected := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 2,
			ts.AddDate(0, 0, -1): 2,
			ts.AddDate(0, 0, 1):  2,
			ts.AddDate(0, 0, 2):  2,
			ts.AddDate(0, 0, 3):  2,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		areEqual, actualValue, expectedValue := schedulesAreEqual(actual, expected)
		if !areEqual {
			t.Errorf("Received %f; want %f", actualValue, expectedValue)
		}
	})

	t.Run("Does not supplement nonbillables when scheduled nonbillables account for all nonbillable time", func(t *testing.T) {
		billableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 3,
			ts.AddDate(0, 0, -1): 3,
			ts.AddDate(0, 0, 1):  3,
			ts.AddDate(0, 0, 2):  3,
			ts.AddDate(0, 0, 3):  3,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		nonbillableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 5,
			ts.AddDate(0, 0, -1): 5,
			ts.AddDate(0, 0, 1):  5,
			ts.AddDate(0, 0, 2):  5,
			ts.AddDate(0, 0, 3):  5,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		timeOffSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 0,
			ts.AddDate(0, 0, -1): 0,
			ts.AddDate(0, 0, 1):  0,
			ts.AddDate(0, 0, 2):  0,
			ts.AddDate(0, 0, 3):  0,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		actual := getAdjustedNonbillablesSchedule(billableSchedule, nonbillableSchedule, timeOffSchedule)
		expected := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 5,
			ts.AddDate(0, 0, -1): 5,
			ts.AddDate(0, 0, 1):  5,
			ts.AddDate(0, 0, 2):  5,
			ts.AddDate(0, 0, 3):  5,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		areEqual, actualValue, expectedValue := schedulesAreEqual(actual, expected)
		if !areEqual {
			t.Errorf("Received %f; want %f", actualValue, expectedValue)
		}
	})

	t.Run("Subtracts time off from billable and nonbillable expectations", func(t *testing.T) {
		billableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 3,
			ts.AddDate(0, 0, -1): 3,
			ts.AddDate(0, 0, 1):  3,
			ts.AddDate(0, 0, 2):  3,
			ts.AddDate(0, 0, 3):  3,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		nonbillableSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 5,
			ts.AddDate(0, 0, -1): 1,
			ts.AddDate(0, 0, 1):  1,
			ts.AddDate(0, 0, 2):  1,
			ts.AddDate(0, 0, 3):  1,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		timeOffSchedule := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 0,
			ts.AddDate(0, 0, -1): 2,
			ts.AddDate(0, 0, 1):  2,
			ts.AddDate(0, 0, 2):  2,
			ts.AddDate(0, 0, 3):  3,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		actual := getAdjustedNonbillablesSchedule(billableSchedule, nonbillableSchedule, timeOffSchedule)
		expected := fc.Schedule{
			// Monday
			ts.AddDate(0, 0, -2): 5,
			ts.AddDate(0, 0, -1): 3,
			ts.AddDate(0, 0, 1):  3,
			ts.AddDate(0, 0, 2):  3,
			ts.AddDate(0, 0, 3):  2,
			ts.AddDate(0, 0, 4):  0,
			// Sunday
			ts.AddDate(0, 0, 5): 0,
		}

		areEqual, actualValue, expectedValue := schedulesAreEqual(actual, expected)
		if !areEqual {
			t.Errorf("Received %f; want %f", actualValue, expectedValue)
		}
	})
}
