package main

import "fmt"

func main() {
	var vetor = []int16{1, 1, 0, 0, 0, 3, 1, 1, 4, 4}
	n := len(vetor)
	max := 0
	count := 1
	for i := 1; i < n; i++ {
		if vetor[i] == vetor[i-1] {
			count++
			if count >= max {
				max = count
			}
		} else {
			count = 1
		}
	}
	fmt.Println(max)
}
