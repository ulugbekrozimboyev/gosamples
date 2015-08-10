package main

import (
	"fmt"
	"time"
)

func Update_stamp(sn, sta string) bool {
	const shrtFrm = "2006-01-02 15:04:05"

	t2, e := time.Parse(shrtFrm, sta)
	fmt.Println("err", e)
	fmt.Println("t2=", t2)
	nt := t2.Local()
	fmt.Println("nt=", nt)
	nx := t2.Unix()
	fmt.Println("NX: ", nx)

	return true
}

func main() {
	tt := "2015-05-21 11:17:27"
	sn := "123"
	Update_stamp(sn, tt)
}
