package main

import (
	"fmt"
	"math"
	"time"
)

type ErrNegativeSqrt struct {
	When time.Time
	What string
}

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, &ErrNegativeSqrt{
			time.Now(),
			fmt.Sprintf("cannot Sqrt negative number: %v", x),
		}
	} else {
		return math.Sqrt(x), nil
	}

}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
