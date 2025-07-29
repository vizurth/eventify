create schema if not exists schema_name;

create table if not exists schema_name.reviews
(
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES schema_name.events(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    username VARCHAR(100),
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);