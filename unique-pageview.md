# Unique pageviews with hyperloglog

To track unique page views, we can use hyperloglog algorithm to minimize storage and still get accurate results.
We need the following field to be store:
```
pfadd [key] [item]
key: A unique key to indicate the item to store, and also the time up to hours (we store hourly records, to get daily record we can merge them or count all the hourly ones), e.g. to store a job page view, we use the key `pv:job:2020010123` where `pv` is page view, and the number is date in the format `yymmddhh`.
item: The id of the item to store. In this case, it will be the job id.
```

Let's add some entry:
```
pfadd job-hr00 jobid-1
pfadd job-hr00 jobid-2
pfadd job-hr01 jobid-1
pfadd job-hr01 jobid-3
```

To get the count of the entry at hour 0 and hour 1:
```
pfcount job-hr00 // Returns 2.
pfcount job-hr01 // Returns 2.
```

We can get the combined count by specifying multiple keys. To get daily count, we just need to specify hours from 0-23.
```
pfcount job-hr00 job-hr01 // Returns 3 (unique instances).
```

We can merge the keys into a single key at the end of the day.
```
pfmerge job-hrall job-hr00 job-hr01 // Merge all keys to a new key.
```
