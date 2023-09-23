package main

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueHandler struct {
	Cache CacheStorage
}

func NewQueueHandler(c CacheStorage) QueueHandler {
	return QueueHandler{c}
}

func (handler QueueHandler) UserAddedHandler(message amqp.Delivery) error {
	newUser := UserResponse{}

	if err := json.Unmarshal(message.Body, &newUser); err != nil {
		fmt.Println("failed to process message", message.Body, err.Error())
		return err
	}

	users := []UserResponse{}
	cacheResult, err := handler.Cache.Get(context.Background(), USERS_CACHE_KEY)
	if err != nil {
		fmt.Println("failed to get cache users", err.Error(), cacheResult)
	}
	json.Unmarshal([]byte(cacheResult), &users)

	users = append(users, newUser)
	val, err := json.Marshal(users)
	if err != nil {
		fmt.Println("failed to marshall users", err.Error())
	}
	handler.Cache.Set(context.Background(), USERS_CACHE_KEY, string(val))

	return nil
}
