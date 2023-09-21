package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStorage struct {
	Client *redis.Client
}

func NewCacheStorage(addr string) CacheStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	return CacheStorage{Client: rdb}
}

func (cache CacheStorage) Set(ctx context.Context, key string, value any) {
	err := cache.Client.Set(ctx, key, value, 1*time.Hour).Err()
	if err != nil {
		panic(err)
	}
}

func (cache CacheStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
