package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
}

type Analytics struct {
	client *redis.Client
}

func NewAnalytics(client *redis.Client) *Analytics {
	return &Analytics{
		client: client,
	}
}

func (a *Analytics) key(date string) string {
	return fmt.Sprintf("visits:daily:%s", date)
}
func (a *Analytics) AddUser(date string, userID int64) error {
	return a.client.SetBit(a.key(date), userID, 1).Err()
}

func (a *Analytics) Count(date string) int64 {
	return a.client.BitCount(a.key(date), &redis.BitCount{Start: 0, End: -1}).Val()
}

func (a *Analytics) ShowUserIDsFromVisit(date string) ([]int, error) {
	b, err := a.client.Get(a.key(date)).Bytes()
	if err != nil {
		return nil, err
	}
	var result []int
	for i := 0; i < len(b); i++ {
		for j := 7; j >= 0; j-- {
			if visited := b[i] >> (uint(i) & 1); visited == 1 {
				userID := i*8 + (7 - j)
				result = append(result, userID)
			} else {
				fmt.Println("visited is", visited)
			}
		}
	}
	return result, nil
}

func main() {
	client := NewClient()
	analytics := NewAnalytics(client)
	date := time.Now().Format("2006-01-02")
	analytics.AddUser(date, 1)
	analytics.AddUser(date, 2)
	fmt.Println(analytics.Count(date))
	// fmt.Println(analytics.ShowUserIDsFromVisit(date))
}
