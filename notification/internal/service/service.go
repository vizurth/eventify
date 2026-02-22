package service

import (
	"context"
	mykafka "eventify/common/kafka"
	"eventify/common/logger"
	"eventify/common/retry"
	"eventify/notification/internal/wsserver"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type NotificationService struct {
	kafkaReaders []*mykafka.Reader
	wsServer     wsserver.WsServer
	log          *logger.Logger
}

func NewNotificationService(
	ctx context.Context,
	kafkaCfg mykafka.Config,
	wsServer wsserver.WsServer,
	log *logger.Logger,
) *NotificationService {
	topics := []string{
		"event-created",
		"review-created",
		"review-updated",
		"review-deleted",
		"registration-created",
		"registration-deleted",
	}
	var readers []*mykafka.Reader
	for _, topic := range topics {
		readers = append(readers, mykafka.NewReader(ctx, kafkaCfg, topic, "notification-group"))
	}
	return &NotificationService{
		kafkaReaders: readers,
		wsServer:     wsServer,
		log:          log,
	}
}

func (s *NotificationService) Start(ctx context.Context) error {
	s.log.Info(ctx, "starting notification service")

	// Канал для сообщений

	// Запускаем чтение Kafka в горутине
	for _, reader := range s.kafkaReaders {
		msgCh := make(chan kafka.Message, 10)
		reader.StartConsuming(ctx, msgCh, retry.Strategy{Attempts: 5, Delay: time.Second, Backoff: 2})
		go s.consumeMessages(ctx, msgCh, reader)
	}

	// Запускаем обработку сообщений

	return nil
}

func topicTitle(topic string) string {
	titles := map[string]string{
		"event-created":        "Новое событие создано",
		"review-created":       "Новый отзыв",
		"review-updated":       "Отзыв обновлён",
		"review-deleted":       "Отзыв удалён",
		"registration-created": "Новая регистрация на событие",
		"registration-deleted": "Регистрация отменена",
	}
	if title, ok := titles[topic]; ok {
		return title
	}
	return topic
}

func (s *NotificationService) consumeMessages(ctx context.Context, msgCh <-chan kafka.Message, reader *mykafka.Reader) {
	for {
		select {
		case <-ctx.Done():
			s.log.Info(ctx, "stopping kafka consumer")
			return
		case msg, ok := <-msgCh:
			if !ok {
				s.log.Info(ctx, "kafka message channel closed")
				return
			}

			s.log.Info(ctx, "received kafka message",
				zap.String("topic", msg.Topic),
				zap.String("key", string(msg.Key)),
			)

			wsMsg := wsserver.NotificationMessage{
				Type:    msg.Topic,
				Title:   topicTitle(msg.Topic),
				Message: string(msg.Value),
			}

			if err := s.wsServer.BroadcastMessage(wsMsg); err != nil {
				s.log.Error(ctx, "failed to broadcast websocket message", zap.Error(err))
			} else {
				s.log.Info(ctx, "successfully broadcasted websocket message", zap.String("topic", msg.Topic))
			}

			// Подтверждаем обработку сообщения в Kafka
			if err := reader.Commit(ctx, msg); err != nil {
				s.log.Error(ctx, "failed to commit kafka message", zap.Error(err))
			}
		}
	}
}

func (s *NotificationService) Stop() error {
	for _, reader := range s.kafkaReaders {
		if err := reader.Close(); err != nil {
			s.log.Error(context.Background(), "failed to close kafka reader", zap.Error(err))
		}
	}
	return nil
}
