package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x uint8 = 10
	// Se descomentar a linha abaixo dará erro:
	// y = 11
	y := 11
	fmt.Println(reflect.TypeOf(x), x, reflect.TypeOf(y), y)

	var z float32 = 22 / 7.0
	pi := 22 / 7.0
	fmt.Println(reflect.TypeOf(z), z, reflect.TypeOf(pi), pi)

	var nome string = "Fulano"
	sobrenome := " de Tal"
	fmt.Println(reflect.TypeOf(nome), nome, reflect.TypeOf(sobrenome), sobrenome)
	fmt.Println(nome + sobrenome)

	mensagem :=
		`Esta é uma mensagem multilinhas, 
	pois abrange mais de uma linha física. `
	fmt.Println(mensagem)

	outra := "Esta também é uma mensagem\nmultilinha."
	fmt.Println(outra)

	var (
		p bool = true
		q bool = false
	)
	r := p && !q
	fmt.Println(r)
	fmt.Println("p XOR q", p != q)
	fmt.Println("p XOR r", p != r)

	const tipo = "*"
	fmt.Println(tipo)

	var notas = [5]float32{5.5, 7.5, 8.0, 5.3, 9.2}
	fmt.Println(notas)
	for ix, val := range notas {
		fmt.Println("Index", ix, "value", val)
	}
	notas[1] = 8.0
	fmt.Println(notas[1])

	primeiras := notas[0:2]
	ultimas := notas[3:]
	fmt.Println("Primeiras", primeiras)
	fmt.Println("Ultimas", ultimas)

	for i := 0; i < 5; i++ {
		fmt.Println("notas", i, notas[i])
	}
}
