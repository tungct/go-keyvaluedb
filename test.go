package main

import (
	"unsafe"
	"fmt"
	"reflect"
	"time"
)

func main (){
	var m map[int] int
	m = make(map[int] int)
	ma := make(map[string]string)
	m[1] = 1
	ma["1"] = "1"
	ma["2"] = "2"
	fmt.Println(reflect.TypeOf(len(ma)))
	fmt.Println(ma["2"])
	if val, ok := ma["3"]; ok {
		fmt.Println(ok)
		fmt.Println(val)
	}
	fmt.Println((int(unsafe.Sizeof(ma["1"]))))
	fmt.Println(reflect.TypeOf(time.Now().Unix()))
	timeNow := time.Now().Unix()
	time.Sleep(1 * time.Second)
	timeAfter := time.Now().Unix()
	fmt.Println(timeAfter - timeNow)
	for k := range ma{
		fmt.Println(k)
	}
}
