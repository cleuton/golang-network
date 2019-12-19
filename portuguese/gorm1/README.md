![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

## Banco de dados com gORM

![](../../images/godb.png)

Acessar banco de dados é uma dor de cabeça para os programadores em qualquer linguagem. Técnicas como [**ORM** - mapeamento objeto relacional](https://pt.wikipedia.org/wiki/Mapeamento_objeto-relacional) ajudam a melhorar a produtividade, reduzindo a [**complexidade acidental**](http://www.obomprogramador.com/2012/12/complexidade-acidental.html) das interfaces com [**SGBDs**](https://pt.wikipedia.org/wiki/Sistema_de_gerenciamento_de_banco_de_dados).

Para a linguagem **Go** temos o popular pacote [**gorm**](http://gorm.io/), que faz este papel muito bem. Porém, devo avisar que o **gorm** possui *idiossincrasias* que podem ser um pouco perturbadoras, para quem está acostumado com outros pacotes **ORM**, como o [**Hibernate**](https://hibernate.org/) para a linguagem **Java**.  

Infelizmente, o site do **gorm** não possui exemplos claros o suficiente e é por isso que eu resolvi escrever este post. A intenção é complementar as informações e mostrar um exemplo mais próximo do que é esperado pelos desenvolvedores. 

Os exemplos comuns utilizam um banco [**sqlite**](https://www.sqlite.org/index.html) local, coisa que considero bem diferente da realidade dos desenvolvedores de software corporativo, que utilizam um servidor remoto de banco de dados. Portanto, resolvi inovar e usar um [**PostgreSQL**](https://www.postgresql.org/) remoto, hospedado no serviço [**ElephantQL**](https://customer.elephantsql.com/instance). Para rodar o exemplo, você precisará ter uma conta no **ElephantQL** e criar uma instância de banco lá. Não se preocupe, pois tem um nível gratuito. 

## Criando o banco de dados

O **gorm** assume que você utilizará uma chave primária chamada **ID** que é do tipo inteiro. Embora ele funcione com quaisquer tipos de chaves primárias, alguns comandos (como o **Find()**) funcionam de maneira diferente quando o tipo de dados da chave primária é string, por exemplo. Desta maneira, resolvi usar o meu banco, com uma tabela que já está criada nele. 

O **gorm** também assume que você vai mapear o banco de dados conforme as suas structs, o que é incorreto. Geralmente, em ambientes corporativos, você terá que trabalhar com um banco de dados pré-existente (a não ser em testes). Portanto, vou usar um banco de dados pré-existente.

Depois de criar uma conta no **ElephantQL**, crie uma instância de banco de dados. O nível gratuito, chamado de **Tiny turtle** só permite criar uma única instância, e não tem instrumentos de administração. Você terá que criar sua tabela usando SQL. Abra a opção **browser** e execute a **query** a seguir: 

```
CREATE TABLE pessoa (
    cpf             varchar(11) CONSTRAINT pkey PRIMARY KEY,
    nome            varchar(50) NOT NULL,
    data_nascimento date NOT NULL,
    funcionario     boolean
);
```

Depois, abra o menu **details** e copie o DNS do servidor (**server**), o nome do usuário e do banco (**User & Default database**), a senha (**Password**). No campo **URL** você pode ver qual é a porta TCP, geralmente 5432.

Crie um código **Go** para testar o acesso, como este: 

```
package main

import (
    "fmt"
    "time"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

    // Open connection to a postgresql database running on ElephantQL:
    db, err := gorm.Open("postgres", "host=elmer.db.elephantsql.com port=5432 user=userdbname dbname=userdbname password=password")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()
    fmt.Println("OK")
}
```

Os campos **user** e **dbname** são iguais.

Substitua os campos no método **Open** de acordo com o que copiou lá no **ElephantQL**. Se tudo estiver certo, você verá um **OK** na console. 

O que é o **defer db.Close()**? O comando **defer** coloca o comando seguinte em uma pilha. Podemos colocar vários comandos nessa pilha. Após a função retornar (terminar), os comandos serão retirados da pilha (FIFO) e executados. Assim, estou empilhando o fechamento da conexão com o SGBD.

O [**exemplo completo**](https://github.com/cleuton/golang-network/tree/master/code//gorm1/godb.go) tem um [**CRUD**](https://pt.wikipedia.org/wiki/CRUD) básico para você ver como funcionam os comandos. 

Antes de começar, você precisa instalar o pacote do **gorm**: 

```
go get github.com/jinzhu/gorm
```

## Mapeando tabelas

Como é o **ORM**, o **gorm** mapeia tabelas SQL em **structs**. Por isso, precisamos definir nossas structs, como esta:

```
type Pessoa struct {
    Cpf string `gorm:"primary_key;type:varchar(11);column:cpf"`
    Name string `gorm:"type:varchar(50);column:nome"`
    BirthDate *time.Time `gorm:"column:data_nascimento"`
    Employee bool `gorm:"column:funcionario"`
}
```
O que são essas anotações depois das declarações dos campos, cercadas por acentos graves? São [**anotações** ou **struct tags**](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go). São declarações extras sobre os elementos, que podem ser lidas e processadas por outro código escrito em **Go**. Semelhantes às anotações da linguagem **Java**.

Podemos declarar muitas coisas nas [**struct tags**](http://gorm.io/docs/models.html), como: chave primária, tipo de dados e nome da coluna na tabela (se for diferente do nome na struct). Eu declarei o campo que é a chave primária, além dos nomes e tipos das colunas.

## Inserindo dados

Para começar, o **gorm** suporta o uso de transações, e eu estou usando isso para inserir os dados: 

```
tx := db.Begin()
dt,_ := time.Parse("2006-01-02", "1979-08-18")
person1 := &Pessoa{"111","Person#1",&dt,false}
tx.Table("public.pessoa").Create(&person1)
person2 := &Pessoa{"222","Person#2",&dt,false}
tx.Table("public.pessoa").Create(&person2)
person3 := &Pessoa{"333","Person#3",&dt,false}
tx.Table("public.pessoa").Create(&person3)
...
tx.Commit()
```

Cada comando é executado dentro de uma transação automaticamente. Se quisermos rodar vários comandos dentro de uma única transação, podemos fazer, como eu fiz no exemplo. O método **Create()** cria um registro no banco, usando os dados da **struct**.

Como eu criei uma conexão com o banco usando a variável **db**, poderia simplesmente escrever: ```db.Create()```. Mas há dois problemas: 

1. O nome da tabela. O **gorm** assume que o nome da tabela é sempre no plural do nome da struct. Se o nome for diferente, você terá que informar o **schema** e o **table name**, usando o método **Table()**;
2. Estou usando transação na variável **tx**, portanto, tenho que usá-la, em vez de **db**.

Existe também o comando **tx.Rollback()** para desfazer uma transação. 

## Selecionando registros

Podemos selecionar registros de várias maneiras, como o [**exemplo**](https://github.com/cleuton/golang-network/tree/master/code//gorm1/godb.go): 

```
var person Pessoa
if result := db.Table("public.pessoa").First(&person); result.Error != nil {
    panic(result.Error)
}
fmt.Println("First person: ",person)
```

Há várias coisas aqui. Para começar, o método **First()** seleciona o primeiro registro de um conjunto. Como não especifiquei condição alguma, será o primeiro registro do banco de dados, de acordo com a chave primária. 

Também estou mostrando como tratar os erros. A função retornará uma **struct** que tem um campo **Error**. Se ele for nulo, então o comando funcionou. Caso contrário, eu posso usar o **panic** e terminar o programa. Isto funciona com o **Create()** também. 

Eu mostrei outras formas de selecinar registros, por exemplo, mais de um registro: 

```
var persons []Pessoa
if result := db.Table("public.pessoa").Find(&persons); result.Error != nil {
    panic(result.Error)
}
fmt.Printf("Persons: %+v\n",persons)
```

O método **Find()** retorna todos os registros encontrados. Como não especifiquei filtro, serão todos mesmo. Por isso, preciso de um **slice** para receber os registros.

Também mostrei como selecionar registros de acordo com um filtro: 

```
var person222 Pessoa
if result := db.Table("public.pessoa").Where("cpf like ?", "222").First(&person222); result.Error != nil {
    panic(result.Error)
}
fmt.Println("Person identified by cpf: ",person222)
```

Aqui usei o método **Where()** para selecionar o primeiro elemento cujo cpf começa por 222. 

## Atualização

A atualização de registros é bem intuitiva: 

```
fmt.Println("Before updating: ",person3.BirthDate)
if result := db.Table("public.pessoa").Where("cpf like ?", "333").First(&person3); result.Error != nil {
    panic(result.Error)
}
newdate, _ := time.Parse("2006-01-02", "1990-11-01")
if result := db.Table("public.pessoa").Model(&person3).Update("data_nascimento", &newdate); result.Error != nil {
    panic(result.Error)
}
fmt.Println("Person data updated",person3.BirthDate)
```

Aqui eu selecionei o registro e atualizei o campo "data_nascimento". 

## Deleção

Deletar registros é igualmente simples: 

```
if result := db.Table("public.pessoa").Where("cpf like ?", "111").First(&person1); result.Error != nil {
    panic(result.Error)
}
fmt.Println(person1)
if result := db.Table("public.pessoa").Delete(&person1); result.Error != nil {
    panic(result.Error)
}
fmt.Println("Person deleted")

```

## O que faltou?

Bom, existem as [**associações**](http://gorm.io/docs/belongs_to.html) entre tabelas, as [**chaves compostas**](http://gorm.io/docs/composite_primary_key.html) e outras características avançadas, mas, com essa base, você pode pequisar na documentação do **gorm**.

