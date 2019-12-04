package main


import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"bhaskarautils"
)

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
	deltaV := bhaskarautils.Delta(coefA, coefB, coefC)
	x1, x2 := bhaskarautils.Roots(deltaV, coefA, coefB)
	fmt.Printf("X1: %f, X2: %f\n", x1, x2)
	stringArr := []byte(fmt.Sprintf("A: %f, B: %f, C: %f, x1: %f, x2: %f", coefA, coefB, coefC, x1, x2))
	// Permission: -rw-r--r--
	err := ioutil.WriteFile("/tmp/arq2.txt", stringArr, 0644)
	check(err)
}