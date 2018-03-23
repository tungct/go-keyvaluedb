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
func PutItemToCache(cache *map[string] string, timeItem *map[string] int, timeItemCache *map[string] int, countRequest *map[string]int, key string, value string, lenghtCache *int, client *redis.Client, db *leveldb.DB){
	if len(*cache) < (maxLenghtCache){
		AddItemToCache(cache, timeItem, countRequest, key, value, lenghtCache, client, db)
	}else {
		redis_storage.SetKeyValueToRedis(client, key, value)
		(*timeItemCache)[key] = int(time.Now().Unix())
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
func GetItemInCache(cache *map[string]string, timeItem *map[string] int, timeItemCache *map[string] int, countRequest *map[string]int, key string, client * redis.Client)string{
	if val, ok := (*cache)[key]; ok {
		GetItemCache(cache, timeItem, countRequest, key, client)
		return val
	}else { // else, get in redis
		(*timeItemCache)[key] = int(time.Now().Unix())
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

func WorkerCache(cach *map[string] string, timeIte *map[string] int, countReques *map[string] int, lenghtCach *int, client * redis.Client, db *leveldb.DB){
	timeNow := time.Now().Unix()
	//countRequest := * countReques
	//lenghtCache := * lenghtCach

	// check time out in cache
	for key := range *timeIte{
		if(int(timeNow) - (*timeIte)[key] >= timeOut){
			fmt.Println(key, " Time out")
			if CheckItemInCache(cach, key) == true {
				RemoveItemCache(cach, timeIte, countReques, key, lenghtCach)
				fmt.Println("Lenght cache : ", *lenghtCach)

			}else { // delete in redis
				redis_storage.DelKeyValueRedis(client, key)
			}
		}
	}

}

func Worker(cache *map[string] string, timeItem *map[string] int, timeItemCache *map[string] int, countRequest *map[string] int, lenghtCache *int, client * redis.Client, db *leveldb.DB) {

	timeNow := time.Now().Unix()

	for key := range *timeItemCache {
		if (float32((*countRequest)[key]) / (float32(int(timeNow) - (*timeItemCache)[key]))) > 3 {
			if redis_storage.CheckItemInRedis(client, key) == true {
				fmt.Println("in redis")
				fmt.Println("Lenght cache : ", *lenghtCache)
				if *lenghtCache < maxLenghtCache {
					value, _ := redis_storage.GetValueFromKeyRedis(client, key)
					PutItemToCache(cache, timeItem, timeItemCache, countRequest, key, value, lenghtCache, client, db)
					redis_storage.DelKeyValueRedis(client, key)
					fmt.Println("Check in cache : ", (*cache)[key])
					PrintCache(cache)
				} else {

					// check item have lowest frequency in cache, delete it
					var kMin string
					check := false
					for k := range *cache {
						if (int(timeNow) - (*timeItem)[k] >= 2) {
							if check == false {
								kMin = k
								check = true
							} else {
								// check item in cache, if item have time >=2 and have lowest-frequency
								if (float32((*countRequest)[k]) / float32(int(timeNow) - (*timeItem)[k]) < float32((*countRequest)[kMin]) / float32(int(timeNow) - (*timeItem)[kMin])) {
									kMin = k
								}
							}
						}
					}
					fmt.Println("Del item : ", kMin, " in cache")
					// delete item have lowest-frequency in cache
					RemoveItemCache(cache, timeItem, countRequest, key, lenghtCache)
					PrintCache(cache)
					for k,v := range *countRequest{
						fmt.Println("K and V in countRequest : ", k,v)
					}
					value, _ := redis_storage.GetValueFromKeyRedis(client, key)

					// put item in redis to cache
					PutItemToCache(cache, timeItem, timeItemCache, countRequest, key, value, lenghtCache, client, db)
					redis_storage.DelKeyValueRedis(client, key)
					fmt.Println("push from redis, Check in cache : ", (*cache)[key])
					for k := range *cache {
						fmt.Println("Key in cache : ", k)
					}

				}
			}
		}
	}
}

