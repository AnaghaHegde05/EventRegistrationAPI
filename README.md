# Event Registration & Ticketing System

**Capstone Project 5 – Go (Golang)**

## Overview

This project implements a backend REST API for an Event Registration & Ticketing System, similar to platforms like Eventbrite. Users can browse events, register for events with their email address, and organizers can create events and manage registrations (view and cancel).

The primary challenge solved here is handling concurrent registrations safely – ensuring that overbooking never occurs, even when multiple users attempt to register for the last available seats simultaneously. The system also now tracks who registered, allowing organizers to see the list of attendees and cancel registrations if needed.

---

## Key Features

- Create events with limited seat capacity
- Browse all available events with real‑time seat availability
- Register for an event by providing an email address – seat count decreases atomically and registration is stored
- List all registrations for a specific event (organizer view)
- Cancel a registration – automatically frees up a seat
- Prevent overbooking under concurrent requests using atomic database updates
- Automated concurrency test simulating many users booking the last seats simultaneously

---

## Tech Stack

- Language: Go (Golang)
- Web Framework: Gin
- Database: PostgreSQL
- Testing: Go testing framework (testing, sync)
- Concurrency Handling: Database‑level atomic operations with transactions

---

## Project Structure

The project follows a clean layered architecture:

- cmd/server/ – entry point, server setup
- internal/db/ – database connection (configurable via env)
- internal/event/ – handlers, service, repository
- internal/test/ – concurrent booking stress test
- migrations/ – SQL schema files for events and registrations tables

---

## Database Schema

The database consists of two tables with strict constraints to maintain data integrity.

- events table: stores event details with checks to ensure total and available seats are never negative.
- registrations table: links a user email to an event, with a foreign key to events.

Together, these tables allow full tracking of who registered for which event.

---

## API Endpoints

### 1. Create Event (Organizer)
- Method: POST
- Path: /events
- Request body: event name and number of seats
- Response: success message

### 2. Browse Events (User)
- Method: GET
- Path: /events
- Response: list of events with id, name, total seats, available seats

### 3. Register for Event (User)
- Method: POST
- Path: /events/:id/register
- Request body: user email
- Response: success message or error if event full

This endpoint atomically decrements available seats and inserts a registration record within a database transaction.

### 4. List Registrations for an Event (Organizer)
- Method: GET
- Path: /events/:id/registrations
- Response: list of registrations (id, event id, user email, timestamp)

### 5. Cancel a Registration (Organizer)
- Method: DELETE
- Path: /registrations/:id
- Response: success message

This endpoint deletes the registration and increments the event's available seats, all within a transaction.

---

## Concurrency Strategy (Core Requirement)

To prevent overbooking, seat allocation is handled using a single atomic SQL update inside a transaction that also inserts the registration. The update only succeeds if seats are still available. PostgreSQL guarantees transactional isolation, so concurrent requests cannot interfere. A CHECK constraint on available seats provides an extra safety layer, making overbooking impossible.

This approach avoids race conditions entirely without relying on application‑level locks.

---

## Concurrency Test

A dedicated test simulates multiple users attempting to register for the same event at the same time. It creates a fresh event, launches many goroutines each trying to register, and verifies that the number of successful registrations never exceeds the event's capacity. It also checks the final seat count and the number of registration records.

Run the test with: go test -v ./...

---

## How to Run the Project

1. Create a PostgreSQL database named eventdb.
2. Run the migration files in order to create the events and registrations tables.
3. Set the environment variable DB_PASSWORD to your PostgreSQL password. Optionally set DB_HOST, DB_PORT, DB_USER, DB_NAME (defaults are provided).
4. Start the server with: go run cmd/server/main.go
5. The API will be available at http://localhost:8080.

---

## Testing the API (PowerShell Examples)

- Create an event: use Invoke-RestMethod with a POST to /events and a JSON body containing name and seats.
- Browse events: use GET to /events.
- Register with email: POST to /events/:id/register with a JSON body containing user_email.
- List registrations: GET to /events/:id/registrations.
- Cancel a registration: DELETE to /registrations/:id.

---

## Conclusion

This project demonstrates:

- Correct REST API design with a clean layered architecture
- Strong database constraints to enforce data integrity
- Safe concurrency handling via atomic SQL updates and transactions
- Realistic concurrent testing using goroutines
- Full registration management – track who registered, list, and cancel

The system is now a complete Event Registration & Ticketing API that meets all capstone requirements and can serve as a foundation for a real‑world application.

---

## License

This project is open source and available under the MIT License.
