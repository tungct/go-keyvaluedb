package messqueue

import (
	pb "github.com/tungct/go-keyvaluedb/grpc"
	"github.com/go-redis/redis"
	"github.com/tungct/go-keyvaluedb/storage_redis"
	"errors"
)

const lengthQueue = 10

func InitMessageQueue() (chan pb.Message){
	messQueue := make(chan pb.Message, lengthQueue)
	return messQueue
}

func PutMessage(queue chan pb.Message, mess pb.Message, client *redis.Client )( messQueue chan pb.Message, er error){
	if len(queue) < lengthQueue{
		queue <- mess
		return queue, nil
	}else {
		storage_redis.SetKeyValue(client, "test", mess.Content)
		return queue, errors.New("Push to redis")
	}
}

func GetMessageInQueue(queue chan pb.Message)(messQueue chan pb.Message, mess pb.Message, er error){
	var message pb.Message
	if len(queue) > 0{
		message = <- queue
		return queue, message, nil
	}else {
		return queue, message, errors.New("Not exits message in queue")
	}
}