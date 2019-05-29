package main

import (
	"fmt"
)

type LikeManager interface {
	Like(userID, mediaID string) error
	Unlike(userID, mediaID string) error
}

type PostLikeManager struct {
	client *redis.Client
}

func NewPostLikeManager(client *redis.Client) *PostLikeManager {
	return &PostLikeManager{client}
}

func (p *PostLikeManager) Like(userID, postID string) error {
	// ? Is Redis Set performant enough to keep track of the unique likes of the user?
	return p.client.SAdd(postID, userID).Err()
}

func (p *PostLikeManager) Like(userID, postID string) error {
	return p.client.SRem(postID, userID).Err()
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
	userID, postID = "john", "hello world"
	fmt.Println(likeManager.Like(userID, postID))
	fmt.Println(likeManager.Unlike(userID, postID))
}
