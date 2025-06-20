create schema if not exists schema_name;

create table if not exists schema_name.event_participants
(
    event_id INT NOT NULL REFERENCES schema_name.events(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    username VARCHAR(255) NOT NULL,
    PRIMARY KEY (event_id, user_id)
);
