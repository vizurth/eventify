package redis

import (
	"context"
	"eventify/common/logger"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Addr        string `yaml:"addr"`
	Password    string `yaml:"password"`
	User        string `yaml:"user"`
	DB          int    `yaml:"db"`
	MaxRetries  int    `yaml:"max_retries"`
	DialTimeout int    `yaml:"dial_timeout"`
	Timeout     int    `yaml:"timeout"`
}

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
	})

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	// retry
	var err error
	retries := 3
	for i := 0; i < retries; i++ {
		err = db.Ping(ctx).Err()
		if err == nil {
			break
		}
		log.Warn(ctx, "failed to ping redis, retrying...", zap.Error(err))
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		log.Error(ctx, "failed to connect to redis", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info(ctx, "successfully connected to redis", zap.String("addr", cfg.Addr), zap.Int("db", cfg.DB))

	return db, nil
}
