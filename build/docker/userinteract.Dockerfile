# --- Stage 1: Build dependencies and binary ---
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем только модули на этом этапе для кэширования
COPY go.mod go.sum ./

# Кэшируем go mod download
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN go build -o userinteract-service ./user-interaction/cmd

EXPOSE 8083

CMD ["./userinteract-service"]
