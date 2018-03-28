package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/tungct/go-keyvaluedb/grpc"
	"github.com/go-redis/redis"
	"github.com/tungct/go-keyvaluedb/redis_storage"
	"google.golang.org/grpc/reflection"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tungct/go-keyvaluedb/leveldb_storage"
	"github.com/tungct/go-keyvaluedb/cache_map"
	"time"
	"github.com/tungct/go-keyvaluedb/utils"
)

var cache map[string] string
var timeItemInit map[string] int // map key and time init per item
var timeItemRedis map[string] int // map key and time redis
var timeItemDb map[string] int // map key and time levelDb
var timeItemCache map[string] int // map key and time in cache
var countRequest map[string] int // map key and count Request item
var count int
var lenghtCache int // lenght of mem cache

// redis client
var connRedis *redis.Client

// levelDb client
var connLevelDb * leveldb.DB

type server struct{}

// implement SendMessage method grpc
func (s *server) SendMessage(ctx context.Context, in *pb.Message)(*pb.MessageResponse, error){
	var err error
	fmt.Println("Rec message from client : ", in)
	// cache item
	if in.Key != "-1"{

		count = count + 1
		cache_map.PutItemToCache(&cache, &timeItemInit, &timeItemCache, &timeItemRedis, &timeItemDb, &countRequest, in.Key, in.Value, &lenghtCache, connRedis, connLevelDb)

		fmt.Println("Lenght Cache : ", lenghtCache)
		fmt.Println(count)
	}else { // get item in cache
		value := cache_map.GetItemInCache(&cache, &timeItemCache, &timeItemRedis, &timeItemDb, &countRequest, in.Value, connRedis, connLevelDb)
		fmt.Println("Value", value)
		return &pb.MessageResponse{Content: "Response from server" + value }, nil
	}

	if err != nil{
		fmt.Println(err)
	}

	// return response to client
	return &pb.MessageResponse{Content: "Response from server" }, nil
}

func main(){

	config := utils.LoadConfigServer("../config/server.conf")

	count = 0
	// init redis client
	connRedis = redis_storage.InitConnRedis()

	// init levelDb client
	connLevelDb , _ = leveldb_storage.InitConnLevelDb(config.PATHLEVELDB, nil)

	cache = cache_map.InitCacheMap()
	timeItemInit = make(map[string] int)
	timeItemCache = make(map[string]int)
	timeItemRedis = make(map[string] int)
	timeItemDb = make(map[string] int)
	countRequest = make(map[string]int)
	lenghtCache = 0

	go func() {
		for{
			cache_map.Worker(&cache, &timeItemInit, &timeItemCache, &timeItemRedis, &timeItemDb, &countRequest, &lenghtCache, connRedis, connLevelDb)
			time.Sleep(1 * time.Second)
		}
	}()


	// init grpc server
	lis, er := net.Listen("tcp", config.PORT)
	fmt.Println("Server listen at port 8888")
	if er != nil{
		log.Fatalf("failed to serve : %v", er)
	}
	s := grpc.NewServer()
	pb.RegisterCServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}