package forecastclient

import (
	"time"

	"github.com/joefitzgerald/forecast"
)

type ForecastClient struct {
	Client *forecast.API
}

type ExpectedHours struct {
	HoursAll         TimeEntry
	HoursBillable    TimeEntry
	HoursNonbillable TimeEntry
}

func (c *ForecastClient) getHoursFromAssignments(startDate time.Time, a Assignments) ExpectedHours {
	projects := c.getProjets()
	evaluatorAll := func(a Assignment) bool { return true }
	evaluatorBillable := projects.isAssignmentBillable
	evaluatorNonbillable := func(a Assignment) bool { return !projects.isAssignmentBillable(a) }

	return ExpectedHours{
		HoursAll:         getEvaluatedHoursFromAssignments(startDate, a, evaluatorAll),
		HoursBillable:    getEvaluatedHoursFromAssignments(startDate, a, evaluatorBillable),
		HoursNonbillable: getEvaluatedHoursFromAssignments(startDate, a, evaluatorNonbillable),
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
