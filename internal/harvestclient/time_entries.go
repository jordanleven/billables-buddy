package harvestclient

import (
	fc "billables-buddy/internal/forecastclient"
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

type TimeEntry struct {
	Total    float64
	Schedule fc.Schedule
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

func getScheduledHoursFromEvaluator(startDate time.Time, t HarvestTimeEntriesResponse, evaluator hoursEvaluator) fc.Schedule {
	loc := startDate.Location()
	schedule := fc.GetWeeklyScheduleFromStartDate(startDate)

	for _, entry := range t.HarvestTimeEntries {
		if evaluator(entry) {
			date, _ := time.Parse("2006-01-02", entry.Date)
			// The dates that come back from Harvest are actually localized to the timezone they
			// were entered in.
			dateL := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
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
