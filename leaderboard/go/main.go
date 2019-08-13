package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

type Leaderboard struct {
	client *redis.Client
	key    string
}

func NewLeaderboard(client *redis.Client, key string) *Leaderboard {
	return &Leaderboard{
		client: client,
		key:    key,
	}
}

func (l *Leaderboard) AddUser(username string, score float64) error {
	return l.client.ZAdd(l.key, redis.Z{Member: username, Score: score}).Err()
}

func (l *Leaderboard) RemoveUser(username string) error {
	return l.client.ZRem(l.key, username).Err()
}

func (l *Leaderboard) GetUserScoreAndRank(username string) (score float64, rank int64) {
	score = l.client.ZScore(l.key, username).Val()
	rank = l.client.ZRevRank(l.key, username).Val() + 1
	return
}

func (l *Leaderboard) ShowTopUsers(quantity int64) ([]redis.Z, error) {
	fmt.Printf("showing the top %d users\n", quantity)
	return l.client.ZRevRangeWithScores(l.key, 0, quantity-1).Result()
}

func main() {

	client := NewClient()
	leaderboard := NewLeaderboard(client, "go.srv/leaderboard")
	fmt.Println(leaderboard.AddUser("alpha", 10))
	leaderboard.AddUser("beta", 101)
	leaderboard.AddUser("charlie", 75)
	leaderboard.AddUser("delta", 25)
	leaderboard.AddUser("eclair", 55)

	leaderboard.RemoveUser("eclair")
	score, rank := leaderboard.GetUserScoreAndRank("beta")
	log.Println("beta score and rank is", score, rank)

	users, err := leaderboard.ShowTopUsers(3)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Println(user.Member, user.Score)
	}
}
