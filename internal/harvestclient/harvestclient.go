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
	Hours Hours
}

func getHoursFromEntries(entries HarvestTimeEntriesResponse) Hours {
	evaluatorTotal := func(e HarvestTimeEntryResponse) bool { return true }
	evaluatorBillable := isEntryBillable
	evaluatorNonbillable := func(e HarvestTimeEntryResponse) bool { return !isEntryBillable(e) }

	return Hours{
		HoursAll:         getEvaluatedHoursFromEntries(entries, evaluatorTotal),
		HoursBillable:    getEvaluatedHoursFromEntries(entries, evaluatorBillable),
		HoursNonbillable: getEvaluatedHoursFromEntries(entries, evaluatorNonbillable),
	}
}

func (c HarvestClient) GetTrackedHoursBetweenDates(startDate time.Time, endDate time.Time) TrackedHours {
	userID := c.getCurrentUserID()
	entries := c.getHarvestTimeEntries(userID, startDate, endDate)

	return TrackedHours{
		Hours: getHoursFromEntries(entries),
	}
}

func GetHarvestAPI(token string, accountID string) *HarvestClient {
	return &HarvestClient{
		Token:     token,
		AccountID: accountID,
	}
}
