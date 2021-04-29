package forecastclient

import (
	"testing"
)

func TestIsBillableProject(t *testing.T) {
	t.Run("Returns false when the project is not billable", func(t *testing.T) {
		actual := isProjectBillable(forecastProject{BudgetBy: "in-house"})
		expected := false

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})

	t.Run("Returns true when the project is billable", func(t *testing.T) {
		actual := isProjectBillable(forecastProject{BudgetBy: "project"})
		expected := true

		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})
}

func TestGetFormattedProjects(t *testing.T) {
	projects := forecastProjects{
		forecastProject{
			ProjectID: 1,
			BudgetBy:  "project",
		},
		forecastProject{
			ProjectID: 2,
			BudgetBy:  "project",
		},
		forecastProject{
			ProjectID: 2,
			BudgetBy:  "in-house",
		},
	}
	actualProjects := getFormattedProjects(projects)

	t.Run("Returns the expected formatted projects when projects are billable", func(t *testing.T) {
		actual := actualProjects[1]
		expected := Project{IsBillable: true}
		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})

	t.Run("Returns the expected formatted projects when projects are not billable", func(t *testing.T) {
		actual := actualProjects[3]
		expected := Project{IsBillable: false}
		if actual != expected {
			t.Errorf("Received %t; want %t", actual, expected)
		}
	})
}
