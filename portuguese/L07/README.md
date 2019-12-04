![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Pacotes

Os programadores **Go** organizam seu código em uma única **workspace**. Uma pasta, geralmente dentro de sua pasta pessoal, apontada pela variável de ambiente **GOPATH**. Eu uso uma pasta assim: 

```
echo $GOPATH
/home/cleuton/go
```

Você pode atualizar a variável de ambiente dentro de: ```/etc/bash.bashrc``` e ```/etc/profile.d```.

A estrutura do **GOPATH** é assim: 

```
go 
+-src: projetos por repositório e pacote
+-bin: comandos (executáveis compilados)
```

Se você baixar um pacote Go, ele será instalado dentro de uma subpasta dentro de **src**: 

```
go get -v golang.org/x/text
```

Se você analisar o **GOPATH** verá que foi baixado o código-fonte em **src**: 

**/home/cleuton/go/src/golang.org/x/text/...**

Este é o código-fonte do pacote baixado. E veremos uma pasta **pkg** que contém o compilado: 

**/home/cleuton/go/pkg/linux_amd64/golang.org/x/text.a**

Como já deve ter notado, os pacotes são agrupados pelo hostname do repositório onde estão. Neste caso, o pacote **text** está em **golang.org**. Não há controle de versão e nem repositório central de pacotes, como em outras linguagens. 

Se você quiser criar um repositório **Go** para os outros usarem, é só publica-lo no **Github** e fornecer seu endereço: github.com/**username**/**nome do repositório**.

## Pacotes executáveis e bibliotecas

Tudo em **Go** deve estar dentro de um pacote. Por exemplo, vamos criar um pacote dentro do $GOPATH com as funções que usamos para calcular a equação do segundo grau (veja o [**arquivo bhaskara**](./codigo/bhaskara.go)).

Dentro de $GOPATH/src crie uma pasta chamada "bhaskarautils" e copie este arquivo para lá. Você terá uma estrutura assim: **$GOPATH/src/bhaskarautils/bhaskara.go**. 

Agora, rode o comando abaixo: 

```
go install bhaskarautils
```

E você verá que foi criado um arquivo **$GOPATH/pkg/linux_amd64/bhaskarautils.a**. Agora, podemos rodar o código [**testbhaskara.go**](./codigo/test/testbhaskara.go) sem problemas. 

```
import (
    ...
	"bhaskarautils"
)
...
	deltaV := bhaskarautils.Delta(coefA, coefB, coefC)
	x1, x2 := bhaskarautils.Roots(deltaV, coefA, coefB)
```

Eu tive que renomear as funções **Delta()** e **Roots()** para inicial maiúscula, pois assim o pacote as exportará. Caso contrário, serão consideradas privadas. 

O pacote **bhaskarautils** é uma biblioteca, pois não contém arquivos com a **func main()**, portanto, será gerado um arquivo de biblioteca (extensão .a). 

Se compilarmos um executável, geraremos um arquivo executável nativo do sistema operacional: 

```
go build testbhaskara.go
```

E podemos executá-lo como fazemos com qualquer programa.

## Repositórios

O ideal é criar seu projeto em um repositório, como o Github, e importá-lo para sua workspace com ```go get <pacote>```.

## Este é o fim?

Claro que não. É apenas o fim deste curso básico. Sei que há muito mais a aprender, porém, acredito que agora você tem os instrumentos indispensáveis para continuar sozinho. Boa sorte e acompanhe [**golang.network**](http://golang.network).