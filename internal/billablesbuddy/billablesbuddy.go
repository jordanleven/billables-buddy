package billablesbuddy

type GetHoursStatisticsArgs struct {
	HarvestAccountId    string
	HarvestAccountToken string
}

type Hours struct {
	HoursActual float64
}

type HoursStatistics struct {
	HoursTotal       Hours
	HoursBillable    Hours
	HoursNonbillable Hours
}

func GetTrackedHoursStatistics(args GetHoursStatisticsArgs) HoursStatistics {
	statDates := getStatisticDates()
	workweekBegin := statDates.WorkweekBegin
	workweekEnd := statDates.WorkweekEnd
	h := getCurrentWeeklyTrackedHours(args, workweekBegin, workweekEnd)

	return HoursStatistics{
		HoursTotal: Hours{
			HoursActual: h.Hours.HoursAll.Total,
		},
		HoursBillable: Hours{
			HoursActual: h.Hours.HoursBillable.Total,
		},
		HoursNonbillable: Hours{
			HoursActual: h.Hours.HoursNonbillable.Total,
		},
	}
}
