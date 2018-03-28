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
	"github.com/tungct/go-keyvaluedb/utils"
)

const (
	defaultContent = "Hello" // if arg is none
)

func main() {
	config := utils.LoadConfigServer("../../config/server.conf")
	// init grpc client
	conn, err := grpc.Dial(config.IP + config.PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCClient(conn)

	content := defaultContent
	if len(os.Args) > 1 {
		content = os.Args[1] // get content message from arg
	}

	// send n message to server grpc
	for i:=0;i<10;i++ {
		//id := rand.Intn(100)
		content = strconv.Itoa(int(7))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		message := &pb.Message{Key: content, Value:content}
		r, err := c.SendMessage(ctx, message)
		fmt.Println("Send message ", *message)
		if err != nil {
			log.Fatalf("fail to send: %v", err)
		}
		log.Printf("Rec from server : %s", r)
		//time.Sleep(1*time.Second)
	}
	time.Sleep(1 * time.Second)
	for i:=0;i<10;i++ {
		//id := rand.Intn(100)
		content = strconv.Itoa(int(8))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		message := &pb.Message{Key: content, Value:content}
		r, err := c.SendMessage(ctx, message)
		fmt.Println("Send message ", *message)
		if err != nil {
			log.Fatalf("fail to send: %v", err)
		}
		log.Printf("Rec from server : %s", r)
		//time.Sleep(1*time.Second)
	}
}