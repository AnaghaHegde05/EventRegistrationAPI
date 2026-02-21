package event

// CreateNewEvent is unchanged
func CreateNewEvent(name string, seats int) error {
	return CreateEvent(Event{
		Name:       name,
		TotalSeats: seats,
	})
}

// Register now accepts a userEmail
func Register(eventID int, userEmail string) error {
	return RegisterForEventWithUser(eventID, userEmail)
}
