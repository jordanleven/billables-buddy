package forecastclient

import (
	"testing"
)

func TestIsBillableProject(t *testing.T) {
	t.Run("Returns false when the project is not billable", func(t *testing.T) {
		actual := isProjectBillable(forecastProjectRemainingBudget{BudgetBy: "none"})
		expected := false

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})

	t.Run("Returns true when the project is billable", func(t *testing.T) {
		actual := isProjectBillable(forecastProjectRemainingBudget{BudgetBy: "project"})
		expected := true

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})
}

func TestGetFormattedProjects(t *testing.T) {
	projectsWithRemainingBudget := forecastProjectsRemainingBudget{
		forecastProjectRemainingBudget{
			ProjectID: 1,
			BudgetBy:  "project",
		},
		forecastProjectRemainingBudget{
			ProjectID: 2,
			BudgetBy:  "project",
		},
		forecastProjectRemainingBudget{
			ProjectID: 3,
			BudgetBy:  "none",
		},
	}

	projectNames := map[int]ForecastProject{
		1: {
			Name:      "Project A",
			HarvestID: 100,
		},
		2: {
			Name:      "Project B",
			HarvestID: 101,
		},
		3: {
			Name:      "Project C",
			HarvestID: 102,
		},
	}

	actualProjects := getFormattedProjects(projectsWithRemainingBudget, projectNames)

	t.Run("For billable projects", func(t *testing.T) {
		actual := actualProjects[1]
		expected := Project{
			IsBillable: true,
			Name:       "Project A",
			HarvestID:  100,
		}

		t.Run("Returns the project as billable", func(t *testing.T) {
			if actual.IsBillable != expected.IsBillable {
				t.Errorf("Received %t; want %t", actual.IsBillable, expected.IsBillable)
			}
		})

		t.Run("Returns the project with the correct name", func(t *testing.T) {
			if actual.Name != expected.Name {
				t.Errorf("Received %s; want %s", actual.Name, expected.Name)
			}
		})

		t.Run("Returns the project with the correct Harvest ID", func(t *testing.T) {
			if actual.HarvestID != expected.HarvestID {
				t.Errorf("Received %d; want %d", actual.HarvestID, expected.HarvestID)
			}
		})
	})

	t.Run("For non-billable projects", func(t *testing.T) {
		actual := actualProjects[3]
		expected := Project{
			IsBillable: false,
			Name:       "Project C",
			HarvestID:  102,
		}

		t.Run("Returns the project as nonbillable", func(t *testing.T) {
			if actual.IsBillable != expected.IsBillable {
				t.Errorf("Received %t; want %t", actual.IsBillable, expected.IsBillable)
			}
		})

		t.Run("Returns the project with the correct name", func(t *testing.T) {
			if actual.Name != expected.Name {
				t.Errorf("Received %s; want %s", actual.Name, expected.Name)
			}
		})

		t.Run("Returns the project with the correct Harvest ID", func(t *testing.T) {
			if actual.HarvestID != expected.HarvestID {
				t.Errorf("Received %d; want %d", actual.HarvestID, expected.HarvestID)
			}
		})
	})
}
