package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

type MediaID string

type LikeManager interface {
	Like(mediaID, userID string) error
	Unlike(mediaID, userID string) error
	Count(mediaID string) int64
}

type PostLikeManager struct {
	client *redis.Client
}

func NewPostLikeManager(client *redis.Client) *PostLikeManager {
	return &PostLikeManager{client}
}

func (p *PostLikeManager) Like(postID, userID string) error {
	// ? Is Redis Set performant enough to keep track of the unique likes of the user?
	return p.client.SAdd(postID, userID).Err()
}

func (p *PostLikeManager) Unlike(postID, userID string) error {
	return p.client.SRem(postID, userID).Err()
}

func (p *PostLikeManager) Count(postID string) int64 {
	return p.client.SCard(postID).Val()
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
	likeManager := NewPostLikeManager(client)
	userID, postID := "john", "hello world"
	fmt.Println(likeManager.Like(postID, userID))
	// Same user will be excluded.
	fmt.Println(likeManager.Like(postID, userID))
	fmt.Println(likeManager.Count(postID))
	fmt.Println(likeManager.Unlike(postID, userID))
	fmt.Println(likeManager.Count(postID))
}
