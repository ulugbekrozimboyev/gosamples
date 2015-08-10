package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	var f []int
	f = append(f, 1)

	return func() int {
		len := len(f)
		//fmt.Println(len)
		if len > 2 {
			f = append(f, f[len-1]+f[len-2])
		} else {
			f = append(f, f[len-1])
		}
		//fmt.Println(f[len])
		return f[len]
	}

}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())

	}
}
