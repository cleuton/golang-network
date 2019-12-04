package bhaskarautils

import (
	"math"
	"fmt"
)

func Roots(d float64, a float64, b float64) (float64, float64) {
	var (
		x1 float64
		x2 float64
	)
	if d > 0 {
		x1 = (-b + math.Sqrt(d)) / (2 * a)
		x2 = (-b - math.Sqrt(d)) / (2 * a)
	} else {
		if d == 0 {
			x1 = -b / (2 * a)
			x2 = x1
		} else {
			x1 = math.Inf(1)
			x2 = x1
		}
	}
	return x1, x2
}

func Delta(a float64, b float64, c float64) float64 {
	return math.Pow(float64(b), 2) - 4*a*c
}
