package leveldb_storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// init levelDb client
func InitConnLevelDb(path string, o *opt.Options) (db *leveldb.DB, err error){
	db, er := leveldb.OpenFile(path, o)
	return db, er
}

// store key-value type string to levelDb
func SetKeyValueToLevelDb(db *leveldb.DB, key string, value string) error{
	err := db.Put([]byte(key), []byte(value), nil)
	return err
}

// get value with key from levelDb
func GetValueFromKeyLevelDb(db *leveldb.DB, key string) (val string, err error){
	var value string
	data, er := db.Get([]byte(key), nil)
	if er == nil{
		value = string(data[:len(data)])
		return value, er
	}
	return value, er
}

// del key-value in levelDb
func DelKeyValueLevelDb(db *leveldb.DB, key string) error{
	del := db.Delete([]byte(key), nil)
	return del
}

// check item exits in levelDb
func CheckItemInDb(db *leveldb.DB, key string) bool{
	_, er := db.Get([]byte(key), nil)
	if er != nil{
		return false
	}else {
		return true
	}
}
