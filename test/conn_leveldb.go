package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
)

func main()  {
	db, _ := leveldb.OpenFile("../leveldb_storage/keyvaluedb", nil)
	//_ = db.Put([]byte("key"), []byte("value"), nil)
	data, _ := db.Get([]byte("test"), nil)
	var value = string(data[:len(data)])
	fmt.Println(value)
	defer db.Close()
}