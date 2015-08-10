package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
  return [1,2,3,4][4,3,2,1]
}

func main() {
	pic.Show(Pic)
}
