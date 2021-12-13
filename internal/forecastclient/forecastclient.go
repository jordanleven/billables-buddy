package forecastclient

import (
	"time"

	"github.com/joefitzgerald/forecast"
)

type ForecastClient struct {
	Client *forecast.API
}

type ExpectedHoursConsolidated struct {
	HoursAll         TimeEntry
	HoursBillable    TimeEntry
	HoursNonbillable TimeEntry
	HoursTimeOff     TimeEntry
}

type ExpectedHoursByProject struct {
	ProjectName string
	HarvestID   int
	Hours       TimeEntry
}

type ExpectedHours struct {
	HoursConsolidated ExpectedHoursConsolidated
	HoursByProject    map[int]ExpectedHoursByProject
}

func getConsolidatedHoursFromAssignments(startDate time.Time, a Assignments, projects Projects) ExpectedHoursConsolidated {
	evaluatorAll := func(a Assignment) bool { return !isAssignmentTimeOff(a) }
	evaluatorBillable := projects.isAssignmentBillable
	evaluatorNonbillable := func(a Assignment) bool { return !projects.isAssignmentBillable(a) && !isAssignmentTimeOff(a) }

	return ExpectedHoursConsolidated{
		HoursAll:         getEvaluatedHoursFromAssignments(startDate, a, evaluatorAll),
		HoursBillable:    getEvaluatedHoursFromAssignments(startDate, a, evaluatorBillable),
		HoursNonbillable: getEvaluatedHoursFromAssignments(startDate, a, evaluatorNonbillable),
		HoursTimeOff:     getEvaluatedHoursFromAssignments(startDate, a, isAssignmentTimeOff),
	}
}

func getCombinedProjectTimeEntry(a TimeEntry, b TimeEntry) TimeEntry {
	updatedTimeEntry := a
	updatedTimeEntry.Total = 0

	for aDate, aValue := range updatedTimeEntry.Schedule {
		bValue := b.Schedule[aDate]
		updatedTimeEntry.Total += aValue + bValue
		updatedTimeEntry.Schedule[aDate] += bValue
	}

	return updatedTimeEntry
}

func getHoursByProjectFromAssignments(startDate time.Time, a Assignments, p Projects) map[int]ExpectedHoursByProject {
	projectsByID := make(map[int]ExpectedHoursByProject)

	for _, assignment := range a {
		assignmentID := assignment.ID
		projectID := assignment.ProjectID
		evaluatorIsCurrentAssignment := func(aLocal Assignment) bool { return aLocal.ID == assignmentID }

		projectName := p[projectID].Name
		harvestID := p[projectID].HarvestID
		hours := getEvaluatedHoursFromAssignments(startDate, a, evaluatorIsCurrentAssignment)

		if harvestID == 0 {
			continue
		}

		// If the project assignment already exists, then we have a split assignment
		if _, projectAssignmentExists := projectsByID[projectID]; projectAssignmentExists {
			hours = getCombinedProjectTimeEntry(projectsByID[projectID].Hours, hours)
		}

		projectsByID[projectID] = ExpectedHoursByProject{
			ProjectName: projectName,
			HarvestID:   harvestID,
			Hours:       hours,
		}
	}

	return projectsByID
}

func (c *ForecastClient) getHoursFromAssignments(startDate time.Time, a Assignments) ExpectedHours {
	projects := c.getProjets()
	hoursConsolidated := getConsolidatedHoursFromAssignments(startDate, a, projects)
	hoursByProject := getHoursByProjectFromAssignments(startDate, a, projects)

	return ExpectedHours{
		HoursConsolidated: hoursConsolidated,
		HoursByProject:    hoursByProject,
	}
}

func (c *ForecastClient) GetExpectedHoursBetweenDates(startDate time.Time, endDate time.Time) ExpectedHours {
	assignments := c.getUserAssignments(startDate, endDate)
	return c.getHoursFromAssignments(startDate, assignments)
}

func GetForecastAPI(token string, accountID string) *ForecastClient {
	c := forecast.New("https://api.forecastapp.com", accountID, token)
	return &ForecastClient{
		Client: c,
	}
}
