package billablesbuddy

type GetHoursStatisticsArgs struct {
	ForecastAccountID   string
	HarvestAccountId    string
	HarvestAccountToken string
}

type HoursStatistic struct {
	HoursExpected float64
	HoursActual   float64
}

type HoursStatistics struct {
	Status
	HoursAll         HoursStatistic
	HoursBillable    HoursStatistic
	HoursNonbillable HoursStatistic
	HoursRemaining
}

func GetTrackedHoursStatistics(args GetHoursStatisticsArgs) HoursStatistics {
	statDates := getStatisticDates()
	h := getActualAndExpectedHours(args, statDates)
	s := getCurrentStatus(h.HoursBillable.Actual, h.HoursBillable.Expected, h.HoursBillable.ExpectedTotal)
	hr := getHoursRemaining(statDates.CurrentTimestamp, h.TodayStartTime, h.HoursBillable, h.HoursNonbillable)

	return HoursStatistics{
		Status:         s,
		HoursRemaining: hr,
		HoursAll: HoursStatistic{
			HoursActual:   h.HoursAll.Actual,
			HoursExpected: h.HoursAll.Expected,
		},
		HoursBillable: HoursStatistic{
			HoursActual:   h.HoursBillable.Actual,
			HoursExpected: h.HoursBillable.Expected,
		},
		HoursNonbillable: HoursStatistic{
			HoursActual:   h.HoursNonbillable.Actual,
			HoursExpected: h.HoursNonbillable.Expected,
		},
	}
}
