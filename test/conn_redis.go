package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"reflect"
	"strconv"

	"time"
	"sync"
)
var mutex = &sync.Mutex{}
func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	err = client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(reflect.TypeOf(client.DbSize().Val()))
	count := 0

	// send n message to server grpc
	for i:=1;i<=10;i++ {
		//id := rand.Intn(100)
		go func() {
			for {
				count = count + 1
				fmt.Println(count)
				content := strconv.Itoa(int(count))
				//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				//defer cancel()
				client.Set(content, content, 0).Err()
			}
		}()
	}
	time.Sleep(5*time.Second)


	//val, err := client.Get("key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//
	//val2, err := client.Get("key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
}


func main() {
	ExampleNewClient()
}
