package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	var z float64 = 1
	for i := 1; i <= int(x); i++ {
		z = z - (z*z-float64(i))/(2*z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(5))
}
