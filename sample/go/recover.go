package main

import (
	"fmt"
)

func test(a int, b int) int {
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println("in test")
	//		fmt.Println(err)
	//	}
	//}()
	return a / b
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in main")
			fmt.Println(err)
		}
	}()

	re := test(2, 0)

	fmt.Println(re)
}