package billablesbuddy

import (
	"testing"
	"time"
)

func TestGetDateFromTime(t *testing.T) {
	t.Run("Returns the correctly formatted date", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 8, 33, 23, 0, time.UTC)

		actual := getDateFromTime(ts)
		expected := time.Date(1984, 1, 24, 0, 0, 0, 0, time.UTC)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}
