![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Datatypes, atribuição, blocos

Go possui um rico conjunto de tipos de dados, entre eles: 

**Números inteiros**

| Tipo | Descrição |
|------------|--------|
| uint8 | 8 bits sem sinal |
| uint16 | 16 bits sem sinal |
| uint32 | 32 bits sem sinal |
| uint64 | 64 bits sem sinal |
| int8   | 8 bits com sinal |
| int16  | 16 bits com sinal |
| int32  | 32 bits com sinal |
| int64  | 64 bits com sinal |

Vejamos um exemplo de declaração de variável com tipo declarado e inferido: 

```
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
    fmt.Println(reflect.TypeOf(x),x,reflect.TypeOf(y),y)
}
...
uint8 10 int 11
```

Neste código, vemos a declaração de duas variáveis: "x", inteiro não sinalizado de 8 bits, e "y", inteiro (**int**). O tipo **int** tem o tamanho dependente da arquitetura do computador onde você está compilando o código. 

Mais interessante é que vemos duas maneiras de declararmos variáveis: Com tipo explícito e implícito. Quando é necessário explicitar o tipo, usamos o comando **var**: 

```
var <nome 1>, <nome 2> <tipo de dados> = <valor inicial 1>, <valor inicial 2>
```

Já na outra declaração, o tipo de dados é inferido pelo valor atribuído. Neste caso, precisamos acrescentar um caractere dois pontos (":") antes do sinal de atribuição, indicando que estamos declarando e inicializando uma variável. Note que a inicialização não é obrigatória quando usamos **var**. Eis a declaração implícita: 

```
y := 11
```

Se você tentar declarar a variável como faz em **Python**, ou seja, simplesmente atribuindo o valor, o compilador dará um erro indicando que a variável não foi declarada. 

Outra coisa interessante foi como eu mostrei o tipo de dados de cada variável: o pacote "reflect", que, entre outras coisas, possui o método "TypeOf()". Para usá-lo, tive que declará-lo na lista de **import**.


**Números reais**

| Tipo | Descrição |
|------------|--------|
| float32 | 32 bits |
| float64 | 64 bits |

Veja o exemplo no [**arquivo tipos.go**](./codigo/tipos.go):

```
var z float32 = 22 / 7.0
pi := 22 / 7.0
fmt.Println(reflect.TypeOf(z), z, reflect.TypeOf(pi), pi)
...
float32 3.142857 float64 3.142857142857143
```

Para começar, eu calculei **PI** dividindo 22 por 7. Note que usei decimal no algarismo 7. Isto foi para transformá-lo em **float**, caso contrário, daria um resultado inteiro. 

Note que a primeira variável, **z**, foi declarada como **float32**, e a segunda, **pi**, foi declarada com tipo implícito, e o compilador assumiu **float64**. Note como a precisão do valor calculado foi maior. Quanto maior o tipo flutuante, maior a precisão. 

**Texto**

Em **Go** temos o tipo **string** para declarar sequências de caracteres: 

```
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
... 
string Fulano string  de Tal
Fulano de Tal
Esta é uma mensagem multilinhas, 
	pois abrange mais de uma linha física. 
Esta também é uma mensagem
multilinha.
```

Podemos declarar literais **string** com aspas duplas ou acentos graves ("`"). Quando delimitamos com acentos graves, podemos expandir o texto por múltiplas linhas. Quando usamos aspas duplas, podemos incluir "\n" para nova linha ou "\t" para tabulação. E podemos concatenar strings como fazemos em **Java**.

**Variáveis lógicas**

Em **Go** temos o tipo **bool** para declararmos variáveis lógicas. E temos as operações lógicas: 

- Conjunção (AND): && 
- Disjunção (OR): || 
- Negação (NOT): !
- Disjunção exclusiva (XOR): Não tem. 

```
var (
    p bool = true
    q bool = false
)
r := p && !q
fmt.Println(r)
fmt.Println("p XOR q", p != q)
fmt.Println("p XOR r", p != r)
...
true
p XOR q true
p XOR r false
```

Aproveitei para mostrar como podemos criar um bloco **var** com múltiplas variáveis sendo declaradas. Eu simulei o **XOR** com um "truque" (p!=q).

**Constantes**

Constantes podem ser declaradas com o comando **const**: 

```
const tipo = "*"
fmt.Println(tipo)
...
*
```

**Arrays e slices**

**Go** tem um jeito especial para lidarmos com variáveis multivaloradas: o **array**, e podemos retirar **fatias** dele com **slicing**. Nota, o primeiro elemento de um array é zero, e o último é **tamanho - 1**.

Declaramos arrays informando o tamanho, entre chaves, e o tipo dos valores. E podemos inicializá-los: 

```
var notas = [5]float32{5.5, 7.5, 8.0, 5.3, 9.2}
```

Podemos declarar **slices** indicando a posição inicial e a posição final (aberta): 

```
primeiras := notas[0:2]
```

Pegamos do primeiro elemento ao terceiro (exclusive).

E podemos declarar omitindo uma das informações: 

```
ultimas := notas[3:]
```

Neste exemplo, pegamos da quarta posição até o final.

E podemos iterar sobre arrays usando o comando **for** com a opção **range**. A cada iteração, receberemos duas variáveis novas: O índice e o valor daquela posição: 

```
for ix, val := range notas {
    fmt.Println("Index", ix, "value", val)
}
```

Eis o trecho completo: 

```
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
...
[5.5 7.5 8 5.3 9.2]
Index 0 value 5.5
Index 1 value 7.5
Index 2 value 8
Index 3 value 5.3
Index 4 value 9.2
8
Primeiras [5.5 8]
Ultimas [5.3 9.2]
notas 0 5.5
notas 1 8
notas 2 8
notas 3 5.3
notas 4 9.2
```

É possível usar a forma tradicional do **for** e também o comando **break**: 

```
for i := 0; i < 5; i++ {
    fmt.Println("notas", i, notas[i])
}
```

**Comparação**

Para finalizar esta lição, temos o comando **if**: 

```
if condição {

}
else {

}
```

Os operadores de comparação em **Go** são:  ```==, !=, <, <=, > e >=```

## Desafio

Dado um vetor contendo "n" inteiros, indique a maior quantidade repetida seguida de números. Exemplo: 

```
{1,1,0,0,0,3,1,1,4,4,4,0,0,7,7,7,7,7,1}
```
Neste exemplo, a resposta seria: 5, pois o número 7 se repete 5 vezes. 

