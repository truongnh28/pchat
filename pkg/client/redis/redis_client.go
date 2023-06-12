package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	HSet(ctx context.Context, key, field string, values interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	Set(
		ctx context.Context,
		key string,
		value interface{},
		expireTime time.Duration,
	) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	HDel(ctx context.Context, key, field string) *redis.IntCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, key string) *redis.IntCmd
}

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) HSet(
	ctx context.Context,
	key, field string,
	value interface{},
) *redis.IntCmd {
	return r.client.HSet(ctx, key, field, value)
}

func (r *RedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, key, field)
}

func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	expireTime time.Duration,
) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expireTime)
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *RedisClient) HDel(ctx context.Context, key, field string) *redis.IntCmd {
	return r.client.HDel(ctx, key, field)
}

func (r *RedisClient) LPush(
	ctx context.Context,
	key string,
	values ...interface{},
) *redis.IntCmd {
	return r.client.LPush(ctx, key, values...)
}

func (r *RedisClient) RPush(
	ctx context.Context,
	key string,
	values ...interface{},
) *redis.IntCmd {
	return r.client.RPush(ctx, key, values...)
}

func (r *RedisClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	return r.client.LPop(ctx, key)
}

func (r *RedisClient) Del(ctx context.Context, key string) *redis.IntCmd {
	return r.client.Del(ctx, key)
}
