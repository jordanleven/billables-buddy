package harvestclient

import (
	"fmt"
	"log"
	"time"
)

const (
	PathURLHarvestTimeEntries = "/time_entries"
)

type HarvestProject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HarvestTimeEntryResponse struct {
	Hours        float64 `json:"hours"`
	HoursRounded float64 `json:"rounded_hours"`
	Project      HarvestProject
	Billable     bool      `json:"billable"`
	Date         string    `json:"spent_date"`
	Notes        string    `json:"notes"`
	TimeStart    time.Time `json:"created_at"`
}

type HarvestTimeEntriesResponse struct {
	HarvestTimeEntries []HarvestTimeEntryResponse `json:"time_entries"`
}

type Schedule map[time.Time]float64

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
			hours += entry.Hours
		}
	}

	return hours
}

func getScheduledHoursFromEvaluator(startDate time.Time, t HarvestTimeEntriesResponse, evaluator hoursEvaluator) Schedule {
	loc := startDate.Location()
	schedule := Schedule{
		startDate.AddDate(0, 0, 0): 0.0,
		startDate.AddDate(0, 0, 1): 0.0,
		startDate.AddDate(0, 0, 2): 0.0,
		startDate.AddDate(0, 0, 3): 0.0,
		startDate.AddDate(0, 0, 4): 0.0,
		startDate.AddDate(0, 0, 5): 0.0,
		startDate.AddDate(0, 0, 6): 0.0,
	}

	for _, entry := range t.HarvestTimeEntries {
		if evaluator(entry) {
			dateL := entry.TimeStart.In(loc)
			dateLF := time.Date(dateL.Year(), dateL.Month(), dateL.Day(), 0, 0, 0, 0, loc)
			schedule[dateLF] += entry.Hours
		}
	}

	return schedule
}

func getEarliestStartTimeFromEntries(ts time.Time, t HarvestTimeEntriesResponse) time.Time {
	const timeFormat = "2006-01-02"
	loc := ts.Location()
	tsF := ts.Format(timeFormat)

	var todayStartTime time.Time
	for _, entry := range t.HarvestTimeEntries {
		startTime := entry.TimeStart
		startTimeL := startTime.In(loc)
		startTimeLF := startTimeL.Format(timeFormat)

		if startTimeLF != tsF {
			continue
		}

		if todayStartTime.IsZero() || todayStartTime.After(startTime) {
			todayStartTime = startTimeL
		}
	}

	return todayStartTime
}

func getEvaluatedHoursFromEntries(startDate time.Time, t HarvestTimeEntriesResponse, evaluator hoursEvaluator) TimeEntry {
	return TimeEntry{
		Total:    getTotalHoursFromEvaluator(t, evaluator),
		Schedule: getScheduledHoursFromEvaluator(startDate, t, evaluator),
	}
}

func (c HarvestClient) getHarvestTimeEntries(userID int, startDate time.Time, endDate time.Time) HarvestTimeEntriesResponse {
	var entries HarvestTimeEntriesResponse
	args := Arguments{
		"user_id": fmt.Sprint(userID),
		"from":    getFormattedHarvestAPIDate(startDate),
		"to":      getFormattedHarvestAPIDate(endDate),
	}
	err := c.get(PathURLHarvestTimeEntries, args, &entries)

	if err != nil {
		log.Fatalln("Error retrieving Harvest time entries:", err)
	}

	return entries
}
