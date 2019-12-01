package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x uint8 = 10
	// An error will happen if you uncomment the following line:
	// y = 11
	y := 11
	fmt.Println(reflect.TypeOf(x), x, reflect.TypeOf(y), y)

	var z float32 = 22 / 7.0
	pi := 22 / 7.0
	fmt.Println(reflect.TypeOf(z), z, reflect.TypeOf(pi), pi)

	var name string = "John"
	surname := " Doe"
	fmt.Println(reflect.TypeOf(name), name, reflect.TypeOf(surname), surname)
	fmt.Println(name + surname)

	message :=
		`This is a multiline string, 
	because there are 2 lines of text. `
	fmt.Println(message)

	other := "this is also a\nmultiline message."
	fmt.Println(other)

	var (
		p bool = true
		q bool = false
	)
	r := p && !q
	fmt.Println(r)
	fmt.Println("p XOR q", p != q)
	fmt.Println("p XOR r", p != r)

	const tipec = "*"
	fmt.Println(tipec)

	var grades = [5]float32{5.5, 7.5, 8.0, 5.3, 9.2}
	fmt.Println(grades)
	for ix, val := range grades {
		fmt.Println("Index", ix, "value", val)
	}
	grades[1] = 8.0
	fmt.Println(grades[1])

	firstOnes := grades[0:2]
	lastOnes := grades[3:]
	fmt.Println("firstOnes", firstOnes)
	fmt.Println("lastOnes", lastOnes)

	for i := 0; i < 5; i++ {
		fmt.Println("grades", i, grades[i])
	}
}
