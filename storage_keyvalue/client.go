package main

import (
	"log"
	"os"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/tungct/go-keyvaluedb/grpc"
)

const (
	address     = "localhost:8888"
	defaultContent = "Hello"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSendMessageClient(conn)

	for i:=0;i<10;i++ {

		content := defaultContent
		if len(os.Args) > 1 {
			content = os.Args[1]
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SendMessage(ctx, &pb.Message{Id:1, Content:content})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r)
		time.Sleep(1*time.Second)
	}
}