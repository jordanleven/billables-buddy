package forecastclient

import (
	"github.com/joefitzgerald/forecast"
)

type UserID = *forecast.CurrentUser

func (c *ForecastClient) getCurrentUserID() UserID {
	uid, _ := c.Client.WhoAmI()
	return uid
}
