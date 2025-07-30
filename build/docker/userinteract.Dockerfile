# build/docker/userinteract.Dockerfile
FROM golang:1.24 AS builder


WORKDIR /app

COPY . .

RUN go mod tidy && go build -o userinteract-service ./user-interaction/cmd

EXPOSE 8083

CMD ["./userinteract-service"]
