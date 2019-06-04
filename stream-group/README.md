## redis stream consumer

A stream needs to be created first:

```bash
$ xadd mystream * sensor-id 1 temperature 100 
$ xgroup create mystream mygroup $
```
