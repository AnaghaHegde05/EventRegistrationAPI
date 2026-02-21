package test

import (
	"sync"
	"testing"

	"event-registration-api/internal/db"
	"event-registration-api/internal/event"
)

func TestConcurrentBooking(t *testing.T) {
	// 🔹 Step 1: Connect to database explicitly for tests
	err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// IMPORTANT:
	// Make sure an event with ID = 1 exists and has limited seats
	eventID := 1

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	totalRequests := 20

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := event.Register(eventID)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// This assertion assumes the event has only 5 seats
	if successCount > 5 {
		t.Fatalf(
			"Overbooking occurred! Expected max 5 successful bookings, got %d",
			successCount,
		)
	}
}
