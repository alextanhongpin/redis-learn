# To remove all keys matching a pattern in redis fast

```python
import time
import logging
from rediscluster import StrictRedisCluster


logger = logger.getLogger(__name__)
client = StrictRedisClient(startup_nodes=hosts,
                           password=password,
                           skip_full_coverage_check=True)

pattern = 'abc:*'
start_time = time.time()
item_count = 0
batch_size = 100_000
keys = []

logger.info('Start scanning keys...')

for k in client.scan_iter(pattern, count=batch_size):
  keys.append(k)
  if len(keys) >= batch_size:
    item_count += len(keys)
    logger.info(f'batch delete to {item_count}...')
    client.delete(*keys)
    keys = []

if len(keys) > 0:
  item_count += len(keys)
  logger.info(f'batch delete to {item_count}...')
  client.delete(*keys)

end_time = time.time()

logger.info(f'deleted {item_count} keys in {end_time-start_time:0.3f} ms')
```
