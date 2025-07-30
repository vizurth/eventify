# build/docker/userinteract.Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o notification-service ./notification/cmd


CMD ["./notification-service"]
