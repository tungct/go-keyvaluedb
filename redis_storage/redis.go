package redis_storage

import (
	"github.com/go-redis/redis"
	"fmt"
)

// init redis client
func InitConnRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

// store key-value type string to redis
func SetKeyValueToRedis(client *redis.Client, key string, value string) *redis.StatusCmd{
	err := client.Set(key, value, 0)
	return err
}

// get value with key in redis
func GetValueFromKeyRedis(client *redis.Client, key string) (value string, er error){
	fmt.Println("Value get from redis")
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, err
}

// delete value with key in redis
func DelKeyValueRedis(client * redis.Client, key string) *redis.IntCmd{
	err := client.Del(key)
	return err
}

func CheckItemInRedis(client * redis.Client, key string) bool{
	_, err := client.Get(key).Result()
	if err != redis.Nil {
		return true
	}else {
		return false
	}
}
