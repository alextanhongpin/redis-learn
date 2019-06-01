package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Granularity struct {
	Name          string
	TTL, Duration int64
}

type Series struct {
	Timestamp, Value int64
}

const (
	Second = 1
	Minute = 60
	Hour   = 60 * 60
	Day    = 24 * 60 * 60
)

var defaultGranularities = map[string]Granularity{
	"1sec":  Granularity{"1sec", 2 * Hour, Second},
	"1min":  Granularity{"1min", 7 * Day, Minute},
	"1hour": Granularity{"1hour", 60 * Day, Hour},
	"1day":  Granularity{"1day", -1, Day},
}

type TimeSeries struct {
	client        *redis.Client
	namespace     string
	granularities map[string]Granularity
}

func NewTimeSeries(client *redis.Client, namespace string) *TimeSeries {
	return &TimeSeries{
		client:        client,
		namespace:     namespace,
		granularities: defaultGranularities,
	}
}

func (t *TimeSeries) Insert(timestampInSeconds int64, id string) error {
	for _, granularity := range t.granularities {
		key := t.key(granularity, timestampInSeconds)
		if err := t.client.PFAdd(key, id).Err(); err != nil {
			return err
		}
		if granularity.TTL > 0 {
			if err := t.client.Expire(key, time.Duration(granularity.TTL)*time.Second).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TimeSeries) key(granularity Granularity, timestampInSeconds int64) string {
	roundedTimestamp := t.roundTimestamp(timestampInSeconds, granularity.Duration)
	return fmt.Sprintf("%s:%s:%d", t.namespace, granularity.Name, roundedTimestamp)
}

func (t *TimeSeries) roundTimestamp(timestampInSeconds, precision int64) int64 {
	return timestampInSeconds - (timestampInSeconds % precision)
}

func (t *TimeSeries) Fetch(granularityName string, startTimestamp, endTimestamp int64) ([]Series, error) {
	granularity := t.granularities[granularityName]
	start := t.roundTimestamp(startTimestamp, granularity.Duration)
	end := t.roundTimestamp(endTimestamp, granularity.Duration)

	var result []*redis.IntCmd
	_, err := t.client.Pipelined(func(pipe redis.Pipeliner) error {
		for ts := start; ts <= end; ts += granularity.Duration {
			key := t.key(granularity, ts)
			result = append(result, pipe.PFCount(key))
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		return nil, errors.Wrap(err, "pipeline failed")
	}
	output := make([]Series, len(result))
	for i, res := range result {
		val := res.Val()
		output[i] = Series{
			Timestamp: startTimestamp + int64(i)*granularity.Duration,
			Value:     val,
		}
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
	// NOTE: Using / will cause an error here.
	timeseries := NewTimeSeries(client, "go:srv:timeseries")
	var startTimestamp int64
	log.Println(timeseries.Insert(startTimestamp, "user:max"))
	timeseries.Insert(startTimestamp, "user:max")
	timeseries.Insert(startTimestamp+1, "user:hugo")
	timeseries.Insert(startTimestamp+1, "user:renata")
	timeseries.Insert(startTimestamp+3, "user:hugo")
	timeseries.Insert(startTimestamp+61, "user:kc")
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
