package event

func CreateNewEvent(name string, seats int) error {
	return CreateEvent(Event{
		Name:       name,
		TotalSeats: seats,
	})
}

func Register(eventID int) error {
	return RegisterForEvent(eventID)
}