package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/CodingFervor/carbon-emission-management/internal/config"
	"github.com/CodingFervor/carbon-emission-management/pkg/logger"
)

var RDB *redis.Client

func Connect(cfg config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	logger.Info("redis connected", "addr", cfg.Addr())
	return nil
}

func Close() {
	if RDB != nil {
		RDB.Close()
		logger.Info("redis connection closed")
	}
}
