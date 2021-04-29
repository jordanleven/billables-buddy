package billablesbuddy

import (
	"testing"
	"time"

	fc "github.com/jordanleven/billables-buddy/internal/forecastclient"
)

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
