package cache_map

import (
	"unsafe"
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

const sizeCache = 200000 * 32
var maxLenghtCache int = sizeCache / ( int(unsafe.Sizeof("")) + int(unsafe.Sizeof("")) )

// init mem cache
func InitCacheMap() (c map[string] string){
	cache := make(map[string] string)
	return cache
}

// remove item in mem cache
func RemoveItemCache(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, countRequest *map[string] int, key string, lenghtCache *int){

	//lenghtCache := *lenghtCach
	delete(*cache, key)
	delete(*timeItemCache, key)
	delete(*timeItemInit, key)
	delete(*countRequest, key)
	*lenghtCache = *lenghtCache - 1
}

// cache a item to mem cache, if over cache
func AddItemToCache(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, countRequest *map[string] int, key string, value string, lenghtCache *int, client *redis.Client, db *leveldb.DB){
	// if not exist item in cache
	if (*cache)[key] != "" {
		return
	}else {
		(*cache)[key] = value
		timeNow := int(time.Now().UnixNano())
		(*timeItemInit)[key] = timeNow
		(*timeItemCache)[key] = timeNow
		(*countRequest)[key] = 1
		*lenghtCache = *lenghtCache + 1
		return
	}
}

// get item in mem cache, if not found, get in redis
func GetItemCache(cache *map[string]string, timeItemCache *map[string] int, countRequest *map[string] int, key string, client * redis.Client){
	// if item in mem cache
	(*timeItemCache)[key] = int(time.Now().UnixNano())
	(*countRequest)[key] = (*countRequest)[key] + 1
}

func CheckItemInCache(cache *map[string]string, key string) bool{
	if _, ok := (*cache)[key]; ok {
		return true
	}else {
		return false
	}
}