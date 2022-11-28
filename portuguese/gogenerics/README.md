![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# GENERICS

Generics é um recurso importante em linguagens de programação modernas e foi introduzido no Golang na versão 1.18, em Março de 2022.

O uso de Generics permite criar funções (ou métodos) que funcionam com tipos diferentes de dados. 

Abra o [**código deste tutorial**](../../code/gogenerics/cmd/main.go) e vamos examinar como isso funciona.

```
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
``` 

A função **Between** utiliza Generics e pode receber qualquer tipo de dados que imponha ordenação. O pacote **constraints** tem vários tipos restrições que podemos utilizar. Neste caso, a função Between precisa receber argumentos de tipos que possam ser comparados com ">=" e "<=".

A assinatura da função diz que o tipo dos argumentos, "T", tem que ser um tipo de dados ordenável: 

```
func Between[T constraints.Ordered](a, b, x T) bool
```

Ela recebe 3 argumentos, todos do tipo "T" e retorna um "bool". O que ela faz é comparar se o argumento "x" está no intervalo fechado entre "a" e "b".

Para utilizar uma função com Generics você precisa fazer como em **Java**, ou seja, tem que informar o tipo da função que vai usar: 

``` 
if Between[int](2, 10, 5)
``` 
Agora, se você instanciar a função, não precisa informar o tipo de dados: 

```
beetween := Between[int]
fmt.Println(beetween(2, 10, 8))
```

É um recurso muito importante e você o verá em muitos pacotes daqui para a frente.