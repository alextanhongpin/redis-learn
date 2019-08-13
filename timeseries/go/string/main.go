package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type Granularity struct {
	Name     string // The level of the granularity.
	TTL      int64  // The duration to keep the data.
	Duration int64  // The time-window to store the data.
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
	"1min":  Granularity{"1min", 7 * 24 * Hour, Minute},
	"1hour": Granularity{"1hour", 60 * 24 * Hour, Hour},
	"1day":  Granularity{"1day", -1, Day},
}

type TimeSeries struct {
	client        *redis.Client
	namespace     string
	granularities map[string]Granularity
}

func NewTimeSeries(client *redis.Client, namespace string, granularities map[string]Granularity) *TimeSeries {
	if len(granularities) == 0 {
		granularities = defaultGranularities
	}
	return &TimeSeries{
		client:        client,
		namespace:     namespace,
		granularities: granularities,
	}
}

func (t *TimeSeries) Insert(timestampInSecs int64) {
	for _, granularity := range t.granularities {
		key := t.key(granularity, timestampInSecs)
		_ = t.client.Incr(key)
		if granularity.TTL > 0 {
			t.client.Expire(key, time.Duration(granularity.TTL)*time.Second)
		}
	}
}

func (t *TimeSeries) key(granularity Granularity, timestampInSecs int64) string {
	roundedTimestamp := t.roundedTimestamp(granularity, timestampInSecs)
	return fmt.Sprintf("%s:%s:%d", t.namespace, granularity.Name, roundedTimestamp)
}

func (t *TimeSeries) roundedTimestamp(granularity Granularity, timestampInSecs int64) int64 {
	return timestampInSecs - (timestampInSecs % granularity.Duration)
}

func (t *TimeSeries) Fetch(name string, startTimestamp, endTimestamp int64) ([]Series, error) {
	granularity, ok := t.granularities[name]
	if !ok {
		return nil, errors.New("granularity does not exist")
	}
	start := t.roundedTimestamp(granularity, startTimestamp)
	end := t.roundedTimestamp(granularity, endTimestamp)
	var keys []string
	for ts := start; ts <= end; ts += granularity.Duration {
		key := t.key(granularity, ts)
		keys = append(keys, key)
	}
	res, err := t.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]Series, len(res))
	for i := 0; i < len(res); i++ {
		var val int64
		if res[i] != nil {
			// Convert from interface to redis string.
			s, _ := res[i].(string)

			// Parse the string value into int64.
			var err error
			val, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		result[i] = Series{
			Timestamp: startTimestamp + int64(i)*granularity.Duration,
			Value:     val,
		}
	}
	return result, nil
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
	timeseries := NewTimeSeries(client, "go.srv/timeseries", nil)
	var startTimestamp int64
	timeseries.Insert(startTimestamp)
	timeseries.Insert(startTimestamp + 1)
	timeseries.Insert(startTimestamp + 1)
	timeseries.Insert(startTimestamp + 3)
	timeseries.Insert(startTimestamp + 61)
	{
		results, err := timeseries.Fetch("1sec", startTimestamp, startTimestamp+3)
		if err != nil {
			log.Fatal(err)
		}
		displayResults("1sec", results)

	}
	{
		results, err := timeseries.Fetch("1min", startTimestamp, startTimestamp+120)
		if err != nil {
			log.Fatal(err)
		}
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
