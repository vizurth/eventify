package main

import (
	"context"
	"eventify/common/kafka"
	"eventify/notification/internal/handler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("[notification-service] starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("[notification-service] shutting down")
		cancel()
	}()

	// Init service & handler
	h := handler.NewNotificationHandler()

	// Kafka consumer (общий топик events)
	consumer := kafka.NewConsumer(
		[]string{"kafka:9092"},
		"events",
		"notification-service",
	)
	defer consumer.Close()

	// Запуск Kafka listener
	go consumer.StartListening(ctx, func(key, value []byte) {
		eventType := string(key)
		h.HandleNotification(eventType, value)
	})

	<-ctx.Done()
	log.Println("[notification-service] stopped")
}
