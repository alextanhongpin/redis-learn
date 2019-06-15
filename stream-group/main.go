package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:       0,
		Password: "123456",
		Addr:     "localhost:6379",
	})
}

func main() {
	var (
		consumerName = "john"
		groupName    = "mygroup"
		streamName   = "mystream"
		block        = 2 * time.Second
		count        = 10
		lastID       = "0-0"
		checkBacklog = false
	)
	client := NewClient()
	defer client.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var myid string
		if checkBacklog {
			myid = lastID
		} else {
			myid = ">"
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Pick the id based on the iteration: the first time we want to read our pending messages, in case we crashed and are recovering.
			// Once we consume our history, we can start getting new messages.
			items, err := client.XReadGroup(&redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamName, myid},
				Count:    int64(count),
				Block:    block,
				NoAck:    true,
			}).Result()
			if err != nil && err == redis.Nil {
				log.Println("read group failed", err)
				continue
			}
			fmt.Println("got items", items)
			if items == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// If we receive an empty reply, it means we were consuming our history
			checkBacklog = len(items[0].Messages) == 0
			// and that history is now empty. Let's start to consume new messages.
			for _, item := range items[0].Messages {
				fmt.Println(item.ID, item.Values)
				client.XAck(streamName, groupName, item.ID)
				lastID = item.ID
			}
		}
	}()
	<-stop
	cancel()
	wg.Wait()
	fmt.Println("terminating application")
}
