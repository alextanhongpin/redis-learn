## Redis Learn

Basic Redis usage.

```bash
# Start docker.
$ docker-compose up -d

# Enter REPL
$ docker exec -it <REDIS_IMAGE_ID> redis-cli

# Shutdown redis.
$ docker-compose down

# Shortcuts
$ make up
$ make down
$ make redis

# Help
redis > help @string
redis > help set
```

## strings

Usecases:

- cache mechanism (SET, GET, MSET, MGET), e.g. cache text, binary data, html page, api responses, image, video up to 512MB.
- cache with automatic expiration (SETEX, EXPIRE, EXPIREAT), e.g. caching database queries.
- counting (INCR, INCRBY, DECR, DECRBY, INCRFLOATBY), e.g. page views, likes, voting system.

## list

Usecases:
- event queue: used by Resque, Celery, and Logstash
- storing most recent user posts: twitter stores the latest tweets of a user in a list
- timeline, last viewed item, recent purchases, recent notifications, recent chat

## hashes

How Instagram use Redis hash instead of strings to store the keys.
https://instagram-engineering.com/storing-hundreds-of-millions-of-simple-key-value-pairs-in-redis-1091ae80f74c

## sets

Usecases:

- data filtering: filtering all the flights that depart from a given city and arrive in another
- data grouping: grouping all users who viewed similar products, e.g. recommendations on Amazon
- membership checking: checking whether a user is on a blacklist
- location checking: checking if a user has been in the location before

You can also use bloomfilter for some of the stuff mentioned here if precision is not what you are looking for (this is especially true when the dataset starts to grow large).

## sorted sets

Usecases:
- build a real-time waiting list for customer service
- show a leaderboard of a massive online game that displays the top players, users with similar scores, or the scores of your friend
- build an autocomplete system using million of words

## bitmaps

Usecases:
- Useful for real-time data analytics, whether a user performed an action e.g. did user X performed Y today? How many times an event occured, e.g. how many users performed action Y this week?

## hyperloglog

Usecases:
- counting the number of unique users who visited the website 
- counting the number of distinct terms that were searched for on your website on a specific date or time
- counting the number of distinct hashtags that were used by a user
- counting the number of distinct words that appear in a book

# Time Series

- usage of specific words or terms in a newspaper over time
- minimum wage year-by-year
- daily changes in stock prices
- product purchased month-by-month
- climate changes

# commands

## pubsub

Usecases:
- news and weather dashboards
- chat application
- push notifications, such as subway delay alert
- remote code execution

## stream

## How to do bitops with redis
## Redis modules

## Good articles

https://tech.trivago.com/2017/01/25/learn-redis-the-hard-way-in-production/

## Cache Prewarm

http://dmitrypol.github.io/redis/2017/03/27/redis-cache-pregen.html
https://sorentwo.com/2015/07/27/optimizing-redis-usage-for-caching.html

## Optimization
https://aws.amazon.com/caching/implementation-considerations/
https://d0.awsstatic.com/whitepapers/performance-at-scale-with-amazon-elasticache.pdf#page40
http://highscalability.com/strategy-break-memcache-dog-pile

## Task Queues

https://charlesleifer.com/blog/multi-process-task-queue-using-redis-streams/
https://slack.engineering/scaling-slacks-job-queue-687222e9d100?gi=14c80a524901

