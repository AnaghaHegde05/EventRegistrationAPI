-- +goose Up
CREATE TABLE IF NOT EXISTS registrations (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_email TEXT NOT NULL,
    registered_at TIMESTAMP DEFAULT NOW()
);

