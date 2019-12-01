package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func roots(d float64, a float64, b float64) (float64, float64) {
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

func delta(a float64, b float64, c float64) float64 {
	return math.Pow(float64(b), 2) - 4*a*c
}

func readFromArgs(a []string, p int) float64 {
	f, err := strconv.ParseFloat(a[p], 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func main() {
	clArguments := os.Args[1:] // We got command line arguments
	coefA := readFromArgs(clArguments, 0)
	coefB := readFromArgs(clArguments, 1)
	coefC := readFromArgs(clArguments, 2)
	vDelta := delta(coefA, coefB, coefC)
	fmt.Println(vDelta)
	fmt.Println(roots(vDelta, coefA, coefB))
}
