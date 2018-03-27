# go-keyvaluedb
Thư viện sử dụng : 
- go-Redis :  https://github.com/go-redis/redis
- go-LevelDb : https://github.com/syndtr/goleveldb/leveldb

## Config Redis 

```
redis_storage/redis.go
```

## Getting started
run server : 

```
keyvalue_storage
$ go run server.go
```

run test put client with 10 goroutine in during 5 time second : 
```
keyvalue_storage
$go run client.go
```
