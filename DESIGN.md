
---

# `DESIGN.md`

```markdown
# Design Document – Preventing Race Conditions in Event Registration

## Problem Statement

In an event registration system, multiple users may attempt to register for the same event at the same time.  
If the system does not handle concurrency correctly, this can lead to **overbooking**, where more users are registered than available seats.

---

## Race Condition Explained

A race condition occurs when:
- Multiple requests read the same available seat count
- Each request assumes seats are still available
- All requests proceed and update the database independently

Example (Incorrect approach):
1. Request A reads available_seats = 1
2. Request B reads available_seats = 1
3. Both register successfully
4. Result: available_seats = -1 (overbooking)

---

## Design Approach

This project prevents race conditions by **delegating concurrency control to the database**, which is the single source of truth.

Instead of:
- Reading available seats
- Checking in application code
- Updating later

We perform **check + update in one atomic SQL operation**.

---

## Atomic Update Strategy

```sql
UPDATE events
SET available_seats = available_seats - 1
WHERE id = $1 AND available_seats > 0;