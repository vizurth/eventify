# build/docker/auth.Dockerfile
FROM golang:1.24 AS builder


WORKDIR /app

COPY . .

RUN go mod tidy && go build -o auth-service ./auth/cmd

EXPOSE 8081

CMD ["./auth-service"]
