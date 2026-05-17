package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// 全局 Redis 客户端
var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// 测试连接
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Panicln("Redis 连接失败:", err)
	}
	log.Println("Redis 连接成功")
}
