# go-keyvaluedb
libs :

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

