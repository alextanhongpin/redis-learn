
https://d0.awsstatic.com/whitepapers/performance-at-scale-with-amazon-elasticache.pdf#page40


# Build a real-time sales analytics dashboard with Redis


## Daily order count

- Each time an order is received, we increment the count of orders by one

```
INCR orders:123456
GET orders:123456
```

## Number of daily unique items sold

SADD orders:items:20180101 itemskuid
EXPIRE orders:items:20180101 60 * 60 * 24 * 7
SCARD orders:items:20180101


## Leaderboard

- View leaderboard with the most purchased items all-time
- Use sorted sets to handle the ranking
- ZRANGE returns items in ascending score, to get the top ten, use ZREVRANGE

```
ZINCRBY orders:items:popular 2 12345678
ZREVRANGE orders:items:popular 0 9
```

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
	revenue[data] = HGETALL "sales:revenue:" + data
```
