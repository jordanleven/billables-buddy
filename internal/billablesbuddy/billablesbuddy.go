package billablesbuddy

import (
	"sort"

	fc "billables-buddy/internal/forecastclient"
)

type GetHoursStatisticsArgs struct {
	ForecastAccountID   string
	HarvestAccountId    string
	HarvestAccountToken string
}

type HoursStatistic struct {
	HoursExpected    float64
	HoursActual      float64
	HoursActualToday float64
}

type HoursProject struct {
	ProjectName  string
	ProjectHours HoursStatistic
}

type HoursStatistics struct {
	Status
	Person           fc.Person
	HoursToday       float64
	HoursAll         HoursStatistic
	HoursBillable    HoursStatistic
	HoursNonbillable HoursStatistic
	HoursRemaining
	HoursByProject []HoursProject
}

func getTrackedHoursByProject(hoursByProject []ProjectHours) []HoursProject {
	var trackedHoursByProject []HoursProject

	for _, project := range hoursByProject {
		hoursProject := HoursProject{
			ProjectName: project.Name,
			ProjectHours: HoursStatistic{
				HoursActual:   project.Hours.ActualCurrent,
				HoursExpected: project.Hours.ExpectedCurrent,
			},
		}
		trackedHoursByProject = append(trackedHoursByProject, hoursProject)
	}

	sort.Slice(trackedHoursByProject, func(i int, j int) bool {
		return trackedHoursByProject[i].ProjectName < trackedHoursByProject[j].ProjectName
	})

	return trackedHoursByProject
}

func GetTrackedHoursStatistics(args GetHoursStatisticsArgs) HoursStatistics {
	statDates := getStatisticDates()
	h := getActualAndExpectedHours(args, statDates)
	hConsolidated := h.HoursConsolidated
	hConsolidatedAll := hConsolidated.HoursAll
	hConsolidatedBillable := hConsolidated.HoursBillable
	hConsolidatedNonBillable := hConsolidated.HoursNonbillable
	hByProject := getTrackedHoursByProject(h.HoursByProject)

	s := getCurrentStatus(hConsolidatedBillable.ActualCurrent, hConsolidatedBillable.ExpectedCurrent, hConsolidatedBillable.ExpectedEndOfWeek)
	hr := getHoursRemaining(statDates.CurrentTimestamp, h.TodayStartTime, hConsolidatedBillable, hConsolidatedNonBillable)

	return HoursStatistics{
		Status:         s,
		Person:         h.Person,
		HoursRemaining: hr,
		HoursAll: HoursStatistic{
			HoursActual:      hConsolidatedAll.ActualCurrent,
			HoursExpected:    hConsolidatedAll.ExpectedCurrent,
			HoursActualToday: hConsolidatedAll.ActualToday,
		},
		HoursBillable: HoursStatistic{
			HoursActual:      hConsolidatedBillable.ActualCurrent,
			HoursExpected:    hConsolidatedBillable.ExpectedCurrent,
			HoursActualToday: hConsolidatedBillable.ActualToday,
		},
		HoursNonbillable: HoursStatistic{
			HoursActual:      hConsolidatedNonBillable.ActualCurrent,
			HoursExpected:    hConsolidatedNonBillable.ExpectedCurrent,
			HoursActualToday: hConsolidatedNonBillable.ActualToday,
		},
		HoursByProject: hByProject,
	}
}
