## Bitmap

We can use bitmap in Redis to store unique visitors. Bitmaps can store up to 2^32 bits (more than 4 billion).


We can only set 0 or 1 for SETBIT. To record that a user has visited a site on a particular date:
```
SETBIT visits:2019-01-01 <id> 1
```


To check whether the user has visited the site on a particular date:
```
GETBIT visits:2019-01-01 <id>
```

To count the number of visits on a particular date:

```
BITCOUNT visits:2019-01-01
```

To get the total count of users that visited the site on multiple dates:

```
BITOP OR total_users visits:2019-01-01 visits:2019-01-02
BITCOUNT total_users
```

