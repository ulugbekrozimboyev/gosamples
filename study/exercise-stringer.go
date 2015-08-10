package main

import "fmt"

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.
func (obj IPAddr) String() string {

	//src := obj[:]
	//fmt.Println(src)
	var str string = ""

	for i, value := range obj {
		str += fmt.Sprintf("%v", value)
		if i != len(obj)-1 {
			str += "."
		}
	}
	//fmt.Println(str)
	return str
}

func main() {
	addrs := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for n, a := range addrs {
		fmt.Printf("%v: %v\n", n, a)
	}
}
