package billablesbuddy

import (
	"time"

	hc "github.com/jordanleven/billables-buddy/internal/harvestclient"
)

type HoursActual = hc.TrackedHours

func getCurrentWeeklyTrackedHours(a GetHoursStatisticsArgs, startDate time.Time, endDate time.Time) HoursActual {
	api := hc.GetHarvestAPI(a.HarvestAccountToken, a.HarvestAccountId)
	return api.GetTrackedHoursBetweenDates(startDate, endDate)
}
