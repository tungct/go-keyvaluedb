package cache_map

import (
	"unsafe"
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tungct/go-keyvaluedb/redis_storage"
	"fmt"
)

const sizeCache = 5 * 32 // 10 kb
var maxLenghtCache int = sizeCache / ( int(unsafe.Sizeof("")) + int(unsafe.Sizeof("")) )

// init mem cache
func InitCacheMap() (c map[string] string){
	cache := make(map[string] string)
	return cache
}

// remove item in mem cache
func RemoveItemCache(cache map[string] string, key string, lenghtCache int) (c map[string] string, len int){
	delete(cache, key)
	lenghtCache --
	return cache, lenghtCache
}

// cache a item to mem cache, if over cache, cache to redis
func AddItemToCache(cache map[string] string, key string, value string, lenghtCache int, client *redis.Client, db *leveldb.DB) (c map[string] string, lenght int){
	if len(cache) < maxLenghtCache {
		// if not exist item in cache
		if cache[key] != "" {
			return cache, lenghtCache
		}
		cache[key] = value
		lenghtCache ++
		return cache, lenghtCache
	}else {
		redis_storage.SetKeyValueToRedis(client, key, value)
		fmt.Println("Cache full, cache to redis")
		return cache, lenghtCache
	}
}

// get item in mem cache, if not found, get in redis
func GetItemCache(cache map[string]string, key string, client * redis.Client) string{
	// if item in mem cache
	if val, ok := cache[key]; ok {
		return val
	}else { // else, get in redis
		val,_ := redis_storage.GetValueFromKeyRedis(client, key)
		return val
	}
}