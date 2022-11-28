package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func Between[T constraints.Ordered](a, b, x T) bool {
	if x <= b && x >= a {
		return true
	}
	return false
}

func main() {
	if Between[int](2, 10, 5) {
		fmt.Println("Yes")
	}
	beetween := Between[int]
	fmt.Println(beetween(2, 10, 8))
}
