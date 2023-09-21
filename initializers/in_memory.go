package initializers

import (
	"context"
	"dev_community_server/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

var InMemoryDB InMemoryDatabase

type InMemoryDatabase interface {
	Set(key string, value string) error
	SetExp(key string, value string, expires time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type Redis struct {
	client *redis.Client
}

func InMemoryConnection() {
	r := &Redis{}
	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	InMemoryDB = r
}

func (r *Redis) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *Redis) SetExp(key string, value string, expires time.Duration) error {
	return r.client.Set(context.Background(), key, value, expires).Err()
}

func (r *Redis) Get(key string) (string, error) {
	str, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", utils.ErrValueNotFound
	}
	return str, err
}

func (r *Redis) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
