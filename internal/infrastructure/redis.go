package infrastructure

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

var Rdb *redis.Client

// NewRedisClient 根据配置创建 Redis 连接，并实现 ping 检测
func NewRedisClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Configs.Redis.Addr,
		Password: config.Configs.Redis.Password,
		DB:       config.Configs.Redis.DB,
	})

	// Ping 检测
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal(nil, "ping redis error: "+err.Error())
		panic(err)
	}
	Rdb = rdb
}
