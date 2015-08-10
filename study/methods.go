package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// second example start

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

// second example end

func main() {
	v := &Vertex{3, 4}
	fmt.Println(v.Abs())

	// second example start
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
	// second example end
}
