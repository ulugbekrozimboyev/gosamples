package main

import (
	"fmt"
	"runtime"
)

func getInterval(age int) string {
	switch age {
	case 1:
		fallthrough
	case 2, 3, 4, 5:
		fallthrough
	case 10:
		return "young"
	default:
		return "adult"
	}
}

func main() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.", os)
	}

	fmt.Println(getInterval(1))
	fmt.Println(getInterval(5))
	fmt.Println(getInterval(23))

	k := 5
	switch k {
	case 4:
		fmt.Println("was <= 4")
		fallthrough
	case 5:
		fmt.Println("was <= 5")
		fallthrough
	case 6:
		fmt.Println("was <= 6")
		fallthrough
	case 7:
		fmt.Println("was <= 7")
		fallthrough
	case 8:
		fmt.Println("was <= 8")
		fallthrough
	default:
		fmt.Println("default case")
	}

}
