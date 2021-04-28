package harvestclient

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestGetFormattedHarvestAPIDate(t *testing.T) {
	t.Run("Returns the correctly formatted date", func(t *testing.T) {
		loc, _ := time.LoadLocation("EST")
		ts := time.Date(1984, 1, 24, 8, 0, 0, 0, loc)

		actual := getFormattedHarvestAPIDate(ts)
		expected := "1984-01-24T13:00:00Z"

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}

func TestIsEntryBillable(t *testing.T) {
	t.Run("Returns false when an entry is not considered to be billable", func(t *testing.T) {
		entry := HarvestTimeEntryResponse{
			Billable: false,
		}

		actual := isEntryBillable(entry)
		expected := false

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})

	t.Run("Returns true when an entry is considered to be billable", func(t *testing.T) {
		entry := HarvestTimeEntryResponse{
			Billable: true,
		}

		actual := isEntryBillable(entry)
		expected := true

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})
}

func TestGetTotalHoursFromEvaluator(t *testing.T) {
	t.Run("Returns the expected number of total hours based on the evaluator", func(t *testing.T) {
		evaluator := func(t HarvestTimeEntryResponse) bool {
			return t.Billable
		}
		args := HarvestTimeEntriesResponse{
			HarvestTimeEntries: []HarvestTimeEntryResponse{
				{
					Billable: false,
					Hours:    1.0,
				},
				{
					Billable: true,
					Hours:    2.5,
				},
				{
					Billable: true,
					Hours:    1.2,
				},
			},
		}

		actual := getTotalHoursFromEvaluator(args, evaluator)
		expected := 3.7

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetScheduledHoursFromEvaluator(t *testing.T) {
	evaluator := func(t HarvestTimeEntryResponse) bool {
		return t.Billable
	}
	args := HarvestTimeEntriesResponse{
		HarvestTimeEntries: []HarvestTimeEntryResponse{
			{
				Billable: false,
				Hours:    1.0,
				Date:     "1984-01-24",
			},
			{
				Billable: true,
				Hours:    2.5,
				Date:     "1984-01-24",
			},
			{
				Billable: true,
				Hours:    3.5,
				Date:     "1984-01-24",
			},
			{
				Billable: true,
				Hours:    1.2,
				Date:     "1984-01-25",
			},
		},
	}

	actualSchedule := getScheduledHoursFromEvaluator(args, evaluator)
	expectedSchedule := Schedule{
		time.Monday:    0.0,
		time.Tuesday:   6,
		time.Wednesday: 1.2,
		time.Thursday:  0.0,
		time.Friday:    0.0,
		time.Saturday:  0.0,
		time.Sunday:    0.0,
	}

	t.Run("Returns the expected schedule when no hours are on a day", func(t *testing.T) {
		actual := actualSchedule[time.Monday]
		expected := expectedSchedule[time.Monday]

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the expected schedule when single entries are on a day", func(t *testing.T) {
		actual := actualSchedule[time.Wednesday]
		expected := expectedSchedule[time.Wednesday]

		spew.Dump(actualSchedule)
		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the expected schedule when multiple entries are on a day", func(t *testing.T) {
		actual := actualSchedule[time.Tuesday]
		expected := expectedSchedule[time.Tuesday]

		spew.Dump(actualSchedule)
		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}
