package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

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
}
