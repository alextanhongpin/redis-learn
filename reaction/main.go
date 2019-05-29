package main

import (
	"fmt"
)

type ReactionManager interface {
	React(userID, mediaID, reactionType string) error
	Unreact(userID, mediaID, reactionType string) error
}

type PostReactionManager struct {
	client *redis.Client
}

func NewPostReactionManager(client *redis.Client) *PostReactionManager {
	return &PostReactionManager{client}
}

func postReaction(postID, reactionType string) string {
	// return strings.Join([]string{postID, reactionType}, ":")
	return fmt.Sprintf("%s:%s", postID, reactionType)
}

func (p *PostReactionManager) React(userID, postID, reactionType string) error {
	// ? Is Redis Set performant enough to keep track of the unique likes of the user?
	return p.client.SAdd(postReaction(postID, reactionType), userID).Err()
}

func (p *PostReactionManager) Unreact(userID, postID, reactionType string) error {
	return p.client.SRem(postReaction(postID, reactionType), userID).Err()
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
	reactionType = "happy" // sad, angry, amazed, like
	fmt.Println(reactionManager.React(userID, postID, reactionType))
	fmt.Println(reactionManager.Unreact(userID, postID, reactionType))
}
