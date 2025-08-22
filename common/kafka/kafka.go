package kafka

import (
	"context"
	"eventify/common/logger"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Host    string   `yaml:"host" env:"HOST" env-default:"kafka"`
	Port    uint16   `yaml:"port" env:"PORT" env-default:"9092"`
	Brokers []string `yaml:"brokers" env:"BROKERS" env-separator:","`
}

func NewReader(ctx context.Context, cfg Config, topic, groupID string) *kafka.Reader {
	l := logger.GetOrCreateLoggerFromCtx(ctx)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	l.Info(ctx, "connected to Kafka topic",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
		zap.String("group_id", groupID),
	)
	return r
}

func NewWriter(ctx context.Context, cfg Config, topic string) *kafka.Writer {
	l := logger.GetOrCreateLoggerFromCtx(ctx)
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
	}

	l.Info(ctx, "created Kafka writer",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("topic", topic),
	)
	return w
}

func CreateTopicIfNotExists(cfg Config, topic string, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return err
	}

	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
}

func CreateTopicWithRetry(cfg Config, topic string, numPartitions, replicationFactor int) error {
	var err error
	for i := 0; i < 10; i++ {
		err = CreateTopicIfNotExists(cfg, topic, numPartitions, replicationFactor)
		if err == nil {
			return nil
		}

		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i))
	}
	return err
}
