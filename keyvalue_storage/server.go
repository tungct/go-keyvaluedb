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
	"github.com/tungct/go-keyvaluedb/messqueue"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tungct/go-keyvaluedb/leveldb_storage"
	"github.com/tungct/go-keyvaluedb/cache_map"
	"strconv"
	"time"
)

const (
	port = ":8888"
        pathLevelDb = "../leveldb_storage/keyvaluedb"
)

var messQueue chan pb.Message

var cache map[string] string
var timeItemInit map[string] int
var timeItemRedis map[string] int // map key and time cache
var timeItemCache map[string] int
var countRequest map[string] int // map key and count Request item

var lenghtCache int

// redis client
var connRedis *redis.Client

// levelDb client
var connLevelDb * leveldb.DB

type server struct{}

// implement SendMessage method grpc
func (s *server) SendMessage(ctx context.Context, in *pb.Message)(*pb.MessageResponse, error){
	var err error
	fmt.Println("Rec message from client : ", in)

	// handle message receive from client with messageQueue, redis and levelDb
	//messQueue, err = messqueue.PutMessage(messQueue, *in, connRedis, connLevelDb)

	// cache item
	if int(in.Id) != -1{
		cache_map.PutItemToCache(&cache, &timeItemInit, &timeItemCache, &timeItemRedis, &countRequest, strconv.Itoa(int(in.Id)), in.Content, &lenghtCache, connRedis, connLevelDb)
		fmt.Println("Lenght Cache : ", lenghtCache)
	}else { // get item in cache
		value := cache_map.GetItemInCache(&cache, &timeItemCache, &timeItemRedis, &countRequest, in.Content, connRedis)
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
	// init redis client
	connRedis = redis_storage.InitConnRedis()

	// init levelDb client
	connLevelDb , _ = leveldb_storage.InitConnLevelDb(pathLevelDb, nil)

	// init message Queue
	messQueue = messqueue.InitMessageQueue()

	cache = cache_map.InitCacheMap()
	timeItemInit = make(map[string] int)
	timeItemCache = make(map[string]int)
	timeItemRedis = make(map[string] int)
	countRequest = make(map[string]int)
	lenghtCache = 0

	// worker to check cache per time second
	go func() {
		for{
			cache_map.TimeOutWorker(&cache,&timeItemInit, &timeItemCache, &timeItemRedis, &countRequest, &lenghtCache, connRedis, connLevelDb)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for{
			cache_map.FrequencyWorker(&cache, &timeItemInit, &timeItemCache, &timeItemRedis, &countRequest, &lenghtCache, connRedis, connLevelDb)
			time.Sleep(1 * time.Second)
		}
	}()


	// init grpc server
	lis, er := net.Listen("tcp", port)
	fmt.Println("Server listen at port 8888")
	if er != nil{
		log.Fatalf("failed to serve : %v", er)
	}
	s := grpc.NewServer()
	pb.RegisterSendMessageServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}