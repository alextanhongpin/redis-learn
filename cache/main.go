package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

const GetUsersCacheKey = "ServerName:UserGetter:GetUsers"

type User struct {
	Name string
	Age  int64
}

type Repository interface {
	GetUsers() ([]User, error)
}

type UserGetter struct {
}

func NewRepository() *UserGetter {
	return &UserGetter{}
}

func (u *UserGetter) GetUsers() ([]User, error) {
	return nil, nil
}

type Cache struct {
	client *redis.Client
	repo   Repository
}

func NewCache(client *redis.Client, repo Repository) *Cache {
	return &Cache{
		client: client,
		repo:   repo,
	}
}

func (c *Cache) GetUsers() ([]User, error) {
	// Check the cache for values.
	b, err := c.client.Get(GetUsersCacheKey).Bytes()
	if err == redis.Nil || len(b) == 0 {
		// Get and set if it does not exist.
		res, err := c.repo.GetUsers()
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		// Expires in 1 minute.
		err = c.client.Set(GetUsersCacheKey, string(b), 1*time.Minute).Err()
		return res, err
	}
	// Return if exists.
	var res []User
	err = json.Unmarshal(b, &res)
	return res, err
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	client := NewClient()
	repo := NewRepository()
	repoWithCache := NewCache(client, repo)
	repoWithCache.GetUsers()
}
