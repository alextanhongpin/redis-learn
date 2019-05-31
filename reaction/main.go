package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Users can react to a post.
// Users can choose one of the given reaction.
// Users can only react once for a post.
// Users can unreact to a post.
// System can return the total number of reactions of a post.
// System can return the total number of count for each reactions. (?)

type ReactionManager interface {
	React(mediaID, userID, reactionType string) error
	Unreact(mediaID, userID string) error
	Count(mediaID string) int64
}

type PostReactionManager struct {
	client *redis.Client
}

func NewPostReactionManager(client *redis.Client) *PostReactionManager {
	return &PostReactionManager{client}
}

func (p *PostReactionManager) React(postID, userID, reactionType string) error {
	// ? Is Redis Set performant enough to keep track of the unique likes of the user?
	return p.client.HSet(postID, userID, reactionType).Err()
}

func (p *PostReactionManager) Unreact(postID, userID string) error {
	return p.client.HDel(postID, userID).Err()
}

func (p *PostReactionManager) Count(postID string) int64 {
	return p.client.HLen(postID).Val()
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
	reactionManager := NewPostReactionManager(client)

	userID := "john"
	postID := "hello world"
	reactionType := "happy" // sad, angry, amazed, like

	{
		fmt.Println(reactionManager.React(postID, userID, reactionType))
		// Second time react with sad emoji.
		fmt.Println(reactionManager.React(postID, userID, "sad"))
	}

	// We can only get the total count of the posts.
	fmt.Println(reactionManager.Count(postID))
	fmt.Println(reactionManager.Unreact(postID, userID))
	fmt.Println(reactionManager.Count(postID))
}
