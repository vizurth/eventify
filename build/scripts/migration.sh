#!/bin/bash
source .env

# если хочешь брать host из переменных окружения, можно так:
export MIGRATION_DSN="host=postgres port=$POSTGRES_PORT dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

sleep 5  # ждем, чтобы postgres поднялся

goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v