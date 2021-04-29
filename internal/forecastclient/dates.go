package forecastclient

import "time"

func getUTCTimeFromLocalTime(ts time.Time) time.Time {
	return ts.In(time.UTC)
}

func getFormattedDate(ts time.Time) string {
	return ts.Format(time.RFC3339)
}
