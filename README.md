# go-keyvaluedb

- Client put and get message key-value from server by grpc protocol
- Server put message to mem-cache (map key-value), if cache full, put message to Redis, LevelDb
- Worker in Server check message in server, if message have time out, delete it
- Worker check message in redis or levelDb, if it have hight-frequency, put it to cache (or redis) by threshold

## Libs 

- go-Redis :  https://github.com/go-redis/redis
- go-LevelDb : https://github.com/syndtr/goleveldb/leveldb
- decode-config : https://github.com/BurntSushi/toml

## Config Redis and Server 

```
config/redis.conf
cofig/server.conf
```

## Getting started
run server : 

```
server/server.go
$ go run server.go
```

run test put client with 10 goroutine in during 5 time second : 
```
client/putclient

$go run putclient.go
```

run test get client with 10 goroutine in during 5 time second : 
```
client/getclient

$go run getclient.go
```

