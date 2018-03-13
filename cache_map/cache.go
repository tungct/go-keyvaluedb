package cache_map

import (

)

func PutItemToCache(cache map[string] string, key string, value string, lenghtCache int)( c map[string] string, len int){
	cache, lenghtCache = AddItemToCache(cache, key, value, lenghtCache)
	return cache, lenghtCache
}
