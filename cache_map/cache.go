package cache_map

import (
	"github.com/go-redis/redis"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
	"fmt"
	"github.com/tungct/go-keyvaluedb/redis_storage"
	"github.com/tungct/go-keyvaluedb/leveldb_storage"
	"log"
	"sync"
)

const timeOut  = 10 // time second
const nano2second = 1000000000 // to convert nano second to second
const frequencyThreshold = 3
const frequencyThresholdDb = 5
const maxLenRedis = 30000
var mutex = &sync.Mutex{}
// cache a item
func PutItemToCache(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, timeItemDb *map[string] int, countRequest *map[string]int, key string, value string, lenghtCache *int, client *redis.Client, db *leveldb.DB){
	// check item exits in cache, redis or levelDb
	mutex.Lock()
	if CheckItemInCache(cache, key) == true { // || redis_storage.CheckItemInRedis(client, key) ==true|| leveldb_storage.CheckItemInDb(db, key) == true
		mutex.Unlock()
		return
	}else {
		//else, put item to cache or redis or levelDb
		if len(*cache) < (maxLenghtCache) {
			AddItemToCache(cache, timeItemInit, timeItemCache, countRequest, key, value, lenghtCache, client, db)
		} else {
			if len(*timeItemRedis) < maxLenRedis {
				redis_storage.SetKeyValueToRedis(client, key, value)
				timeNow := int(time.Now().UnixNano())
				(*timeItemRedis)[key] = timeNow
				(*timeItemInit)[key] = timeNow
				(*countRequest)[key] = 1
				fmt.Println("Cache full, cache to redis")
			} else {
				er := leveldb_storage.SetKeyValueToLevelDb(db, key, value)
				if er == nil {
					timeNow := int(time.Now().UnixNano())
					(*timeItemDb)[key] = timeNow
					(*timeItemInit)[key] = timeNow
					(*countRequest)[key] = 1
					fmt.Println("Push to levelDb")
				} else {
					log.Println(er)
				}
			}
		}
		mutex.Unlock()
		return
	}
}

// print all item in memcache
func PrintCache(cache *map[string] string){
	mutex.Lock()
	for k := range *cache{
		fmt.Println(k, " in cache ")
	}
	mutex.Unlock()
}

// get item in cache with key
func GetItemInCache(cache *map[string]string, timeItemCache *map[string] int, timeItemRedis *map[string] int, timeItemDb *map[string] int, countRequest *map[string]int, key string, client * redis.Client, db * leveldb.DB)string{
	mutex.Lock()
	if val, ok := (*cache)[key]; ok {
		GetItemCache(cache, timeItemCache, countRequest, key, client)
		mutex.Unlock()
		return val
	}else if redis_storage.CheckItemInRedis(client, key) == true{ // else, get in redis
		(*timeItemRedis)[key] = int(time.Now().UnixNano())
		(*countRequest)[key] = (*countRequest)[key] + 1
		val,_ := redis_storage.GetValueFromKeyRedis(client, key)
		mutex.Unlock()
		return val
	}else if leveldb_storage.CheckItemInDb(db, key) == true{
		(*timeItemDb)[key] = int(time.Now().UnixNano())
		(*countRequest)[key] = (*countRequest)[key] + 1
		val,_ := leveldb_storage.GetValueFromKeyLevelDb(db, key)
		mutex.Unlock()
		return val
	}
	mutex.Unlock()
	return ""
}

// worker to check and delete timeout item
func TimeOutWorker(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, timeItemDb *map[string] int, countRequest *map[string] int, lenghtCache *int, client * redis.Client, db *leveldb.DB){
	timeNow := time.Now().UnixNano()
	mutex.Lock()

	// check time out in cache
	for key := range *timeItemCache{
		if(float32(int(timeNow) - (*timeItemCache)[key]) / nano2second >= timeOut){
			fmt.Println(key, " Time out Cache")
			if CheckItemInCache(cache, key) == true {
				RemoveItemCache(cache, timeItemInit, timeItemCache, countRequest, key, lenghtCache)
			}
			PrintCache(cache)
			fmt.Println("\n")
		}
	}
	for key := range *timeItemRedis{
		if(float32(int(timeNow) - (*timeItemRedis)[key]) / nano2second >= timeOut){
			fmt.Println(key, " Time out Redis")
			redis_storage.DelKeyValueRedis(client, key)
			delete(*timeItemRedis, key)
			delete(*timeItemInit, key)
			delete(*countRequest, key)
			PrintCache(cache)
			fmt.Println("\n")
		}
	}
	for key := range *timeItemDb{
		if(float32(int(timeNow) - (*timeItemDb)[key]) / nano2second >= timeOut){
			fmt.Println(key, "Time out Db")
			leveldb_storage.DelKeyValueLevelDb(db, key)
			delete(*timeItemDb, key)
			delete(*timeItemInit, key)
			delete(*countRequest, key)
			PrintCache(cache)
			fmt.Println("\n")
		}
	}
	mutex.Unlock()

}

// check item in redis have high-frequency, put it to mem cache
func FrequencyWorker(cache *map[string] string, timeItemInit *map[string] int, timeItemCache *map[string] int, timeItemRedis *map[string] int, timeItemDb *map[string]int, countRequest *map[string] int, lenghtCache *int, client * redis.Client, db *leveldb.DB) {

	timeNow := time.Now().UnixNano()

	// only check item in redis
	mutex.Lock()
	if len(*timeItemRedis) > 0 {
		for key := range *timeItemRedis {
			if ((float32((int(timeNow) - (*timeItemInit)[key])) / nano2second) >= 2) && (((float32((*countRequest)[key])) / ((float32(int(timeNow) - (*timeItemRedis)[key])) / nano2second)) > frequencyThreshold ) {
				if redis_storage.CheckItemInRedis(client, key) == true {
					fmt.Println("in redis")
					fmt.Println("Lenght cache : ", *lenghtCache)
					if *lenghtCache < maxLenghtCache {
						value, _ := redis_storage.GetValueFromKeyRedis(client, key)
						PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, timeItemDb, countRequest, key, value, lenghtCache, client, db)
						redis_storage.DelKeyValueRedis(client, key)
						delete((*timeItemRedis), key)
						fmt.Println("Put item from redis to cache ")
						PrintCache(cache)
					} else {
						// check item have lowest frequency in cache, delete it
						var kMin string
						check := false
						for k := range *cache {
							if (float32((int(timeNow) - (*timeItemCache)[k])) / nano2second) >= 2 {
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
						PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, timeItemDb, countRequest, key, value, lenghtCache, client, db)
						redis_storage.DelKeyValueRedis(client, key)
						delete(*timeItemRedis, key)
						fmt.Println("push from redis, Check in cache : ", (*cache)[key])
						PrintCache(cache)
						fmt.Println("\n")

					}
				}
			}
		}
	}

	// check item in levelDb
	if len(*timeItemDb) > 0 {
		for key := range *timeItemDb {
			if ((float32((int(timeNow) - (*timeItemInit)[key])) / nano2second) >= 2) && (((float32((*countRequest)[key])) / ((float32(int(timeNow) - (*timeItemDb)[key])) / nano2second)) > frequencyThreshold ) {
				fmt.Println("in Db")
				fmt.Println("Lenght cache : ", *lenghtCache)
				if *lenghtCache < maxLenghtCache {
					value, _ := leveldb_storage.GetValueFromKeyLevelDb(db, key)

					PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, timeItemDb, countRequest, key, value, lenghtCache, client, db)
					leveldb_storage.DelKeyValueLevelDb(db, key)
					delete((*timeItemDb), key)
					PrintCache(cache)
				} else {
					if ((float32(int(timeNow) - (*timeItemDb)[key])) / nano2second) > frequencyThresholdDb {

						// check item have lowest frequency in cache, delete it
						var kMin string
						check := false
						for k := range *cache {
							if (float32((int(timeNow) - (*timeItemCache)[k])) / nano2second) >= 2 {
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

						value, _ := leveldb_storage.GetValueFromKeyLevelDb(db, key)

						// put item in redis to cache
						PutItemToCache(cache, timeItemInit, timeItemCache, timeItemRedis, timeItemDb, countRequest, key, value, lenghtCache, client, db)
						leveldb_storage.DelKeyValueLevelDb(db, key)
						delete(*timeItemDb, key)
						fmt.Println("push from redis, Check in cache : ", (*cache)[key])
						for k := range *cache {
							fmt.Println("Key in cache : ", k)
						}
					} else {
						// check item have lowest frequency in cache, delete it
						if client.DbSize().Val() >= 4 {
							var kMin string
							check := false
							for k := range *timeItemRedis {
								if (*timeItemRedis)[k] > 2 && (float32((int(timeNow) - (*timeItemRedis)[k])) / nano2second) >= 2 {
									if check == false {
										kMin = k
										check = true
									} else {
										// check item in cache, if item have time >=2 and have lowest-frequency
										if (float32((*countRequest)[k]) / float32(int(timeNow) - (*timeItemRedis)[k]) < float32((*countRequest)[kMin]) / float32(int(timeNow) - (*timeItemRedis)[kMin])) {
											kMin = k
										}
									}
								}
							}
							fmt.Println("Del item : ", kMin, " in cache")
							// delete item have lowest-frequency in cache
							redis_storage.DelKeyValueRedis(client, key)
							delete(*timeItemRedis, kMin)
							delete(*timeItemInit, kMin)
							delete(*countRequest, kMin)
						}

						value, _ := leveldb_storage.GetValueFromKeyLevelDb(db, key)

						// put item in redis to cache
						redis_storage.SetKeyValueToRedis(client, key, value)
						timeNow := int(time.Now().UnixNano())
						(*timeItemRedis)[key] = timeNow
						(*timeItemInit)[key] = timeNow
						(*countRequest)[key] = 1

						leveldb_storage.DelKeyValueLevelDb(db, key)
						delete(*timeItemDb, key)
					}

				}

			}

		}
	}
	mutex.Unlock()
}

