package main

import "fmt"

func s() {
	defer fmt.Println("11111111")
	fmt.Println("2222222")
}

func main() {
	defer fmt.Println("world")
	defer s()

	fmt.Println("hello")
	s2()
}

func s2() {
	defer fmt.Println("s2--1")
	fmt.Println("s2--2")
}
