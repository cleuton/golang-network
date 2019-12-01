![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Modularização do código

Modularizar um código-fonte é dividi-lo em camadas de abstração diferentes. Temos uma camada de mais alto nível e uma ou mais camadas de baixo nível, contendo funções ou classes. 

## Funções

Uma **função** em **Go** pode ou não receber argumentos e pode ou não retornar um valor. Funções que nada retornam são conhecidas como **side-effect functions** pois alteram o estado do programa, já que nada retornam.

```
func delta(a float64, b float64, c float64) float64 {
	return math.Pow(float64(b), 2) - 4*a*c
}
```

A função **delta** recebe 3 argumentos reais (float 64) e retorna um valor real (float64). Ela pode ser invocada assim: 

```
vDelta := delta(coefA, coefB, coefC)
```

O programa [**function.go**](./codigo/function.go) demonstra o uso de funções: 

```
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
```

**Lendo argumentos de linha de comando**

Note que usamos o pacote **os** para ler argumentos da linha de comando. A função **readFromArgs()** faz isso com o comando: 

```
clArguments := os.Args[1:]
...
f, err := strconv.ParseFloat(a[p], 64)
```

O objeto **strconv** converte dados para strings e vice-versa. O método **ParseFloat** converte strings em floats, e podemos especificar a precisão, neste caso 64 bits. Nós pegamos o vetor **os.Args[]** com um **slice** a partir da segunda posição (os.Args[1:]), para evitar pegar o nome do programa, que é sempre o primeiro argumento. Dai para a frente, teremos um slice começando em zero, com cada argumento passado após o nome do programa. 

Outra coisa interessante a notar é a exponenciação, provida pelo método **Pow()** do objeto **math**. Este método exige que as variáveis sejam **float64** e é por isso que estou utilizando este tipo de dados. 

**Tratamento de erros**

Em **Go** não usamos **SEH** (Structured Exception Handling) como fazemos em C++, Java ou Python (try). Nós simplesmente testamos se houve algum erro na função. Funções podem retornar múltiplos valores, por exemplo: 

```
	
func f() (int, int) {
    return 1, 2
}
...
a, b := f()
```

Algumas funções também retornam erros. Veja só como eu testo o erro de conversão de string para float: 

```
f, err := strconv.ParseFloat(a[p], 64)
if err != nil {
    log.Fatal(err)
}
```

O método **ParseFloat()** tem a assinatura: 

```
func ParseFloat(s string, bitSize int) (float64, error)
```

Portanto, ele pode retornar um erro. Eu posso testar se o erro é **nil**, ou seja, nulo. Se for, então não existe erro. Porém, se houver um erro, eu uso o método **Fatal()** do **log** que mostra a mensagem e paraliza a execução.

Para executar: 

```
go run function.go 1 -5 6
```

Porém, para executar esse código, você necessitará do **Go** instalado, pois no [**MyCompiler**](https://www.mycompiler.io/) não conseguirá passar argumentos de linha de comando. 


## Instalando Go

[**Instalar o Go**](https://golang.org/doc/install) é muito simples. Faça de acordo com o seu sistema operacional. Aqui, estou usando **Ubuntu**: 

Baixe o [**pacote**](https://golang.org/dl/) e extraia em /usr/local, crie uma árvore em /usr/local/go. Por exemplo: 

```
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```
Escolha o pacote apropriado para sua instalação. Por exemplo, se estiver instalando a versão Go 1.2.1, para Linux x86 64 bits, o pacote é **go1.2.1.linux-amd64.tar.gz**.

(Estes comandos devem ser executados com **sudo**)

Adicione /usr/local/go/bin à variável de ambiente **PATH**. 

Se utiliza **Windows** ou **MacOS** [consulte instruções](https://golang.org/doc/install) específicas. 

## Comando Go

Execute o comando **Go**: 

```
go
...
$ go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildmode   build modes
        c           calling between Go and C
        cache       build and test caching
        environment environment variables
        filetype    file types
        go.mod      the go.mod file
        gopath      GOPATH environment variable
        gopath-get  legacy GOPATH go get
        goproxy     module proxy protocol
        importpath  import path syntax
        modules     modules, module versions, and more
        module-get  module-aware go get
        module-auth module authentication using go.sum
        module-private module configuration for non-public modules
        packages    package lists and patterns
        testflag    testing flags
        testfunc    testing functions

Use "go help <topic>" for more information about that topic.
```

Dos dois comandos básicos são: ```build``` e ```run```. Para executar o código exemplo, [**function.go**](./codigo/function.go), abra a pasta em um terminal e digite: 

```
go run function.go 1 -5 6
...
1
```

## Desafio

Complete o programa criando uma função que calcule as raízes da equação de segundo grau. Lembre-se que o valor do delta é determinante para isso!

Dica: use math.Inf para retornar infinito (se conhece equação do segundo grau, sabe que isso poderá acontecer)