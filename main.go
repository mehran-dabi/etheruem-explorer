package main

import (
	"fmt"
	"math"
)

func main() {
	var i int
	var count int
	for {
		fib := Fib(i)
		sqrt := math.Sqrt(float64(fib))
		if math.Ceil(sqrt) == sqrt {
			fmt.Println(fib)
			count++
		}
		if count == 10 {
			break
		}
		i++
	}
}

func Fib(i int) int {
	if i == 0 || i == 1 {
		return 1
	}

	return Fib(i-1) + Fib(i-2)
}
