package harvestclient

import "log"

const (
	PathURLHarvestUserMe = "/users/me"
)

type HarvestUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"first_name"`
}

func (c HarvestClient) getCurrentUserID() int {
	var user HarvestUserResponse
	err := c.get(PathURLHarvestUserMe, DefaultArgs(), &user)

	if err != nil {
		log.Fatalln("Error retrieving Harvest user ID:", err)
	}

	return user.ID
}
