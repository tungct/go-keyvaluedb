package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/tungct/go-keyvaluedb/grpc"
	"github.com/go-redis/redis"
	"github.com/tungct/go-keyvaluedb/storage_redis"
	"google.golang.org/grpc/reflection"
	"github.com/tungct/go-keyvaluedb/messqueue"
	"fmt"
)

const (
	port = ":8888"
)

var messQueue chan pb.Message
var connRedis *redis.Client

type server struct{}

func (s *server) SendMessage(ctx context.Context, in *pb.Message)(*pb.MessageResponse, error){
	var err error
	fmt.Println("Rec message from client : ", in)
	messQueue, err = messqueue.PutMessage(messQueue, *in, connRedis)
	if err != nil{
		fmt.Println(err)
	}
	return &pb.MessageResponse{Content: "Response from server"}, nil
}

func main(){
	connRedis = storage_redis.InitConnRedis()
	messQueue = messqueue.InitMessageQueue()
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