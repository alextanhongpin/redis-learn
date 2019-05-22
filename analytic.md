## Types of Analytics

Transactional
- authorization
- authentication
- price management
- advertising bids
- messaging
- location based processing
- user session management

Analytics
- counting 
- leaderboard
- scores and rating
- page ranking 
- recommendation engine
- time-series analytics
- session analysis
- secondary index
- personalized offers
- location based ads

Operational
- accelerated reporting
- real-time attribution
- search
- order history
- inventory tracking

## Data structure use cases

- Hyperloglog: probabilistic estimates of counts for anomaly detection
- Sorted Sets: real-time range analyses, top scorers, bid rangers
- Sets: cardinality for fraud detection
- Geospatial indexes: location based searches
- Bitmaps: real-time population counting for activity monitoring


## Useful redis modules

Inline-analytics:
- Redis-cell: rate limiting with Generic cell rate algorithm
- rebloom: probabilistic membership queries
- t-digest: rank based statistic estimator
- topk: track the top-k most frequent elements in a stream
- countminsketch: approximate frequency counter

Operational analytics with modules
- redisearch: Extremely fast text based search, used for secondary indexing
- redis graph: graph query processing
- time-series: range analyses, build in aggregation (min, max, sum, avg)


Machine learning/deep learning modules
- neural redis
- redis-ml

## Fraud mitigation example

User geographic check
- is the user logging from a place that seems rational? Someone who lives in the US suddenly logs in from Mongolia
- determine rationality
  - simple: Check if the person ever visited the place before. If not, then require additional verification
  - advanced: Check bordering regions (someone from the US can login from Mexico or Canada)
- how it works:
  - every valid login location is added to the bloom filter
  - before login, check to see if the location for the login attempts is in the bloom filter
  - properties of bloom filter means that you'll know with certainty if a user has never been to this location
    - false positive is possible, but controllable
  - rebloom module


User out-of-hours check
- people generally follow patterns in life
  - morning: check bank account, top up transit pass
  - work: employee dashboard, partner reports
  - evening: food delivery, games, dating app
- how it works:
  - record the time period of each valid login in a bloom filter
  - check surrounding times in the bloom filter on each login
  - can check multiple dimensions, time of day, day of week, holidays
  
Resource traffic anomaly
- if one resource is being hit frequently by a small number of actors, this indicates problem like
  - poorly behaving clients
  - fraudulent usage
  - DDOS/Bots
- you need to record unique actors and raw counts for a specific time period
- a high ratio raw counts to unique actors indicates a traffic anomaly
- how it works:
  - key based on time period
  - record each visit to in a hyperloglog unique for each
    - hyperloglog estimate cardinality in a very small space with a standard error rate
  - increment a counter for each visit to each resource
  
## References
- https://www.slideshare.net/RedisLabs/collaborative-realtime-editing-shane-carr?next_slideshow=1
