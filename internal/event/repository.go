package event

import (
	"errors"

	"event-registration-api/internal/db"
)

type Event struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TotalSeats     int    `json:"total_seats"`
	AvailableSeats int    `json:"available_seats"`
}

func CreateEvent(e Event) error {
	query := `
		INSERT INTO events (name, total_seats, available_seats)
		VALUES ($1, $2, $2)
	`
	_, err := db.DB.Exec(query, e.Name, e.TotalSeats)
	return err
}

func GetAllEvents() ([]Event, error) {
	rows, err := db.DB.Query("SELECT id, name, total_seats, available_seats FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Name, &e.TotalSeats, &e.AvailableSeats)
		events = append(events, e)
	}
	return events, nil
}

func RegisterForEvent(eventID int) error {
	result, err := db.DB.Exec(`
		UPDATE events
		SET available_seats = available_seats - 1
		WHERE id = $1 AND available_seats > 0
	`, eventID)

	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("event full")
	}
	return nil
}
