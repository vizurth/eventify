create schema if not exists schema_name;

create table if not exists schema_name.events
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    city VARCHAR(100),
    venue VARCHAR(100),
    address VARCHAR(200),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    organizer_id INT,
    organizer_name VARCHAR(255),
    organizer_email VARCHAR(255),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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