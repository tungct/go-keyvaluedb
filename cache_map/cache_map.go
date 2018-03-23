package cache_map

import (
	"unsafe"
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
	"fmt"
)

const sizeCache = 5 * 32
var maxLenghtCache int = sizeCache / ( int(unsafe.Sizeof("")) + int(unsafe.Sizeof("")) )

// init mem cache
func InitCacheMap() (c map[string] string){
	cache := make(map[string] string)
	return cache
}

// remove item in mem cache
func RemoveItemCache(cach *map[string] string, timeIte *map[string] int, countReques *map[string] int, key string, lenghtCach *int){
	cache := *cach
	timeItem := *timeIte
	countRequest := *countReques
	//lenghtCache := *lenghtCach
	delete(cache, key)
	delete(timeItem, key)
	delete(countRequest, key)
	fmt.Println("Cache lenght : ", *lenghtCach)
	*lenghtCach = *lenghtCach - 1
	fmt.Println("After, lenght cache : ", *lenghtCach)
}

// cache a item to mem cache, if over cache
func AddItemToCache(cache *map[string] string, timeItem *map[string] int, countRequest *map[string] int, key string, value string, lenghtCache *int, client *redis.Client, db *leveldb.DB){
	// if not exist item in cache
	if (*cache)[key] != "" {
		return
	}else {
		(*cache)[key] = value
		(*timeItem)[key] = int(time.Now().Unix())
		(*countRequest)[key] = 1
		*lenghtCache = *lenghtCache + 1
		return
	}
}

// get item in mem cache, if not found, get in redis
func GetItemCache(cache *map[string]string, timeItem *map[string] int, countRequest *map[string] int, key string, client * redis.Client){
	// if item in mem cache
	(*timeItem)[key] = int(time.Now().Unix())
	(*countRequest)[key] = (*countRequest)[key] + 1
}

func CheckItemInCache(cache *map[string]string, key string) bool{
	if _, ok := (*cache)[key]; ok {
		return true
	}else {
		return false
	}
}