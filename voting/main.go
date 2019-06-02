package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type Vote struct {
	President string
	Votes     int64
}

type VotingManager interface {
	GetVotes() ([]*Vote, error)
	GetVote(id string) (*Vote, error)
	VoteFor(id string) error
}

type VotingManagerImpl struct {
	client     *redis.Client
	presidents []string
}

func NewVotingManager(client *redis.Client, presidents []string) VotingManager {
	return &VotingManagerImpl{
		client:     client,
		presidents: presidents,
	}
}

func (v *VotingManagerImpl) VoteFor(president string) error {
	err := v.client.Incr(president).Err()
	return err
}

func (v *VotingManagerImpl) GetVotes() ([]*Vote, error) {
	var votes []*Vote
	for _, president := range v.presidents {
		vote, err := v.GetVote(president)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	return votes, nil
}

func (v *VotingManagerImpl) GetVote(president string) (*Vote, error) {
	// TODO: Validate if the president is a valid key.
	val, err := v.client.Get(president).Int64()
	if err != nil {
		return nil, err
	}
	return &Vote{
		President: president,
		Votes:     val,
	}, nil
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // No password set.
		DB:       0,        // Use default db.
	})
}

func main() {
	client := NewClient()
	votingManager := NewVotingManager(client, []string{"a", "b", "c"})
	_ = votingManager.VoteFor("a")
	_ = votingManager.VoteFor("a")
	_ = votingManager.VoteFor("b")
	_ = votingManager.VoteFor("c")
	_ = votingManager.VoteFor("c")
	votes, err := votingManager.GetVotes()
	if err != nil {
		log.Fatal(err)
	}
	for _, vote := range votes {
		fmt.Println(vote.President, vote.Votes)
	}
}
