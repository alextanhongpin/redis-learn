package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:       0,
		Password: "123456",
		Addr:     ":6379",
	})
}
func main() {
	// redisdb := NewClient()
	// NOTE: This will not work in MacOS, cause the IP returned is pointing
	// inside docker.
	redisdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master01",
		SentinelAddrs: []string{"127.0.0.1:26379"},
	})
	fmt.Println("pinging")
	res := redisdb.Ping()
	fmt.Println(res.Result())
	fmt.Println("getting result")
	fmt.Println(redisdb.Get("name").Result())
}
