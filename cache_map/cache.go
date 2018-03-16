package cache_map

import (
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
)

// cache a item
func PutItemToCache(cache map[string] string, key string, value string, lenghtCache int, client *redis.Client, db *leveldb.DB)( c map[string] string, len int){
	cache, lenghtCache = AddItemToCache(cache, key, value, lenghtCache, client, db)
	return cache, lenghtCache
}

// get item in cache with key
func GetItemInCache(cache map[string]string, key string, client * redis.Client)string{
	return GetItemCache(cache, key, client)
}
