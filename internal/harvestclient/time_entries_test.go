package harvestclient

import (
	"testing"
	"time"
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
	startTime := time.Date(1984, 01, 23, 0, 0, 0, 0, time.UTC)
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

	actualSchedule := getScheduledHoursFromEvaluator(startTime, args, evaluator)
	expectedSchedule := Schedule{
		startTime.AddDate(0, 0, 0): 0.0,
		startTime.AddDate(0, 0, 1): 6,
		startTime.AddDate(0, 0, 2): 1.2,
		startTime.AddDate(0, 0, 3): 0.0,
		startTime.AddDate(0, 0, 4): 0.0,
		startTime.AddDate(0, 0, 5): 0.0,
		startTime.AddDate(0, 0, 6): 0.0,
	}

	t.Run("Returns the expected schedule when no hours are on a day", func(t *testing.T) {
		monday := startTime.AddDate(0, 0, 0)
		actual := actualSchedule[monday]
		expected := expectedSchedule[monday]

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the expected schedule when single entries are on a day", func(t *testing.T) {
		wednesday := startTime.AddDate(0, 0, 2)

		actual := actualSchedule[wednesday]
		expected := expectedSchedule[wednesday]

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})

	t.Run("Returns the expected schedule when multiple entries are on a day", func(t *testing.T) {
		wednesday := startTime.AddDate(0, 0, 3)

		actual := actualSchedule[wednesday]
		expected := expectedSchedule[wednesday]

		if actual != expected {
			t.Errorf("Received %f; want %f", actual, expected)
		}
	})
}

func TestGetEarliestStartTimeFromEntries(t *testing.T) {
	loc, _ := time.LoadLocation("EST")
	entries := HarvestTimeEntriesResponse{
		HarvestTimeEntries: []HarvestTimeEntryResponse{
			{
				TimeStart: time.Date(1984, 1, 24, 6, 35, 0, 0, time.UTC),
			},
			{
				TimeStart: time.Date(1984, 1, 24, 6, 29, 0, 0, time.UTC),
			},
			{
				TimeStart: time.Date(1984, 1, 24, 6, 40, 0, 0, time.UTC),
			},
			{
				TimeStart: time.Date(1984, 1, 22, 6, 0, 0, 0, time.UTC),
			},
		},
	}

	t.Run("Returns the correct start time when providing local timezone", func(t *testing.T) {
		ts := time.Date(1984, 1, 24, 11, 0, 0, 0, loc)

		actual := getEarliestStartTimeFromEntries(ts, entries)
		expected := time.Date(1984, 1, 24, 1, 29, 0, 0, loc)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correct start time when providing a single entry that matches", func(t *testing.T) {
		ts := time.Date(1984, 1, 22, 6, 35, 0, 0, loc)

		actual := getEarliestStartTimeFromEntries(ts, entries)
		expected := time.Date(1984, 1, 22, 1, 0, 0, 0, loc)

		if actual != expected {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns the correct zero start time when providing no entries that match", func(t *testing.T) {
		ts := time.Date(1984, 1, 23, 6, 1, 0, 0, loc)

		actual := getEarliestStartTimeFromEntries(ts, entries)

		if !actual.IsZero() {
			t.Errorf("Received %s; want zero", actual)
		}
	})
}
