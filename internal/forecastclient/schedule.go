package forecastclient

import (
	t "time"
)

type Schedule map[t.Time]float64

func GetWeeklyScheduleFromStartDate(startDate t.Time) Schedule {
	return Schedule{
		startDate.AddDate(0, 0, 0): 0.0,
		startDate.AddDate(0, 0, 1): 0.0,
		startDate.AddDate(0, 0, 2): 0.0,
		startDate.AddDate(0, 0, 3): 0.0,
		startDate.AddDate(0, 0, 4): 0.0,
		startDate.AddDate(0, 0, 5): 0.0,
		startDate.AddDate(0, 0, 6): 0.0,
	}
}
