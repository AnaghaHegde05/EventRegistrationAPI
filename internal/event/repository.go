package event

import (
	"database/sql"
	"errors"
	"time"

	"event-registration-api/internal/db"
)

type Event struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TotalSeats     int    `json:"total_seats"`
	AvailableSeats int    `json:"available_seats"`
}

type Registration struct {
	ID           int       `json:"id"`
	EventID      int       `json:"event_id"`
	UserEmail    string    `json:"user_email"`
	RegisteredAt time.Time `json:"registered_at"`
}

// CreateEvent inserts a new event (unchanged)
func CreateEvent(e Event) error {
	query := `
		INSERT INTO events (name, total_seats, available_seats)
		VALUES ($1, $2, $2)
	`
	_, err := db.DB.Exec(query, e.Name, e.TotalSeats)
	return err
}

// GetAllEvents returns all events (unchanged)
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

// RegisterForEventWithUser atomically decrements seats and inserts a registration
func RegisterForEventWithUser(eventID int, userEmail string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Atomic seat decrement
	result, err := tx.Exec(`
		UPDATE events
		SET available_seats = available_seats - 1
		WHERE id = $1 AND available_seats > 0
	`, eventID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("event full")
	}

	// Insert registration
	_, err = tx.Exec(`
		INSERT INTO registrations (event_id, user_email)
		VALUES ($1, $2)
	`, eventID, userEmail)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetRegistrationsByEventID returns all registrations for a given event
func GetRegistrationsByEventID(eventID int) ([]Registration, error) {
	rows, err := db.DB.Query(`
		SELECT id, event_id, user_email, registered_at
		FROM registrations
		WHERE event_id = $1
		ORDER BY registered_at
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regs []Registration
	for rows.Next() {
		var r Registration
		err := rows.Scan(&r.ID, &r.EventID, &r.UserEmail, &r.RegisteredAt)
		if err != nil {
			return nil, err
		}
		regs = append(regs, r)
	}
	return regs, nil
}

// CancelRegistration removes a registration and increments available seats
func CancelRegistration(regID int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the event_id for this registration
	var eventID int
	err = tx.QueryRow("SELECT event_id FROM registrations WHERE id = $1", regID).Scan(&eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("registration not found")
		}
		return err
	}

	// Delete the registration
	_, err = tx.Exec("DELETE FROM registrations WHERE id = $1", regID)
	if err != nil {
		return err
	}

	// Increment available seats
	_, err = tx.Exec(`
		UPDATE events
		SET available_seats = available_seats + 1
		WHERE id = $1
	`, eventID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
