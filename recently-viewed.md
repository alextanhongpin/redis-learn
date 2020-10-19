## Recently Viewed

__Objective__: Implement top 10 recently viewed items with SortedSet.

We want to show users the last 10 items that they have just viewed, sorted from most recent to least recent. Everytime the user visited an item, their history will be recorded. If the user viewed the same item some time later, the score will be updated. If there are more than 10 items in the list, we will periodically prune them. We don't have to do it all the time, just conditionally remove them (20% probability).

To implement it with redis, we need several fields:

```
zadd [key] [score] [id]

key: The key to indicate the recently viewed item per user, e.g. history:jobs:john-doe
score: The negative timestamp, so that newer item will score lower. This makes sorting much easier with redis since scores are sorted from high to low.
id: The id of the item to store.
```

Sample items:
```
zadd job -1 job-1
zadd job -2 job-2
zadd job -3 job-3
zadd job -4 job-4
zadd job -5 job-5
```

Return the list of jobs from lowest score (latest date entry) to highest score.

```
# Take from the lowest score to the highest score (will return all items)
zrange job 0 -1

# Or if we just want to take the 10 lowest score. (This will return 10 items, from rank 0-9)
zrange job 0 9
```

Keep the last 3 item, and delete the rest.

```
# Remove all items starting from position 3 (in a list of 5 items, item 4 and 5 will be removed)
zremrangebyrank job 3 -1
```
