package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

type Queue struct {
	client  *redis.Client
	key     string
	timeout time.Duration
}

func NewQueue(key string, client *redis.Client) *Queue {
	return &Queue{
		client:  client,
		key:     key,
		timeout: 0,
	}

}

func (q *Queue) Size() int64 {
	return q.client.LLen(q.key).Val()
}

func (q *Queue) Push(value string) error {
	return q.client.LPush(q.key, value).Err()
}

func (q *Queue) Pop() []string {
	return q.client.BRPop(q.timeout, q.key).Val()
}

func main() {
	client := NewClient()
	queue := NewQueue("go.srv/queue", client)
	fmt.Println("size is", queue.Size())
	fmt.Println("push", queue.Push("hello world"))
	fmt.Println("pop", queue.Pop())
}
