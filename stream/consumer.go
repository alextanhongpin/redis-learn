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
	res := client.XRead(&redis.XReadArgs{
		Streams: []string{"mystream", "$"},
		Count:   0,
		Block:   0,
	})
	fmt.Println(res.Result())
}
