-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id SERIAL PRIMARY KEY ,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table if not exists refresh_tokens (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    token text not null unique,
    expires_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
drop table if exists refresh_tokens;
-- +goose StatementEnd
