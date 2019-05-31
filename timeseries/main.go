package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Granularity struct {
	Name     string
	TTL      int64
	Duration int64
}

type Series struct {
	Timestamp, Value int64
}

var granularities = map[string]Granularity{
	"1sec":  Granularity{"1sec", 2 * time.Hour().Seconds(), time.Second.Seconds()},
	"1min":  Granularity{"1min", 7 * 24 * time.Hour().Seconds(), time.Minute.Seconds()},
	"1hour": Granularity{"1hour", 60 * 24 * time.Hour().Seconds(), time.Hour.Seconds()},
	"1day":  Granularity{"1day", -1, 24 * time.Hour.Seconds()},
}

type TimeSeries struct {
	client        *redis.Client
	namespace     string
	granularities map[string]Granularity
}

func (t *TimeSeries) Insert(timestampInSecs time.Duration) {
	for _, granularity := range t.granularities {
		key := t.key(granularity, timestampInSecs)
		_ = this.client.Incr(key)
		if granularity.TTL > 0 {
			this.client.Expire(key, granularity.TTL)
		}
	}
}

func (t *TimeSeries) key(granularity Granularity, timestampInSecs int64) {
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
		key := this.key(granularity, ts)
		keys = append(keys, key)
	}
	res, err := t.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]Series, len(res))
	for i := 0; i < len(res); i += 1 {
		val, _ := res[i].(int64)
		result[i] = startTimestamp + val*granularity.Duration
	}
	return result, nil
}

func main() {
}
