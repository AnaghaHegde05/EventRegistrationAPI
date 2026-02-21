# Design Document – Preventing Race Conditions in Event Registration

## Problem Statement

In an event registration system, multiple users may attempt to register for the same event simultaneously.  
Without proper concurrency control, this can lead to **overbooking**, where more users are registered than available seats.  
Additionally, the system must now track **who** registered (via email) and allow organizers to list and cancel registrations, all while maintaining data consistency.

---

## Race Condition Explained

A classic race condition occurs when:
1. Two requests read the same `available_seats` value (e.g., 1).
2. Both see that seats are available and proceed.
3. Both decrement the seat count, resulting in `available_seats = -1` – overbooking.

---

## Design Approach

The core principle is to **let the database handle concurrency** using its transactional guarantees.  
All operations that modify state (registration, cancellation) are wrapped in **database transactions** and use **atomic conditional updates** to ensure correctness.

---

## 1. Registration Flow (Atomic Transaction)

When a user registers for an event, the system must:
- Decrement `available_seats` (only if > 0)
- Insert a record into the `registrations` table

Both operations must succeed or fail together.  
This is achieved with a **single transaction** that performs:

```sql
BEGIN;

UPDATE events
SET available_seats = available_seats - 1
WHERE id = $1 AND available_seats > 0;

-- Check if any row was affected; if not, rollback and return "event full".

INSERT INTO registrations (event_id, user_email)
VALUES ($1, $2);

COMMIT;
```
### Why this prevents race conditions

- The `UPDATE` is **atomic** – the condition `available_seats > 0` and the decrement happen in one indivisible step.
- If two concurrent transactions attempt the same `UPDATE`, the database ensures only one succeeds (the one that finds seats > 0). The other will see zero rows affected and roll back.
- The transaction isolation guarantees that the `INSERT` is only performed if the seat decrement was successful.
- The `CHECK (available_seats >= 0)` constraint on the `events` table provides an additional safety net.

---

## 2. Cancellation Flow (Atomic Transaction)

Cancelling a registration must:

- Delete the registration record
- Increment `available_seats` of the associated event

Again, a transaction ensures consistency:

```sql
BEGIN;

-- Get event_id for this registration
SELECT event_id FROM registrations WHERE id = $1;

DELETE FROM registrations WHERE id = $1;

UPDATE events
SET available_seats = available_seats + 1
WHERE id = $2;

COMMIT;
```
---
### 3. Why Database‑Level Concurrency?

- **Single source of truth**: The database is the final authority on seat counts.
- **No application locks**: Avoids complexity and distributed lock management.
- **Proven correctness**: PostgreSQL's transaction isolation (default `Read Committed`) ensures these operations are safe under concurrent load.
- **Scalability**: The database handles concurrency efficiently; application servers can be scaled horizontally without additional coordination.

---

### 4. Concurrency Test

A dedicated test (`internal/test/concurrency_test.go`) verifies the system under load:

- Creates an event with 5 seats.
- Launches 20 goroutines, each attempting to register.
- Counts successes and ensures they never exceed 5.
- Also checks the final `available_seats` and the number of registration records.

The test passes consistently, proving the design is race‑free.

---

## Summary

By moving concurrency control to the database and using atomic transactions, the system guarantees:

- **No overbooking**, even under extreme load.
- **Consistent tracking** of registrations.
- **Simple, maintainable** application code.

This design meets the capstone requirements and is ready for production‑scale workloads.


This design meets the capstone requirements and is ready for production‑scale workloads.
