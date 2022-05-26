package main

import(
	"fmt"
	"unsafe"
	"math/big"
)

func main(){
	n := big.NewInt(123)
	fmt.Println(unsafe.Sizeof(n))
}
