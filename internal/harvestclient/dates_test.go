package harvestclient

import (
	"testing"
	"time"
)

func TestGetUTCTimeFromLocalTime(t *testing.T) {
	t.Run("Returns the expected time", func(t *testing.T) {
		loc, _ := time.LoadLocation("EST")
		ts := time.Date(1984, 01, 24, 8, 0, 0, 0, loc)
		expected := time.Date(1984, 01, 24, 13, 0, 0, 0, time.UTC)

		actual := getUTCTimeFromLocalTime(ts)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}

func TestGetFormattedDate(t *testing.T) {
	t.Run("Returns the expected time", func(t *testing.T) {
		ts := time.Date(1984, 01, 24, 8, 6, 29, 0, time.UTC)
		expected := "1984-01-24T08:06:29Z"

		actual := getFormattedDate(ts)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}
