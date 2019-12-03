package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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

func getAnumber(msg string) float64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", msg)
	nfs, _ := reader.ReadString('\n')
	nf, errf := strconv.ParseFloat(strings.TrimSpace(nfs), 64)
	if errf != nil {
		log.Fatal("Invalid float number")
	}
	return nf
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	coefA := getAnumber("Enter coefficient a")
	coefB := getAnumber("Enter coefficient b")
	coefC := getAnumber("Enter coefficient c")
	deltaV := delta(coefA, coefB, coefC)
	x1, x2 := roots(deltaV, coefA, coefB)
	fmt.Printf("X1: %f, X2: %f\n", x1, x2)
	stringArr := []byte(fmt.Sprintf("A: %f, B: %f, C: %f, x1: %f, x2: %f", coefA, coefB, coefC, x1, x2))
	// Permission: -rw-r--r--
	err := ioutil.WriteFile("/tmp/arq2.txt", stringArr, 0644)
	check(err)
}
