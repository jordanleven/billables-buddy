package billablesbuddy

import (
	"testing"
	"time"

	"bou.ke/monkey"
)

func TestGetWeekStartDateFromDate(t *testing.T) {
	location, _ := time.LoadLocation("EST")

	// Monday with no time
	expected := time.Date(1984, 01, 23, 0, 0, 0, 0, location)

	t.Run("Returns correct start week date on a Monday", func(t *testing.T) {
		currentTimeMonday := time.Date(1984, 01, 23, 0, 0, 0, 0, location)
		actual := getWeekStartDateFromDate(currentTimeMonday)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns correct start week date midday", func(t *testing.T) {
		currentTimeMidday := time.Date(1984, 01, 26, 8, 0, 0, 0, location)
		actual := getWeekStartDateFromDate(currentTimeMidday)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns correct start week date midweek", func(t *testing.T) {
		currentTimeMidweek := time.Date(1984, 01, 26, 0, 0, 0, 0, location)
		actual := getWeekStartDateFromDate(currentTimeMidweek)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns correct start week date on a Sunday", func(t *testing.T) {
		currentTimeSunday := time.Date(1984, 01, 29, 23, 59, 59, 0, location)
		actual := getWeekStartDateFromDate(currentTimeSunday)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}

func TestGetStatisticDates(t *testing.T) {
	location, _ := time.LoadLocation("EST")
	ts := time.Date(1984, 01, 24, 0, 0, 0, 0, location)

	// Stub our current timestamp function
	monkey.Patch(getCurrentTimestamp, func() time.Time {
		return ts
	})

	expectedDates := StatisticDates{
		CurrentTimestamp: ts,
		// Monday
		WorkweekBegin: time.Date(1984, 01, 23, 0, 0, 0, 0, location),
		// Friday
		WorkweekEnd: time.Date(1984, 01, 30, 0, 0, 0, 0, location),
	}

	actualDates := getStatisticDates()

	t.Run("Returns correct current timestamp", func(t *testing.T) {
		actual := actualDates.CurrentTimestamp
		expected := expectedDates.CurrentTimestamp
		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns correct workweek begin", func(t *testing.T) {
		actual := actualDates.WorkweekBegin
		expected := expectedDates.WorkweekBegin
		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns correct workweek end", func(t *testing.T) {
		actual := actualDates.WorkweekEnd
		expected := expectedDates.WorkweekEnd
		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}
