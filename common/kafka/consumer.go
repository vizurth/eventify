//package kafka
//
//import (
//	"context"
//	"github.com/segmentio/kafka-go"
//	"go.uber.org/zap"
//	"eventify/common/logger"
//)
//
//type KafkaConsumer struct {
//	reader *kafka.Reader
//}
//
//func NewConsumer(brokers []string, topic, groupId string) *KafkaConsumer {
//	reader := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:  brokers,
//		Topic:    topic,
//		GroupID:  groupId,
//		MinBytes: 1e3,
//		MaxBytes: 1e6,
//	})
//
//	return &KafkaConsumer{reader: reader}
//}
//
//func (c *KafkaConsumer) Consume(ctx context.Context, handler func(context.Context, []byte) error) {
//	for {
//		msg, err := c.reader.ReadMessage(ctx)
//		if err != nil {
//			logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to read message", zap.Error(err))
//			continue
//		}
//
//		if err := handler(ctx, msg.Value); err != nil {
//			logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to handle message", zap.Error(err))
//		}
//	}
//}
//
//func (c *KafkaConsumer) Close() error {
//	return c.reader.Close()
//}

package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) StartListening(ctx context.Context, handler func(key, value []byte)) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}
		log.Printf("received message key=%s value=%s", string(m.Key), string(m.Value))
		handler(m.Key, m.Value)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
