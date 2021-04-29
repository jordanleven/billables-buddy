package billablesbuddy

import (
	"time"
)

const (
	// The weekday we consider for weeks to begin on
	weekdayBeginWorkWeek = time.Monday
)

type StatisticDates struct {
	CurrentTimestamp time.Time
	WorkweekBegin    time.Time
	WorkweekEnd      time.Time
}

func getCurrentTimestamp() time.Time {
	return time.Now().Local()
}

func getWeekStartDateFromDate(date time.Time) time.Time {
	offset := int(weekdayBeginWorkWeek - date.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := date.AddDate(0, 0, offset)
	startWeekNoTime := time.Date(weekStartDate.Year(), weekStartDate.Month(), weekStartDate.Day(), 0, 0, 0, 0, weekStartDate.Location())
	return startWeekNoTime
}

func getStatisticDates() StatisticDates {
	ts := getCurrentTimestamp()
	startWeekDay := getWeekStartDateFromDate(ts)
	endDate := startWeekDay.AddDate(0, 0, 6)

	return StatisticDates{
		CurrentTimestamp: ts,
		WorkweekBegin:    startWeekDay,
		WorkweekEnd:      endDate,
	}
}
