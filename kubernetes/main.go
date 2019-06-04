package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:       0,
		Password: "",
		Addr:     "localhost:31567",
	})
}

func main() {
	client := NewClient()
	fmt.Println(client.Ping().Result())
}
