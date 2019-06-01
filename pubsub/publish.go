package main

import "github.com/go-redis/redis"

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:       0,
		Password: "",
		Addr:     "localhost:6379",
	})
}

func main() {
	client := NewClient()
	defer client.Close()
	client.Publish("hello", "world")
}
