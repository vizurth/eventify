//package kafka
//
//import (
//	"context"
//	"github.com/segmentio/kafka-go"
//	"time"
//)
//
//type KafkaProducer struct {
//	writer *kafka.Writer
//}
//
//func NewProducer(brokers []string, topic string) *KafkaProducer {
//	return &KafkaProducer{
//		writer: &kafka.Writer{
//			Addr:         kafka.TCP(brokers...),
//			Topic:        topic,
//			Balancer:     &kafka.LeastBytes{},
//			RequiredAcks: kafka.RequireAll,
//		},
//	}
//}
//
//func (p *KafkaProducer) Send(ctx context.Context, msg *kafka.Message) error {
//	msgToSend := kafka.Message{
//		Key:   msg.Key,
//		Value: msg.Value,
//		Time:  time.Now(),
//	}
//
//	if err := p.writer.WriteMessages(ctx, msgToSend); err != nil {
//		//logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to send message to kafka", zap.Error(err))
//		return err
//	}
//	return nil
//}
//
//func (p *KafkaProducer) Close() error {
//	return p.writer.Close()
//}

package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("failed to write message: %v", err)
	}
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
