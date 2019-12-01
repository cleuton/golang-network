package main

import (
	"fmt"
	"time"
)

func main() {

	timeNow := time.Now()
	fmt.Println("Current time: ", timeNow.String())
}
