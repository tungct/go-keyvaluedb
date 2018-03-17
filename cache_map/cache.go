package cache_map

import (
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
	"fmt"
	"github.com/tungct/go-keyvaluedb/redis_storage"
)

const timeOut  = 10 // time second

// cache a item
func PutItemToCache(cache map[string] string, timeItem map[string] int, countRequest map[string]int, key string, value string, lenghtCache int, client *redis.Client, db *leveldb.DB)( c map[string] string, l int){
	if len(cache) < maxLenghtCache{
		cache, lenghtCache = AddItemToCache(cache, timeItem, countRequest, key, value, lenghtCache, client, db)
	}else {
		redis_storage.SetKeyValueToRedis(client, key, value)
		timeItem[key] = int(time.Now().Unix())
		countRequest[key] = 1
		fmt.Println("Cache full, cache to redis")
		return cache, lenghtCache
	}
	return cache, lenghtCache
}

// get item in cache with key
func GetItemInCache(cache map[string]string, timeItem map[string] int, countRequest map[string]int, key string, client * redis.Client)string{
	if val, ok := cache[key]; ok {
		GetItemCache(cache, timeItem, countRequest, key, client)
		return val
	}else { // else, get in redis
		timeItem[key] = int(time.Now().Unix())
		countRequest[key] = countRequest[key] + 1
		val,_ := redis_storage.GetValueFromKeyRedis(client, key)
		return val
	}
}

func RemoveItemInCache(cache map[string] string, timeItem map[string] int, countRequest map[string] int, key string, lenghtCache int, client * redis.Client){
	// if key exist in mem cache

	if cache[key] != "" {
		RemoveItemCache(cache, timeItem, countRequest, key, lenghtCache)
	}else {
		er := redis_storage.DelKeyValueRedis(client, key)
		fmt.Println(er)
		delete(timeItem, key)
		delete(countRequest, key)
	}
}

func WorkerCache(cache map[string] string, timeItem map[string] int, countRequest map[string] int, lenghtCache int, client * redis.Client){
	timeNow := time.Now().Unix()

	// check time out in cache
	for key := range timeItem{
		if(int(timeNow) - timeItem[key] >= timeOut){
			fmt.Println(key, " Time out")
			delete(timeItem, key)
			RemoveItemInCache(cache, timeItem, countRequest, key, lenghtCache, client)
		}else {
			if countRequest[key] > 3{
				fmt.Println(key, " Have count : ", countRequest[key])
			}
		}
	}

}
