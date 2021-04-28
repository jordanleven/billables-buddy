package harvestclient

const (
	PathURLHarvestUserMe = "/users/me"
)

type HarvestUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"first_name"`
}

func (c HarvestClient) getCurrentUserID() int {
	var user HarvestUserResponse
	c.get(PathURLHarvestUserMe, DefaultArgs(), &user)

	return user.ID
}
