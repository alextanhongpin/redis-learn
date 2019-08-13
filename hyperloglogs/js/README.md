## HyperLogLogs

We can use HyperLogLogs to store the unique visits of users to the website. The HyperLogLogs consumes a lot less storage than the equivalent Sets when storing user ids (UUID).


```
PFADD visits:2019-01-01 a b c d
PFADD visits:2019-01-02 a b c
```

To count the visits on a particular date:

```
PFCOUNT visits:2019-01-01
```

If we specify more than 1 key when getting the count, it will return the union of counts:

```
PFCOUNT visits:2019-01-01 visits:2019-01-02
```

We can also store the union of the results in a different destination:

```
PFMERGE visits:total visits:2019-01-01 visits:2019-01-02
PFCOUNT visits:total
```
