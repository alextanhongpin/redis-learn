## Round robin

Naive approach of distributing cache keys. A hash function (such as CRC32) is used to compute the key, and modulo of the number of cache node is used to distribute the key to a random node in the list.

```
cache node list = [
	redishost1,
	redishost2
]

cache index = hash(key) % length(cache node list)
cache_node = cache_node_list[cache_index]
```

This approach has a fatal flaw. When we add or remove new nodes, the modulo of the key will cause the keys to be remapped to a new node.

http://www.tom-e-white.com/2007/11/consistent-hashing.html

## Caching Strategy

### Lazy caching

A.k.a lazy population or cache-aside. The cache is populated only when an object is requested by the application.

```python
def get_user(user_id):
	# Check the cache.
	record = cache.get(user_id)
	if record is None:
		# Run a db query.
		record = db.query(‘SELECT * FROM user WHERE id = ?’, user_id)
		# Populate the cache.
		cache.set(user_id, record)
	return record

# App code.
user = get_user(17)
```

### Write on through

The cache is updated in real-time when the database is updated. This prevents unnecessary cache-miss.

```python
def save_user(user_id, values):
	# Save to db.
	record = db.query(‘UPDATE user … WHERE id = ?’, user_id, values)
	
	# Push into cache
	cache.set(user_id, record)
	return record
	
# App code.
user = save_user(17, {“name”: “John”})
```


## Thundering herd

Also known as dog piling.

## Ranking

- View leaderboard with the most purchased items all-time
- Use sorted sets to handle the ranking
- ZRANGE returns items in ascending score, to get the top ten, use ZREVRANGE
```

ZINCRBY orders:items:popular 2 12345678
ZREVRANGE orders:items:popular 0 9
```

## Game leaderboards

We can use Sorted Sets to store the ranking for each user
```
ZADD leaderboard 556 Andy
ZADD leaderboard 819 Barry
ZADD leaderboard 105 Carl
ZADD leaderboard 1312 Derek

ZREVRANGE leaderboard 0 -1
1. Derek
2. Barry
3. Andy
4. Carl

ZRANK leaderboard barry
2
```
## Recommendation Engine

Redis counters can be used to increment or decrement the number of likes/dislikes for a given item. Redis HASHes can be used to maintain a list of everyone who has liked or disliked that item
```
INCR item:123456:likes
HSET item:123456:ratings Susan 1
INCR item:123456:dislikes
HSET item:123456:ratings Tommy -1
```
## Chat and Messaging

Pub-sub can be used for simple chat and messaging needs. Use cases include in-app messaging, web chat window, online game invites and chat, and real-time comment streams.
```
SUBSCRIBE chat:123456
PUBLISH chat:114 “Hello world” {“message”, “chat:114”, “hello all”}
UNSUBSCRIBE chat:114
```
## Queues

We can use LIST as queues.

Build a real-time sales analytics dashboard with Redis


## Daily order count

- Each time an order is received, we increment the count of orders by one
```
INCR orders:123456
GET orders:123456
```
## Number of daily unique items sold
```
SADD orders:items:20180101 itemskuid
EXPIRE orders:items:20180101 60 * 60 * 24 * 7
SCARD orders:items:20180101
```

## History

## Latest product purchased

- Use LIST to track a running list of the last 100 items sold
- Alternatively we can use sorted sets too, with the unix timestamp as the score
- After pushing new item, we trim the length to a total of 100
```
LPUSH orders:items:latest 123456
LTRIM orders:items:latest 0 99
```
To retrieve the last 10 item sold from the head, we use LRANGE:
```
LRANGE orders:items:latest 0 9
```
For fetching the next page:
```
LRANGE orders:items:latest 10 19
```

## Historical Sales Revenue

- View one week of historical data, tracked on an hourly basis
- The dates and the revenue is tracked in two separate data structure
- Sorted Set is used to store the revenue data, whereby each member is unique
- Member’s score in sorted set is set to 0, so the results are sorted by the member key, which is ascending in date
- ZREMRANGEBYRANK is used to prune unused data

```
ZADD sales:revenue:days 0 20180101
ZREMRANGEBYRANK sales:revenue:days 0 7
```
We increment the hourly revenue in Redis HASH.
- Hashes are used to store the mapping between a key (the current hour) and a value (revenue)
- For each date we keep a separate hash of data, expiring it after seven days
```
HINCRBYFLOAT sales:revenue:20180101 10 244.56
EXPIRES sales:revenu:20180101 60*60*24*7
```
To retrieve sales data for an entire week
- Retrieve the list of dates for which we have data from a sorted set
- Iterate through each date, and retrieve the hourly revenue from the hash
```
ZREVRANGE sales:revenue:days 0 6
HGETALL sales:revenue:20180101
```
Pseudocode:
```
var dates ZRANGE sales:revenue:days 0 6
var revenue

foreach date in dates:
	revenue[data] = HGETALL sales:revenue: + data
```

- What is russian doll caching?
## References

https://d0.awsstatic.com/whitepapers/performance-at-scale-with-amazon-elasticache.pdf#page40
