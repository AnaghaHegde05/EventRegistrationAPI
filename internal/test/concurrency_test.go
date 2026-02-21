package test

import (
	"sync"
	"testing"

	"event-registration-api/internal/db"
	"event-registration-api/internal/event"
)

func TestConcurrentBooking(t *testing.T) {
	err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	eventName := "Concurrency Test Event"
	totalSeats := 5
	err = event.CreateNewEvent(eventName, totalSeats)
	if err != nil {
		t.Fatalf("failed to create test event: %v", err)
	}

	var eventID int
	err = db.DB.QueryRow(`
		SELECT id FROM events 
		WHERE name = $1 
		ORDER BY id DESC LIMIT 1
	`, eventName).Scan(&eventID)
	if err != nil {
		t.Fatalf("failed to get event ID: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	attempts := 20

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Use a dummy email; in a real test you might vary it
			err := event.Register(eventID, "test@example.com")
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successCount > totalSeats {
		t.Fatalf("❌ Overbooking! Expected ≤ %d successful registrations, got %d",
			totalSeats, successCount)
	}

	var available int
	err = db.DB.QueryRow("SELECT available_seats FROM events WHERE id = $1", eventID).Scan(&available)
	if err != nil {
		t.Fatalf("failed to query final available seats: %v", err)
	}
	expectedAvailable := totalSeats - successCount
	if available != expectedAvailable {
		t.Errorf("available_seats mismatch: expected %d, got %d", expectedAvailable, available)
	}

	// Optional: verify number of registration records
	var regCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM registrations WHERE event_id = $1", eventID).Scan(&regCount)
	if err != nil {
		t.Fatalf("failed to count registrations: %v", err)
	}
	if regCount != successCount {
		t.Errorf("registration count mismatch: expected %d, got %d", successCount, regCount)
	}

	t.Logf("✅ Test passed: %d successful registrations", successCount)
}
