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
)

const (
	port = ":8888"
        pathLevelDb = "../leveldb_storage/keyvaluedb"
)

var messQueue chan pb.Message

var cache map[string] string
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
		cache, lenghtCache = cache_map.PutItemToCache(cache, strconv.Itoa(int(in.Id)), in.Content, lenghtCache, connRedis, connLevelDb)
		fmt.Println("Lenght Cache : ", lenghtCache)
	}else { // get item in cache
		value := cache_map.GetItemInCache(cache, in.Content, connRedis)
		fmt.Println("Value", value)
		return &pb.MessageResponse{Content: "Response from server" + value }, nil
	}

	fmt.Println(len(cache))
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
	lenghtCache = 0

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