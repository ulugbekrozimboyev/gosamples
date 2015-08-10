package main

import (
	"fmt"
	"time"
)

func main() {
	var d int = 0
	var beginTime = time.Now()
	for i := 0; i < 100000; i++ {

		for c := 0; c < 100; c++ {
			d = 1
		}
	}
	var endTime = time.Now()
	fmt.Println(beginTime, endTime)
	fmt.Println(time.Since(beginTime))
	fmt.Println(d)
}
