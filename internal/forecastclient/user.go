package forecastclient

type UserID = int

type Person struct {
	NameFirst string
	NameLast  string
}

func (c *ForecastClient) getCurrentUserID() UserID {
	currentUser, _ := c.Client.WhoAmI()
	return currentUser.ID
}

func (c *ForecastClient) getCurrentPerson() Person {
	currentUserID := c.getCurrentUserID()
	person, _ := c.Client.Person(currentUserID)

	return Person{
		NameFirst: person.FirstName,
		NameLast:  person.LastName,
	}
}
