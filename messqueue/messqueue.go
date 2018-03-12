package messqueue

import (
	pb "github.com/tungct/go-keyvaluedb/grpc"
	"github.com/go-redis/redis"
	"github.com/tungct/go-keyvaluedb/redis_storage"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tungct/go-keyvaluedb/leveldb_storage"
	"strconv"
)

// max length queue
const lengthQueue = 5

// init message queue
func InitMessageQueue() (chan pb.Message){
	messQueue := make(chan pb.Message, lengthQueue)
	return messQueue
}

// push message to message Queue if lenght queue < max Lenght Queue, else push message to redis and levelDb
func PutMessage(queue chan pb.Message, mess pb.Message, client *redis.Client, db *leveldb.DB)( messQueue chan pb.Message, er error){
	if len(queue) < lengthQueue{
		// push message to queue
		queue <- mess
		return queue, nil
	}else {
		// convert key type int32 to string, push
		var key string = strconv.Itoa(int(mess.Id))
		redis_storage.SetKeyValueToRedis(client, key, mess.Content)
		leveldb_storage.SetKeyValueToLevelDb(db, key, mess.Content)
		return queue, errors.New("Full Queue, Push to redis, levelDb")
	}
}

// get message in queue
func GetMessageInQueue(queue chan pb.Message)(messQueue chan pb.Message, mess pb.Message, er error){
	var message pb.Message
	if len(queue) > 0{
		message = <- queue
		return queue, message, nil
	}else {
		return queue, message, errors.New("Not exits message in queue")
	}
}