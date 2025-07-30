# build/docker/event.Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o event-service ./event/cmd

EXPOSE 8082

CMD ["./event-service"]
