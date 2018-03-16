package main

import (
	"log"
	"os"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/tungct/go-keyvaluedb/grpc"
	//"math/rand"
	"fmt"
	"strconv"
)

const (
	address     = "localhost:8888"
	defaultContent = "Hello" // if arg is none
)

func main() {
	// init grpc client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSendMessageClient(conn)

	content := defaultContent
	if len(os.Args) > 1 {
		content = os.Args[1] // get content message from arg
	}

	// send n message to server grpc
	for i:=0;i<10;i++ {
		//id := rand.Intn(100)
		content = strconv.Itoa(int(i))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		message := &pb.Message{Id:int32(-1), Content:content}
		r, err := c.SendMessage(ctx, message)
		fmt.Println("Send message ", *message)
		if err != nil {
			log.Fatalf("fail to send: %v", err)
		}
		log.Printf("Rec from server : %s", r)
		time.Sleep(1*time.Second)
	}
}