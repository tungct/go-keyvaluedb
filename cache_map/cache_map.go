package cache_map

import (
	"unsafe"
)

const sizeCache = 10 * 1024 // 10 kb

func InitCacheMap() (c map[string] string, len int){
	cache := make(map[string] string)
	lenCache := sizeCache / ( int(unsafe.Sizeof("")) + int(unsafe.Sizeof(cache[""])) )
	return cache, lenCache
}

func RemoveItemCache(cache map[string] string, key string, lenghtCache int) (c map[string] string, len int){
	delete(cache, key)
	lenghtCache --
	return cache, lenghtCache
}

func AddItemToCache(cache map[string] string, key string, value string, lenghtCache int) (c map[string] string, len int){
	if cache[key] != ""{
		return cache, lenghtCache
	}
	cache[key] = value
	lenghtCache ++
	return cache, lenghtCache
}