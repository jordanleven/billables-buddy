package forecastclient

import (
	"github.com/joefitzgerald/forecast"
)

type Projects map[int]Project

type forecastProjectRemainingBudget = forecast.RemainingBudgetedHoursItem
type forecastProjectsRemainingBudget = forecast.RemainingBudgetedHours

type Project struct {
	IsBillable bool
	Name       string
	HarvestID  int
}

type ForecastProject struct {
	Name      string
	HarvestID int
}

const (
	TimeOffProjectID = 19829
)

func (p Projects) isAssignmentBillable(a Assignment) bool {
	return p[a.ProjectID].IsBillable
}

func isAssignmentTimeOff(a Assignment) bool {
	return a.ProjectID == TimeOffProjectID
}

func isProjectBillable(p forecastProjectRemainingBudget) bool {
	return p.BudgetBy != "none"
}

func getFormattedProjects(allProjects forecastProjectsRemainingBudget, projectNames map[int]ForecastProject) Projects {
	projects := make(Projects)

	for _, project := range allProjects {
		projectID := project.ProjectID

		projectF := Project{
			IsBillable: isProjectBillable(project),
			Name:       projectNames[projectID].Name,
			HarvestID:  projectNames[projectID].HarvestID,
		}
		projects[project.ProjectID] = projectF
	}

	return projects
}

func (c *ForecastClient) getProjectsWithRemainingBudged() forecastProjectsRemainingBudget {
	allProjectsWithRemainingBudget, _ := c.Client.RemainingBudgetedHours()
	return allProjectsWithRemainingBudget
}

func (c *ForecastClient) getProjectNames() map[int]ForecastProject {
	allProjects, _ := c.Client.Projects()
	projectsByName := make(map[int]ForecastProject)

	for _, project := range allProjects {
		projectID := project.ID
		projectName := project.Name
		projectHarvestID := project.HarvestID
		isArchived := project.Archived

		if !isArchived {
			projectsByName[projectID] = ForecastProject{
				Name:      projectName,
				HarvestID: projectHarvestID,
			}
		}
	}

	return projectsByName
}

func (c *ForecastClient) getProjets() Projects {
	allProjectsWithRemainingBudget := c.getProjectsWithRemainingBudged()
	projectNames := c.getProjectNames()
	return getFormattedProjects(allProjectsWithRemainingBudget, projectNames)
}
