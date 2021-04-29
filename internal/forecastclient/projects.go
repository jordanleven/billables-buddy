package forecastclient

import (
	"github.com/joefitzgerald/forecast"
)

type Project struct {
	IsBillable bool
}

type Projects map[int]Project

type forecastProject = forecast.RemainingBudgetedHoursItem
type forecastProjects = forecast.RemainingBudgetedHours

func (p Projects) isAssignmentBillable(a Assignment) bool {
	return p[a.ProjectID].IsBillable
}

func isProjectBillable(p forecastProject) bool {
	return p.BudgetBy == "project"
}

func getFormattedProjects(allProjects forecastProjects) Projects {
	projects := make(Projects)

	for _, project := range allProjects {
		projectF := Project{
			IsBillable: isProjectBillable(project),
		}
		projects[project.ProjectID] = projectF
	}

	return projects
}

func (c *ForecastClient) getProjectsWithRemainingBudged() forecastProjects {
	allProjects, _ := c.Client.RemainingBudgetedHours()
	return allProjects
}

func (c *ForecastClient) getProjets() Projects {
	allProjects := c.getProjectsWithRemainingBudged()
	return getFormattedProjects(allProjects)
}
