package cache

import (
	"chat-app/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

var redisClient *RedisClient
var redisOnce sync.Once

func GetRedisClient(redisConfig *config.RedisConfig) *RedisClient {
	redisOnce.Do(func() {
		redisClient = &RedisClient{
			client: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
				Password: redisConfig.Password,
			}),
		}
		if err := redisClient.client.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Errorf("unable to connect to redis: %v", err.Error()))
		}
	})
	return redisClient
}
