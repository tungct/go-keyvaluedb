package storage_redis

import (
	"github.com/go-redis/redis"
	"fmt"
)

func InitConnRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

func SetKeyValue(client *redis.Client, key string, value string) *redis.StatusCmd{
	err := client.Set(key, value,0)
	return err
}

func GetValueFromKey(client *redis.Client, key string) (value string, er error){
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, err
}
