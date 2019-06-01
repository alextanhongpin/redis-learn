package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

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
	sub := client.Subscribe("hello")
	res, err := sub.Receive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	for {
		select {
		case res := <-sub.Channel():
			fmt.Println(res)
		}
	}
}
