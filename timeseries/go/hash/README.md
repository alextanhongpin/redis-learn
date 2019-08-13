# Hash 

Using hash versus string. Say we create 1 events every second, so in 1 day we would have 87,865 keys in redis:

```
86,400 keys for 1sec granularity 
1,440 keys for 1min granularity
24 keys for 1hour granularity
1 key for 1day granularity
```


The size when storing in redis string vs hash is 11MB vs 800kb.

Small hashes are encoded in a different data structure called ziplist, which is memory optimized. When any of the conditions below are exceeded, it will be converted to a hash table:

- hash-max-ziplist-entries. The default value is 512 entries.
- hash-max-ziplist-values. The default size is 64 bytes.
