package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const (
	Second = 1
	Minute = 60
	Hour   = 60 * 60
	Day    = 24 * 60 * 60
)

type Series struct {
	Timestamp, Value int64
}

type Granularity struct {
	Name     string
	TTL      int64
	Duration int64
	Quantity int64
}

// These values are based on the hash-max-ziplist-entries of 512 so that data
// will be stored in the memory-optimized ziplist.
var defaultGranularities = map[string]Granularity{
	// Stores 300 timestamps of 1 second each.
	"1sec": Granularity{"1sec", 2 * Hour, Second, 5 * Minute},
	// Stores 480 timestamps of 1 minute each.
	"1min": Granularity{"1min", 7 * Day, Minute, 8 * Hour},
	// Stores 240 timestamps of 1 hour each.
	"1hour": Granularity{"1hour", 60 * Day, Hour, 10 * Day},
	// Stores a maximum of 30 timestamps of 1 day each.
	"1day": Granularity{"1day", -1, Day, 30 * Day},
}

type TimeSeries struct {
	namespace     string
	client        *redis.Client
	granularities map[string]Granularity
}

func NewTimeSeries(client *redis.Client, namespace string) *TimeSeries {
	return &TimeSeries{
		namespace:     namespace,
		client:        client,
		granularities: defaultGranularities,
	}
}

func (t *TimeSeries) key(granularity Granularity, timestampInSeconds int64) string {
	roundedTimestamp := t.roundTimestamp(timestampInSeconds, granularity.Quantity)
	return fmt.Sprintf("%s:%s:%d", t.namespace, granularity.Name, roundedTimestamp)
}

func (t *TimeSeries) roundTimestamp(timestampInSeconds, precision int64) int64 {
	return timestampInSeconds - (timestampInSeconds % precision)
}

func (t *TimeSeries) Insert(timestampInSeconds int64) {
	for _, granularity := range t.granularities {
		key := t.key(granularity, timestampInSeconds)
		field := strconv.FormatInt(t.roundTimestamp(timestampInSeconds, granularity.Duration), 10)
		fmt.Println(key, field)
		t.client.HIncrBy(key, field, 1)
		if granularity.TTL > 0 {
			t.client.Expire(key, time.Duration(granularity.TTL)*time.Second)
		}
	}
}

func (t *TimeSeries) Fetch(granularityName string, startTimestamp, endTimestamp int64) ([]Series, error) {
	granularity, ok := t.granularities[granularityName]
	if !ok {
		return nil, errors.New("granularity does not exist")
	}
	start := t.roundTimestamp(startTimestamp, granularity.Duration)
	end := t.roundTimestamp(endTimestamp, granularity.Duration)
	fmt.Println("start, end", start, end)

	var result []*redis.StringCmd
	_, err := t.client.Pipelined(func(pipe redis.Pipeliner) error {
		for ts := start; ts <= end; ts += granularity.Duration {
			key := t.key(granularity, ts)
			field := strconv.FormatInt(t.roundTimestamp(ts, granularity.Duration), 10)
			fmt.Println("key, field", key, field)
			result = append(result, pipe.HGet(key, field))
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		return nil, errors.Wrap(err, "pipeline error")
	}
	output := make([]Series, len(result))
	for i, res := range result {
		val := res.Val()
		var count int64
		if val != "" {
			var err error
			count, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
		}

		ts := startTimestamp + int64(i)*granularity.Duration
		output[i] = Series{Timestamp: ts, Value: count}
	}
	return output, nil
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	client := NewClient()
	timeseries := NewTimeSeries(client, "go.srv/timeseries")
	var startTimestamp int64
	timeseries.Insert(startTimestamp)
	timeseries.Insert(startTimestamp + 1)
	timeseries.Insert(startTimestamp + 1)
	timeseries.Insert(startTimestamp + 3)
	timeseries.Insert(startTimestamp + 61)
	{
		fmt.Println("operation 1")
		results, err := timeseries.Fetch("1sec", startTimestamp, startTimestamp+3)
		if err != nil {
			log.Fatal(err)
		}
		displayResults("1sec", results)

	}
	{
		fmt.Println("operation 2")
		results, err := timeseries.Fetch("1min", startTimestamp, startTimestamp+120)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("displaying results")
		displayResults("1min", results)
	}
}

func displayResults(granularityName string, results []Series) {
	fmt.Println("result from ", granularityName)
	for _, result := range results {
		fmt.Println(result.Timestamp, result.Value)
	}
	fmt.Println()
}
