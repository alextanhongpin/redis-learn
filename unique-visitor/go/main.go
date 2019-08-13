package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
}

type UniqueCounter struct {
	client *redis.Client
	keyFmt func(string) string
}

func NewUniqueCounter(client *redis.Client, keyFmt func(string) string) *UniqueCounter {
	return &UniqueCounter{
		client: client,
		keyFmt: keyFmt,
	}
}

func (u *UniqueCounter) AddVisit(date, userID string) error {
	return u.client.PFAdd(u.keyFmt(date), userID).Err()
}

func (u *UniqueCounter) GetVisits(date string) int64 {
	return u.client.PFCount(u.keyFmt(date)).Val()
}

func (u *UniqueCounter) AggregateDate(dest string, dates ...string) (int64, error) {
	key := u.keyFmt(dest)
	formattedDates := make([]string, len(dates))
	for i, date := range dates {
		formattedDates[i] = u.keyFmt(date)
	}
	err := u.client.PFMerge(key, formattedDates...).Err()
	if err != nil {
		return -1, err
	}
	return u.client.PFCount(key).Val(), nil
}

func main() {
	var (
		maxUsers    = 200
		totalVisits = 1000
	)

	client := NewClient()
	keyFmt := func(key string) string {
		return fmt.Sprintf("go.srv/visits/daily/%s", key)
	}
	counter := NewUniqueCounter(client, keyFmt)
	for i := 0; i < totalVisits; i++ {
		username := fmt.Sprintf("user_%d", rand.Intn(maxUsers))
		date := fmt.Sprintf("%s:%d", time.Now().Format("2006-01-02"), rand.Intn(24))
		_ = counter.AddVisit(date, username)
	}
	allHours := make([]string, 24)
	for i := 0; i < 24; i++ {
		allHours[i] = fmt.Sprintf("%s:%d", time.Now().Format("2006-01-02"), rand.Intn(24))
	}
	for _, hr := range allHours {
		counts := counter.GetVisits(hr)
		fmt.Printf("visits found in %s is %d\n", hr, counts)
	}
	total, err := counter.AggregateDate("total", allHours...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("total aggregate counter is", total)
}
