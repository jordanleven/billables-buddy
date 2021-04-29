package harvestclient

import (
	"time"
)

type HarvestClient struct {
	Token     string
	AccountID string
}

type Hours struct {
	HoursAll         TimeEntry
	HoursBillable    TimeEntry
	HoursNonbillable TimeEntry
}

type TrackedHours struct {
	Hours          Hours
	TodayStartTime time.Time
}

func getHoursFromEntries(startDate time.Time, entries HarvestTimeEntriesResponse) Hours {
	evaluatorTotal := func(e HarvestTimeEntryResponse) bool { return true }
	evaluatorBillable := isEntryBillable
	evaluatorNonbillable := func(e HarvestTimeEntryResponse) bool { return !isEntryBillable(e) }

	return Hours{
		HoursAll:         getEvaluatedHoursFromEntries(startDate, entries, evaluatorTotal),
		HoursBillable:    getEvaluatedHoursFromEntries(startDate, entries, evaluatorBillable),
		HoursNonbillable: getEvaluatedHoursFromEntries(startDate, entries, evaluatorNonbillable),
	}
}

func (c HarvestClient) GetTrackedHoursBetweenDates(startDate time.Time, endDate time.Time) TrackedHours {
	ts := time.Now().Local()
	userID := c.getCurrentUserID()
	entries := c.getHarvestTimeEntries(userID, startDate, endDate)

	return TrackedHours{
		Hours:          getHoursFromEntries(startDate, entries),
		TodayStartTime: getEarliestStartTimeFromEntries(ts, entries),
	}
}

func GetHarvestAPI(token string, accountID string) *HarvestClient {
	return &HarvestClient{
		Token:     token,
		AccountID: accountID,
	}
}
