package redisfifo

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client = nil

func SetFifoRdb(r *redis.Client) {
	rdb = r
}

func FifoPush(ctx context.Context, list string, v interface{}) (interface{}, error) {
	return rdb.Do(ctx, "RPUSH", list, v).Result()
}

func FifoPop(ctx context.Context, list string) (interface{}, error) {
	val, err := rdb.Do(ctx, "BLPOP", list, 0).Result()
	if err != nil {
		return nil, err
	}

	return val.([]interface{})[1], nil
}
