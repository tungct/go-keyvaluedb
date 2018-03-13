package main

import (
	"unsafe"
	"fmt"
)

func main (){
	var m map[int] int
	m = make(map[int] int)
	ma := make(map[string]string)
	m[1] = 1
	ma["1"] = "1"
	fmt.Println(len(m))
	fmt.Println(ma["2"])
	fmt.Println((int(unsafe.Sizeof(ma["1"]))))
}
