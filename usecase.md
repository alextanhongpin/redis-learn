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

## Game learderboards

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


## TODO

- What is russian doll caching?
## References

https://d0.awsstatic.com/whitepapers/performance-at-scale-with-amazon-elasticache.pdf#page40
