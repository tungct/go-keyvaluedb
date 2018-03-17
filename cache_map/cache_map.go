package cache_map

import (
	"unsafe"
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

const sizeCache = 5 * 32
var maxLenghtCache int = sizeCache / ( int(unsafe.Sizeof("")) + int(unsafe.Sizeof("")) )

// init mem cache
func InitCacheMap() (c map[string] string){
	cache := make(map[string] string)
	return cache
}

// remove item in mem cache
func RemoveItemCache(cache map[string] string, timeItem map[string] int, countRequest map[string] int, key string, lenghtCache int) (c map[string] string, len int){
	delete(cache, key)
	delete(timeItem, key)
	delete(countRequest, key)
	lenghtCache --
	return cache, lenghtCache
}

// cache a item to mem cache, if over cache
func AddItemToCache(cache map[string] string, timeItem map[string] int, countRequest map[string] int, key string, value string, lenghtCache int, client *redis.Client, db *leveldb.DB) (c map[string] string, lenght int){
	// if not exist item in cache
	if cache[key] != "" {
		return cache, lenghtCache
	}
	cache[key] = value
	timeItem[key] = int(time.Now().Unix())
	countRequest[key] = 1
	lenghtCache ++
	return cache, lenghtCache
}

// get item in mem cache, if not found, get in redis
func GetItemCache(cache map[string]string, timeItem map[string] int, countRequest map[string] int, key string, client * redis.Client){
	// if item in mem cache
	timeItem[key] = int(time.Now().Unix())
	countRequest[key] = countRequest[key] + 1
}