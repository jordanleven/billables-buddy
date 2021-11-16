package harvestclient

import (
	"time"
)

type HarvestClient struct {
	Token     string
	AccountID string
}

type HoursConsolidated struct {
	HoursAll         TimeEntry
	HoursBillable    TimeEntry
	HoursNonbillable TimeEntry
}

type HarvestID = int

type ProjectHours struct {
	ProjectName string
	Hours       TimeEntry
}

type HoursByProject map[HarvestID]ProjectHours

type TrackedHours struct {
	HoursConsolidated HoursConsolidated
	HoursByProject    HoursByProject
	TodayStartTime    time.Time
}

func getHoursFromEntries(startDate time.Time, entries HarvestTimeEntriesResponse) HoursConsolidated {
	evaluatorTotal := func(e HarvestTimeEntryResponse) bool { return true }
	evaluatorBillable := isEntryBillable
	evaluatorNonbillable := func(e HarvestTimeEntryResponse) bool { return !isEntryBillable(e) }

	return HoursConsolidated{
		HoursAll:         getEvaluatedHoursFromEntries(startDate, entries, evaluatorTotal),
		HoursBillable:    getEvaluatedHoursFromEntries(startDate, entries, evaluatorBillable),
		HoursNonbillable: getEvaluatedHoursFromEntries(startDate, entries, evaluatorNonbillable),
	}
}

func getHoursByProjectFromEntries(startDate time.Time, entries HarvestTimeEntriesResponse) HoursByProject {
	projectsByID := make(HoursByProject)

	for _, entry := range entries.HarvestTimeEntries {
		projectID := entry.Project.ID
		projectName := entry.Project.Name
		if _, found := projectsByID[projectID]; found {
			continue
		}

		evaluatorIsCurrentProject := func(entryLocal HarvestTimeEntryResponse) bool { return entryLocal.Project.ID == projectID }

		projectsByID[projectID] = ProjectHours{
			ProjectName: projectName,
			Hours:       getEvaluatedHoursFromEntries(startDate, entries, evaluatorIsCurrentProject),
		}
	}

	return projectsByID
}

func (c HarvestClient) GetTrackedHoursBetweenDates(startDate time.Time, endDate time.Time) TrackedHours {
	ts := time.Now().Local()
	userID := c.getCurrentUserID()
	entries := c.getHarvestTimeEntries(userID, startDate, endDate)

	return TrackedHours{
		TodayStartTime:    getEarliestStartTimeFromEntries(ts, entries),
		HoursConsolidated: getHoursFromEntries(startDate, entries),
		HoursByProject:    getHoursByProjectFromEntries(startDate, entries),
	}
}

func GetHarvestAPI(token string, accountID string) *HarvestClient {
	return &HarvestClient{
		Token:     token,
		AccountID: accountID,
	}
}
