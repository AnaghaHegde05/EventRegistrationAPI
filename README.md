# Event Registration & Ticketing System

**Capstone Project 5 – Go (Golang)**

## Overview

This project implements a backend REST API for an **Event Registration & Ticketing System**, similar to platforms like Eventbrite.
Users can browse events, register for events with limited capacity, and organizers can create events.

The **primary focus** of this project is handling **concurrent registrations safely**, ensuring that **overbooking never occurs**, even when multiple users attempt to register for the last available seats simultaneously.

---

## Key Features

* Create events with limited seat capacity
* Browse all available events
* Register users for events
* Prevent overbooking under concurrent requests
* Automated concurrency test simulating multiple users

---

## Tech Stack

* **Language:** Go (Golang)
* **Web Framework:** Gin
* **Database:** PostgreSQL
* **Testing:** Go testing framework (`testing`, `sync`)
* **Concurrency Handling:** Database-level atomic operations

---

## Project Structure

```
event-registration-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── event/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   └── test/
│       └── concurrency_test.go
├── migrations/
│   └── 001_create_events.sql
├── README.md
├── DESIGN.md
├── go.mod
└── go.sum
```

---

## Database Schema

The system uses a single `events` table with strict constraints to maintain data integrity.

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    total_seats INT NOT NULL CHECK (total_seats > 0),
    available_seats INT NOT NULL CHECK (available_seats >= 0),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## API Endpoints

### 1. Create Event (Organizer)

**POST /events**

Request Body:

```json
{
  "name": "Tech Conference",
  "seats": 5
}
```

Response:

```json
{
  "message": "event created"
}
```

---

### 2. Browse Events (User)

**GET /events**

Response:

```json
[
  {
    "id": 1,
    "name": "Tech Conference",
    "total_seats": 5,
    "available_seats": 5
  }
]
```

This endpoint fulfills the requirement:
**“Users can browse events.”**

---

### 3. Register for Event (User)

**POST /events/:id/register**

Success Response:

```json
{
  "message": "registration successful"
}
```

Failure Response (Event Full):

```json
{
  "error": "event full"
}
```

---

## Concurrency Strategy (Core Requirement)

To prevent overbooking, seat allocation is handled using a **single atomic SQL update**:

```sql
UPDATE events
SET available_seats = available_seats - 1
WHERE id = $1 AND available_seats > 0;
```

### Why this works:

* The availability check and seat decrement happen in **one atomic operation**
* PostgreSQL guarantees transactional safety
* Concurrent requests cannot reduce seats below zero
* Overbooking is impossible

This approach avoids race conditions without relying on application-level locks.

---

## Concurrency Test

A concurrency test simulates multiple users attempting to register for the same event at the same time.

* Uses Go **goroutines**
* Uses `sync.WaitGroup`
* Counts successful registrations
* Fails if successful registrations exceed available seats

Run the test:

```bash
go test ./...
```

Successful output:

```
ok   event-registration-api/internal/test
```

This proves the system is safe under concurrent load.

---

## How to Run the Project

### 1. Create Database

```sql
CREATE DATABASE eventdb;
```

### 2. Run Migration

```bash
psql -U postgres -d eventdb -f migrations/001_create_events.sql
```

### 3. Set Database Password

```bash
setx DB_PASSWORD "your_postgres_password"
```

Restart the terminal after setting the variable.

### 4. Start Server

```bash
go run cmd/server/main.go
```

Server runs on:

```
http://localhost:8080
```

---

## Testing Without Postman

You can test APIs using **PowerShell**.

Create Event:

```powershell
Invoke-RestMethod `
  -Uri http://localhost:8080/events `
  -Method POST `
  -ContentType "application/json" `
  -Body '{ "name": "Tech Conference", "seats": 5 }'
```

Browse Events:

```powershell
Invoke-RestMethod http://localhost:8080/events
```

Register for Event:

```powershell
Invoke-RestMethod `
  -Uri http://localhost:8080/events/1/register `
  -Method POST
```

---

## Conclusion

This project demonstrates:

* Correct REST API design
* Strong database constraints
* Safe concurrency handling
* Real concurrent testing using Go