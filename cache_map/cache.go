package cache_map

import (
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
	"fmt"
	"github.com/tungct/go-keyvaluedb/redis_storage"
)

const timeOut  = 10 // time second
const nano2second = 1000000000

// cache a item
func PutItemToCache(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, countRequest *map[string]int, key string, value string, lenghtCache *int, client *redis.Client, db *leveldb.DB){
	if len(*cache) < (maxLenghtCache){
		AddItemToCache(cache, timeItemInit, timeItemCache, countRequest, key, value, lenghtCache, client, db)
	}else {
		redis_storage.SetKeyValueToRedis(client, key, value)
		timeNow := int(time.Now().UnixNano())
		(*timeItemRedis)[key] = timeNow
		(*timeItemInit)[key] = timeNow
		(*countRequest)[key] = 1
		fmt.Println("Cache full, cache to redis")
	}
	return
}

func PrintCache(cache *map[string] string){
	for k := range *cache{
		fmt.Println(k)
	}
}

// get item in cache with key
func GetItemInCache(cache *map[string]string, timeItemCache *map[string] int, timeItemRedis *map[string] int, countRequest *map[string]int, key string, client * redis.Client)string{
	if val, ok := (*cache)[key]; ok {
		GetItemCache(cache, timeItemCache, countRequest, key, client)
		return val
	}else { // else, get in redis
		(*timeItemRedis)[key] = int(time.Now().UnixNano())
		(*countRequest)[key] = (*countRequest)[key] + 1
		val,_ := redis_storage.GetValueFromKeyRedis(client, key)
		return val
	}
}


//func RemoveItemInCache(cache map[string] string, timeItem map[string] int, countRequest map[string] int, key string, lenghtCache int, client * redis.Client){
//	// if key exist in mem cache
//
//	if cache[key] != "" {
//		RemoveItemCache(cache, timeItem, countRequest, key, lenghtCache)
//	}else {
//		er := redis_storage.DelKeyValueRedis(client, key)
//		fmt.Println(er)
//		delete(timeItem, key)
//		delete(countRequest, key)
//	}
//}

func TimeOutWorker(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, countRequest *map[string] int, lenghtCache *int, client * redis.Client, db *leveldb.DB){
	timeNow := time.Now().UnixNano()

	// check time out in cache
	for key := range *timeItemCache{
		if(float32(int(timeNow) - (*timeItemCache)[key]) / nano2second >= timeOut){
			fmt.Println(key, " Time out")
			if CheckItemInCache(cache, key) == true {
				RemoveItemCache(cache, timeItemInit, timeItemCache, countRequest, key, lenghtCache)
			}
		}
	}
	for key := range *timeItemRedis{
		if(float32(int(timeNow) - (*timeItemRedis)[key]) / nano2second >= timeOut){
			fmt.Println(key, " Time out 2")
			redis_storage.DelKeyValueRedis(client, key)
			delete(*timeItemRedis, key)
			delete(*timeItemInit, key)
			delete(*countRequest, key)
		}
	}

}

func FrequencyWorker(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, countRequest *map[string] int, lenghtCache *int, client * redis.Client, db *leveldb.DB) {

	timeNow := time.Now().UnixNano()

	for key := range *timeItemRedis {
		if ((float32((int(timeNow) - (*timeItemInit)[key]) ) / nano2second) >= 2) && (((float32((*countRequest)[key])) / ((float32(int(timeNow) - (*timeItemRedis)[key])) / nano2second)) > 3 ){
			if redis_storage.CheckItemInRedis(client, key) == true {
				fmt.Println("in redis")
				fmt.Println("Lenght cache : ", *lenghtCache)
				if *lenghtCache < maxLenghtCache {
					value, _ := redis_storage.GetValueFromKeyRedis(client, key)
					PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, countRequest, key, value, lenghtCache, client, db)
					redis_storage.DelKeyValueRedis(client, key)
					PrintCache(cache)
				} else {
					// check item have lowest frequency in cache, delete it
					var kMin string
					check := false
					for k := range *cache {
						if (float32((int(timeNow) - (*timeItemCache)[k]) ) / nano2second) >= 2 {
							if check == false {
								kMin = k
								check = true
							} else {
								// check item in cache, if item have time >=2 and have lowest-frequency
								if (float32((*countRequest)[k]) / float32(int(timeNow) - (*timeItemCache)[k]) < float32((*countRequest)[kMin]) / float32(int(timeNow) - (*timeItemCache)[kMin])) {
									kMin = k
								}
							}
						}
					}
					fmt.Println("Del item : ", kMin, " in cache")
					// delete item have lowest-frequency in cache
					RemoveItemCache(cache, timeItemInit, timeItemCache, countRequest, kMin, lenghtCache)
					PrintCache(cache)
					value, _ := redis_storage.GetValueFromKeyRedis(client, key)

					// put item in redis to cache
					PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, countRequest, key, value, lenghtCache, client, db)
					redis_storage.DelKeyValueRedis(client, key)
					delete(*timeItemRedis, key)
					fmt.Println("push from redis, Check in cache : ", (*cache)[key])
					for k := range *cache {
						fmt.Println("Key in cache : ", k)
					}

				}
			}
		}
	}
}

