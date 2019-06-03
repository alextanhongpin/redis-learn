package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:       0,
		Password: "123456",
		Addr:     "localhost:6379",
	})
}

func main() {

	client := NewClient()

	res := client.XAdd(&redis.XAddArgs{
		Stream: "mystream",
		// MaxLen: "",
		// MaxLenApprox: "",
		// ID: "",
		Values: map[string]interface{}{
			"hello": "world",
		},
	})
	fmt.Println(res)
}
