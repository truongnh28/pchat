package redis

import (
	"chat-app/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type IRedisClient interface {
	Set(
		ctx context.Context,
		key string,
		value interface{},
		expireTime time.Duration,
	) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	HDel(ctx context.Context, key, field string) *redis.IntCmd
	HSet(ctx context.Context, key, field string, values interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	LRem(ctx context.Context, key string, count int64, value string) *redis.IntCmd
	Del(ctx context.Context, key string) *redis.IntCmd
	Publish(ctx context.Context, channel string, message any) *redis.IntCmd
	Subscribe(ctx context.Context, channel []string) *redis.PubSub
	Ping(ctx context.Context) *redis.StatusCmd
	LRang(ctx context.Context, key string, startIndex, endIndex int64) *redis.StringSliceCmd
	SAdd(ctx context.Context, key string, value ...any) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SIsMember(ctx context.Context, key, value string) *redis.BoolCmd
	SRem(ctx context.Context, key, value string) *redis.IntCmd
}

type redisClientImpl struct {
	client *redis.Client
}

var redisClient *redisClientImpl
var redisOnce sync.Once

func GetRedisClient(redisConfig *config.RedisConfig) IRedisClient {
	redisOnce.Do(func() {
		redisClient = &redisClientImpl{
			client: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
				Password: redisConfig.Password,
			}),
		}
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Errorf("unable to connect to redis: %v", err.Error()))
		}
	})

	return redisClient
}

func (r *redisClientImpl) Ping(ctx context.Context) *redis.StatusCmd {
	return r.client.Ping(ctx)
}

func (r *redisClientImpl) Publish(
	ctx context.Context,
	channel string,
	message any,
) *redis.IntCmd {
	return r.client.Publish(ctx, channel, message)
}

func (r *redisClientImpl) Subscribe(ctx context.Context, channel []string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel...)
}

func (r *redisClientImpl) HSet(
	ctx context.Context,
	key, field string,
	value interface{},
) *redis.IntCmd {
	return r.client.HSet(ctx, key, field, value)
}

func (r *redisClientImpl) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, key, field)
}

func (r *redisClientImpl) Set(
	ctx context.Context,
	key string,
	value interface{},
	expireTime time.Duration,
) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expireTime)
}

func (r *redisClientImpl) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *redisClientImpl) HDel(ctx context.Context, key, field string) *redis.IntCmd {
	return r.client.HDel(ctx, key, field)
}

func (r *redisClientImpl) LPush(
	ctx context.Context,
	key string,
	values ...interface{},
) *redis.IntCmd {
	return r.client.LPush(ctx, key, values...)
}

func (r *redisClientImpl) RPush(
	ctx context.Context,
	key string,
	values ...interface{},
) *redis.IntCmd {
	return r.client.RPush(ctx, key, values...)
}

func (r *redisClientImpl) LPop(ctx context.Context, key string) *redis.StringCmd {
	return r.client.LPop(ctx, key)
}

func (r *redisClientImpl) Del(ctx context.Context, key string) *redis.IntCmd {
	return r.client.Del(ctx, key)
}

func (r *redisClientImpl) RPop(ctx context.Context, key string) *redis.StringCmd {
	return r.client.RPop(ctx, key)
}

func (r *redisClientImpl) LRem(
	ctx context.Context,
	key string,
	count int64,
	value string,
) *redis.IntCmd {
	return r.client.LRem(ctx, key, count, value)
}

func (r *redisClientImpl) LRang(
	ctx context.Context,
	key string,
	startIndex, endIndex int64,
) *redis.StringSliceCmd {
	return r.client.LRange(ctx, key, startIndex, endIndex)
}

func (r *redisClientImpl) SAdd(ctx context.Context, key string, value ...any) *redis.IntCmd {
	return r.client.SAdd(ctx, key, value)
}

func (r *redisClientImpl) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.SMembers(ctx, key)
}

func (r *redisClientImpl) SIsMember(ctx context.Context, key, value string) *redis.BoolCmd {
	return r.client.SIsMember(ctx, key, value)
}

func (r *redisClientImpl) SRem(ctx context.Context, key, value string) *redis.IntCmd {
	return r.client.SRem(ctx, key, value)
}
