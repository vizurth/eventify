# Этап сборки
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости
COPY . .
RUN go mod download

# Копируем код и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/main.go

# Этап исполнения
#FROM alpine:latest
#
#RUN apk --no-cache add ca-certificates
#
#WORKDIR /root/
#
#COPY --from=builder /app/service .

EXPOSE 8081 8082 8083 8084

CMD ["./service"]
