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
	s1 := "foo"
	s2 := "foobar"

	fmt.Printf("s1 size: %T, %d\n", s1, unsafe.Sizeof(s1))
	fmt.Printf("s2 size: %T, %d\n", s2, unsafe.Sizeof(s2))
	fmt.Printf("s1 len: %T, %d\n", s1, len(s1))
	fmt.Printf("s2 len: %T, %d\n", s2, len(s2))
	//timeNow := time.Now().UnixNano()
	//time.Sleep(1 * time.Second)
	//timeAfter := (time.Now().UnixNano() )
	//fmt.Println("Time")
	//fmt.Println(float32(timeAfter - ((timeNow)))/ 1000000000)
	//for k := range ma{
	//	fmt.Println(k)
	//}
	//fmt.Println("_____")
	//fmt.Println(float64(1)/3)
	//if float32(4)/3 > 1{
	//	fmt.Println("DOne")
	//}


}
