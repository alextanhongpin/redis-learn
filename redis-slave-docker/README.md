## Redis Master-Slave

Run docker:

```
$ docker-compose up -d
Creating network "redis-slave-docker_redis-network" with driver "bridge"
Creating redis-slave-docker_redis-master_1   … done
Creating redis-slave-docker_redis-slave-02_1 ...
Creating redis-slave-docker_redis-slave-02_1 … done
Creating redis-slave-docker_redis-slave-01_1 … done
```

Verify that a redis master and two slaves are created:

```bash
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                  PORTS                       NAMES
c9da735bd594        redis:5.0.5         "docker-entrypoint.s…"   2 seconds ago       Up Less than a second   127.0.0.1:32814->6379/tcp   redis-slave-docker_redis-slave-01_1
b38feeef83e8        redis:5.0.5         "docker-entrypoint.s…"   2 seconds ago       Up Less than a second   127.0.0.1:32813->6379/tcp   redis-slave-docker_redis-slave-02_1
fce7e26c3520        redis:5.0.5         "docker-entrypoint.s…"   2 seconds ago       Up 1 second             127.0.0.1:6379->6379/tcp    redis-slave-docker_redis-master_1
```

## Check replication

Connect to the master and run the command `info replication`. We can see that both slaves are connected to the master instance.

```bash
$ docker exec -it fce redis-cli
127.0.0.1:6379> auth 123456
OK
127.0.0.1:6379> info replication
# Replication
role:master
connected_slaves:2
slave0:ip=192.168.32.3,port=6379,state=online,offset=28,lag=0
slave1:ip=192.168.32.4,port=6379,state=online,offset=28,lag=1
master_replid:c7db13e6d12a5d537d0640589ba1e5fb0912320a
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:28
second_repl_offset:-1
repl_backlog_active:1
repl_backlog_size:1048576
repl_backlog_first_byte_offset:1
repl_backlog_histlen:28
127.0.0.1:6379>
```

## Write to master

Set some values in the master:
```
127.0.0.1:6379> set name john
OK
127.0.0.1:6379> set age 100
OK
```
   
## Read from slave

Connect to any of the slave. We can read the values set from the master, but
cannot set any values here.
```bash
$ docker exec -it b38 redis-cli
127.0.0.1:6379> get name
"john"
127.0.0.1:6379> get age
"100"
127.0.0.1:6379> set name 100
(error) READONLY You can't write against a read only replica.
127.0.0.1:6379>
```
