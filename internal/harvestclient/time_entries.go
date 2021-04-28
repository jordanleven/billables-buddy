package harvestclient

import (
	"fmt"
	"time"
)

const (
	PathURLHarvestTimeEntries = "/time_entries"
)

type HarvestTimeEntryResponse struct {
	Hours        float64   `json:"hours"`
	HoursRounded float64   `json:"rounded_hours"`
	Billable     bool      `json:"billable"`
	Date         string    `json:"spent_date"`
	Notes        string    `json:"notes"`
	TimeStart    time.Time `json:"created_at"`
}

type HarvestTimeEntriesResponse struct {
	HarvestTimeEntries []HarvestTimeEntryResponse `json:"time_entries"`
}

type ScheduleStatistics interface {
	getTotal() float64
}

type Schedule map[time.Weekday]float64

type TimeEntry struct {
	Total    float64
	Schedule Schedule
}

type hoursEvaluator func(HarvestTimeEntryResponse) bool

func getFormattedHarvestAPIDate(ts time.Time) string {
	tsUTC := getUTCTimeFromLocalTime(ts)
	tsUTCFormatted := getFormattedDate(tsUTC)
	return tsUTCFormatted
}

func isEntryBillable(t HarvestTimeEntryResponse) bool {
	return t.Billable
}

func getTotalHoursFromEvaluator(t HarvestTimeEntriesResponse, evaluator hoursEvaluator) float64 {
	var hours float64 = 0.0
	for _, entry := range t.HarvestTimeEntries {
		if evaluator(entry) {
			hours = hours + entry.Hours
		}
	}

	return hours
}

func getScheduledHoursFromEvaluator(t HarvestTimeEntriesResponse, evaluator hoursEvaluator) Schedule {
	schedule := Schedule{
		time.Monday:    0.0,
		time.Tuesday:   0.0,
		time.Wednesday: 0.0,
		time.Thursday:  0.0,
		time.Friday:    0.0,
		time.Saturday:  0.0,
		time.Sunday:    0.0,
	}

	for _, entry := range t.HarvestTimeEntries {
		if evaluator(entry) {
			dateP, _ := time.Parse("2006-01-02", entry.Date)
			weekDay := dateP.Weekday()
			schedule[weekDay] += entry.Hours
		}
	}

	return schedule
}

func getEvaluatedHoursFromEntries(t HarvestTimeEntriesResponse, evaluator hoursEvaluator) TimeEntry {
	return TimeEntry{
		Total:    getTotalHoursFromEvaluator(t, evaluator),
		Schedule: getScheduledHoursFromEvaluator(t, evaluator),
	}
}

func (c HarvestClient) getHarvestTimeEntries(userID int, startDate time.Time, endDate time.Time) HarvestTimeEntriesResponse {
	var entries HarvestTimeEntriesResponse
	args := Arguments{
		"user_id": fmt.Sprint(userID),
		"from":    getFormattedHarvestAPIDate(startDate),
		"to":      getFormattedHarvestAPIDate(endDate),
	}
	c.get(PathURLHarvestTimeEntries, args, &entries)

	return entries
}
