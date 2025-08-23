package kafka

import (
	"context"
	"eventify/common/logger"
	"eventify/common/retry"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"strings"
)

type Config struct {
	Host    string   `yaml:"host" env:"HOST" env-default:"kafka"`
	Port    uint16   `yaml:"port" env:"PORT" env-default:"9094"`
	Brokers []string `yaml:"brokers" env:"BROKERS" env-separator:","`
}

type Writer struct {
	*kafka.Writer
}

type Reader struct {
	*kafka.Reader
}

func NewReader(ctx context.Context, cfg Config, topic, groupID string) *Reader {
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
	return &Reader{r}
}

func (c *Reader) Fetch(ctx context.Context) (kafka.Message, error) {
	return c.Reader.FetchMessage(ctx)
}

func (c *Reader) Commit(ctx context.Context, msg kafka.Message) error {
	return c.Reader.CommitMessages(ctx, msg)
}

func (c *Reader) Close() error {
	return c.Reader.Close()
}

func (c *Reader) FetchWithRetry(ctx context.Context, strat retry.Strategy) (kafka.Message, error) {
	var msg kafka.Message
	err := retry.Do(func() error {
		m, e := c.Fetch(ctx)
		if e == nil {
			msg = m
		}
		return e
	}, strat)
	return msg, err
}

func (c *Reader) StartConsuming(ctx context.Context, out chan<- kafka.Message, strat retry.Strategy) {
	go func() {
		defer close(out)
		for {
			msg, err := c.FetchWithRetry(ctx, strat)
			if err != nil {
				// Можно добавить логирование ошибки
				break
			}
			select {
			case out <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
}

func NewWriter(ctx context.Context, cfg Config, topic string) *Writer {
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
	return &Writer{w}
}

func (p *Writer) Send(ctx context.Context, key, value []byte) error {
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (p *Writer) Close() error {
	return p.Writer.Close()
}

func (p *Writer) SendWithRetry(ctx context.Context, strat retry.Strategy, key, value []byte) error {
	return retry.Do(func() error {
		return p.Send(ctx, key, value)
	}, strat)
}
func CreateTopicIfNotExists(cfg Config, topic string, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		// Проверка на уже существующий топик
		if strings.Contains(err.Error(), "Topic with this name already exists") ||
			strings.Contains(err.Error(), "TOPIC_ALREADY_EXISTS") {
			// Просто логируем и не считаем ошибкой
			return nil
		}
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

func CreateTopicWithRetry(cfg Config, topic string, numPartitions, replicationFactor int, strat retry.Strategy) error {
	return retry.Do(func() error {
		return CreateTopicIfNotExists(cfg, topic, numPartitions, replicationFactor)
	}, strat)
}
