![](./thumb.jpg)

# BDD em Golang com Godog

Esta é uma amostra grátis do meu curso **Engenharia de Testes de Software**, que tem exemplos em: **Java**, **Python** e **Golang**. Saiba mais no meu site: [**cleutonsampaio.com**](http://cleutonsampaio.com). Para esta amostra, o código está todo aqui, neste repositório, junto com este **README**.

Antes de começarmos, é bom deixar claro que a **BDD** ainda não está totalmente madura na plataforma **Golang**. Utilizamos o **godog** que seria baseado no **cucumber** da plataforma **Java**, mas ainda está na versão **0.14** na época em que fizemos os testes.

Behavior-Driven Development (**BDD**) é uma metodologia de desenvolvimento de software que visa melhorar a comunicação entre desenvolvedores, QA (Quality Assurance), e não-técnicos ou membros da equipe de negócios. O BDD foca em obter uma compreensão clara do comportamento desejado do software através de discussões com stakeholders antes de qualquer desenvolvimento começar, utilizando exemplos concretos para descrever o comportamento do sistema. Esses exemplos são então traduzidos em testes automatizados, garantindo que o software desenvolvido esteja alinhado com as expectativas dos stakeholders.

### Golang e BDD

Golang, também conhecido como Go, é uma linguagem de programação open source desenvolvida pela Google. É conhecida por sua simplicidade, eficiência e desempenho. Embora o Go não venha com uma biblioteca padrão específica para BDD, há várias ferramentas e frameworks disponíveis que permitem implementar BDD em projetos Go. Uma dessas ferramentas é o Godog.

### Godog

**Godog** é um framework para **BDD** em **Go** que se inspira no **Cucumber** (uma das ferramentas de BDD mais populares, originalmente escrita para Ruby). Ele permite que você escreva especificações de comportamento utilizando a linguagem **Gherkin**, que é uma linguagem de domínio específico (**DSL**) legível por humanos, permitindo a descrição de comportamentos de software sem detalhar como esse comportamento é implementado.

Godog serve como uma ponte entre os cenários escritos em Gherkin e o código Go que implementa esses comportamentos. Ao usar Godog, você pode definir cenários de teste em arquivos `.feature` utilizando a sintaxe Gherkin, que depois são ligados ao código Go que os valida. Isso ajuda a garantir que o software funcione conforme o esperado e facilita a comunicação sobre os requisitos do sistema.

### Como o Godog funciona?

1. **Escrever Cenários em Gherkin**: Você começa escrevendo cenários para os comportamentos desejados do seu sistema em arquivos `.feature`. Esses cenários são escritos em Gherkin, que usa uma linguagem natural para descrever os requisitos e as expectativas.

2. **Implementar os Passos**: Após definir os cenários, você implementa os passos em Go. Cada passo de um cenário é mapeado para uma função em Go que executa o teste associado àquele passo.

3. **Executar os Testes**: Com os cenários e passos definidos, você usa o Godog para executar os testes. O Godog lê os arquivos `.feature`, executa os passos correspondentes em Go e reporta os resultados.

4. **Refinar com Base nos Resultados**: Após a execução, você pode ver quais cenários passaram ou falharam. Isso permite ajustar o código conforme necessário para atender aos comportamentos esperados.

BDD com Golang e Godog pode ser uma abordagem poderosa para desenvolvimento de software, especialmente em equipes multidisciplinares. Ela permite a definição clara de requisitos através de uma linguagem natural, garantindo que todos os envolvidos no projeto tenham uma compreensão comum dos objetivos. Além disso, ao vincular esses requisitos diretamente aos testes automatizados, as equipes podem desenvolver com confiança, sabendo que o software atende às expectativas dos stakeholders.

## Pontos importantes a se considerar

Embora eu pude utilizar quase na íntegra o mesmo arquivo de **features**, houve alguns detalhes que podem comprometer seu prazo, se não tomar cuidado: 

1) Se uma coluna é **string** sempre a coloque entre aspas duplas dentro do arquivo de **feature**. Eu tive que mudar isso: 

De: 
```gherkin
    Then as formas de pagamento disponíveis devem ser <formas_pagamento>
```

Para: 
```gherkin
    Then as formas de pagamento disponíveis devem ser "<formas_pagamento>"
```

2) Expressões **cucumber** ainda não funcionam. Você não pode utilizar `{int}`, `{float}` ou `{string}`. Tem que utilizar expressões regulares:

No arquivo de passos, em vez de: 
```java
@Given("um cliente com {int} meses de cadastro, situação de crédito {string}, valor do último pedido {float} e último pagamento {float}")
```

Use: 
```golang
ctx.Step(`^um cliente com (\d+) meses de cadastro, situação de crédito "([^"]*)", valor do último pedido (\d+\.\d+) e último pagamento (\d+\.\d+)$`, tc.umClienteComMesesDeCadastroSituacaoDeCreditoValorDoUltimoPedidoEUltimoPagamento)
```

3) É preciso criar um código **bootstrap** para o comando `go test` conseguir executar os testes **godog**: 
```golang
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "pretty",                    // Use "pretty" or another supported format
		Output:    colors.Colored(os.Stdout),   // Use colored output
		Paths:     []string{"../features"},     // Path to feature files
		Randomize: time.Now().UTC().UnixNano(), // Optional: randomize scenario execution
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
```

4) Não utilizamos **anotações** para demarcar os métodos correspondentes aos **steps**, até porque não existe esse conceito de **annotations** nessa linguagem, embora possa utilizar comentários especiais em campos de **structs**. Então, precisamos criar uma função para inicializar os cenários. Ela é mencionada no código de **bootstrap** (**InitializeScenario**): 
```golang
    ...
	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()
```

E este é o código onde colocamos a expressão regular que invocará a função correspondente: 
```golang
func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := NewTestContext()
	fmt.Println("****************************************************Inicializando cenário")
	ctx.Step(`^um cliente com (\d+) meses de cadastro, situação de crédito "([^"]*)", valor do último pedido (\d+\.\d+) e último pagamento (\d+\.\d+)$`, tc.umClienteComMesesDeCadastroSituacaoDeCreditoValorDoUltimoPedidoEUltimoPagamento)
	ctx.Step(`^um pedido com valor de entrada (\d+\.\d+) e valor total (\d+\.\d+)$`, tc.umPedidoComValorDeEntradaEValorTotal)
	ctx.Step(`^verificar as formas de pagamento disponíveis$`, tc.verificarAsFormasDePagamentoDisponiveis)
	ctx.Step(`^as formas de pagamento disponíveis devem ser "([^"]*)"$`, tc.asFormasDePagamentoDisponiveisDevemSer)
}
```

## Nossa demonstração

Mais uma vez, mostrarei aqui as regras de negócio que desejamos testar: 
```code
1) Adicionar: 'Pagamento à vista'.
2) Se cliente.tempo_cliente >= 6 meses e cliente.situacao_credito != "ruim": 
3)      se cliente.situacao_credito = "boa" ou pedido.valor_entrada >= 20%:
4)          - Adicionar: 'Pagamento em 2 vezes com juros'
5)          Se cliente.situacao_credito == 'boa' e cliente.tempo_cliente >= 1 ano:
6)              - Adicionar: 'Pagamento em 3 vezes sem juros'
7)              se o cliente.valor_ultimo_pedido >= 1000 e cliente.valor_ultimo_pagamento >= 500:
8)                  - Adicionar: 'Pagamento em 6 vezes sem juros'
```

E aqui está a **tabela de decisão** que queremos testar: 

| Crédito | Tempo | Último pedido | Último pagamento | entrada | Vista | 2 x cj | 3 x sj | 6 x sj |
| ------- | ----- | ------------- | ---------------- | --- | --- | --- | --- | --- |
| qualquer | < 6 meses | qualquer | qualquer | qualquer | S | N | N | N |
| ruim | qualquer | qualquer | qualquer | qualquer | S | N | N | N |
| regular | 6 meses | qualquer | qualquer | <20% | S | N | N | N |
| regular | 6 meses | qualquer | qualquer | >=20% | S | S | N | N |
| boa | 6 meses | qualquer | qualquer | qualquer | S | S | N | N |
| boa | 1 ano | qualquer | qualquer | qualquer | S | S | S | N |
| boa | 1 ano | >= 1000 | >= 500 |  qualquer | S | S | S | S |

## Comece instalando o godog

Instale o **godog** e crie a variável **GOPATH** (mesmo que esteja utilizando **modulos**). Coloque o caminho no **PATH**:
```shell
go install github.com/cucumber/godog/cmd/godog@latest
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
```

## Altere o arquivo de features

Seu arquivo de **features** deve ficar assim: 
```gherkin
Feature: Verificar formas de pagamento disponíveis para o cliente

  Scenario Outline: Cliente elegível para múltiplas formas de pagamento
    Given um cliente com <tempo_cliente> meses de cadastro, situação de crédito "<situacao_credito>", valor do último pedido <valor_ultimo_pedido> e último pagamento <valor_ultimo_pagamento>
    And um pedido com valor de entrada <valor_entrada> e valor total <valor_pedido>
    When verificar as formas de pagamento disponíveis
    Then as formas de pagamento disponíveis devem ser "<formas_pagamento>"

    Examples:
      | tempo_cliente | situacao_credito | valor_ultimo_pedido | valor_ultimo_pagamento | valor_entrada | valor_pedido | formas_pagamento                                         |
      | 6             | boa              | 500.0                 | 250.0                    | 100.0           | 500.0          | Pagamento à vista, Pagamento em 2 vezes com juros    |
      | 12            | boa              | 1500.0                | 750.0                    | 300.0           | 1500.0         | Pagamento à vista, Pagamento em 2 vezes com juros, Pagamento em 3 vezes sem juros, Pagamento em 6 vezes sem juros |
      | 7             | regular          | 800.0                 | 400.0                    | 200.0           | 1000.0         | Pagamento à vista, Pagamento em 2 vezes com juros |
```

E coloque-o em uma pasta **features** logo abaixo da raiz do projeto. 

O arquivo de **Steps** deve ser como este: 
```golang
...
func (tc *TestContext) umClienteComMesesDeCadastroSituacaoDeCreditoValorDoUltimoPedidoEUltimoPagamento(tempoCliente int, situacaoCredito string, valorUltimoPedido float64, valorUltimoPagamento float64) error {
	// Aqui você configura o cliente com base nos parâmetros fornecidos

	tc.cliente = internal.Cliente{
		Id:              1,
		DataCadastro:    time.Now().AddDate(0, -tempoCliente, 0),
		SituacaoCredito: situacaoCredito,
	}
	tc.pedido = internal.Pedido{
		Numero: 1,
		Valor:  valorUltimoPedido,
	}
	tc.pagamentos = []internal.Pagamento{
		{
			Numero:    1,
			Valor:     valorUltimoPagamento,
			Liquidado: true,
		},
	}
	tc.mockPedidoService = &MockPedidoService{}
	tc.mockPagamentoService = &MockPagamentoService{}
	tc.mockPedidoService.On("GetPedidos", 1).Return([]internal.Pedido{{Numero: 1, Valor: 0}, {Numero: 2, Valor: valorUltimoPedido}})
	tc.mockPagamentoService.On("GetPagamentosCliente", 1).Return([]internal.Pagamento{{Numero: 1, Valor: 10, Liquidado: true},
		{Numero: 2, Valor: valorUltimoPagamento, Liquidado: true}})
	return nil
}

func (tc *TestContext) umPedidoComValorDeEntradaEValorTotal(valorEntrada, valorPedido float64) error {
	// Configure o pedido aqui
	tc.valorUltimoPedido = valorPedido
	tc.entrada = valorEntrada
	return nil
}

func (tc *TestContext) verificarAsFormasDePagamentoDisponiveis() error {
	// Aqui você chamaria a lógica para verificar as formas de pagamento
	parametros := internal.Parametros{ValorLimiteUltimoPedido: 1000, ValorLimiteUltimoPagamento: 500}
	formasPagamentoService := internal.NewFormasPagamentoService(tc.mockPedidoService, tc.mockPagamentoService, parametros)
	tc.resultado = formasPagamentoService.VerificarFormasPagamento(tc.cliente, tc.entrada, tc.pedido.Valor)

	return nil
}

func (tc *TestContext) asFormasDePagamentoDisponiveisDevemSer(formasPagamento string) error {
	// Verifique se o resultado é o esperado
	formasEsperadas := strings.Split(formasPagamento, ", ")
	// Aqui você faria a asserção de que tc.resultado é igual a formasEsperadas
	if !reflect.DeepEqual(formasEsperadas, tc.resultado) {
		return fmt.Errorf("esperado: %v, obtido: %v", formasEsperadas, tc.resultado)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := NewTestContext()
	fmt.Println("****************************************************Inicializando cenário")
	ctx.Step(`^um cliente com (\d+) meses de cadastro, situação de crédito "([^"]*)", valor do último pedido (\d+\.\d+) e último pagamento (\d+\.\d+)$`, tc.umClienteComMesesDeCadastroSituacaoDeCreditoValorDoUltimoPedidoEUltimoPagamento)
	ctx.Step(`^um pedido com valor de entrada (\d+\.\d+) e valor total (\d+\.\d+)$`, tc.umPedidoComValorDeEntradaEValorTotal)
	ctx.Step(`^verificar as formas de pagamento disponíveis$`, tc.verificarAsFormasDePagamentoDisponiveis)
	ctx.Step(`^as formas de pagamento disponíveis devem ser "([^"]*)"$`, tc.asFormasDePagamentoDisponiveisDevemSer)
}
...
```

Aqui temos a inicialização do **cenário** (**InitializeScenario**), que aponta cada **Step** da **feature** para uma função (no nosso caso, método) correspondente. E temos a implementação dos métodos que executarão os **steps**.

Eu crie uma **struct** com campos para armazenar os resultados intermediários dos passos (**steps**): 
```golang
type TestContext struct {
	cliente              internal.Cliente
	pedido               internal.Pedido
	pagamentos           []internal.Pagamento
	entrada              float64
	valorUltimoPedido    float64
	resultado            []string
	mockPedidoService    *MockPedidoService
	mockPagamentoService *MockPagamentoService
}
```

E todas as funções de **step** são métodos dessa **struct**.

Este código deve ficar ao lado do código principal, na mesma pasta, porém em pacote diferente (**internal_test**), como é costume em **golang**. 

## Como executar os testes

Meu código todo está dentro de **internal**, então, ajustamos o caminho da pasta onde está a **feature** no **bootstrap**: 
```golang
...
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "pretty",                    // Use "pretty" or another supported format
		Output:    colors.Colored(os.Stdout),   // Use colored output
		Paths:     []string{"../features"},     // Path to feature files
        ...
```

Entramos na pasta **internal** e executamos: 
```shell
go test 

Feature: Verificar formas de pagamento disponíveis para o cliente
****************************************************Inicializando cenário
      | 12            | boa              | 1500.0              | 750.0                  | 300.0         | 1500.0       | Pagamento à vista, Pagamento em 2 vezes com juros, Pagamento em 3 vezes sem juros, Pagamento em 6 vezes sem juros |
****************************************************Inicializando cenário
      | 7             | regular          | 800.0               | 400.0                  | 200.0         | 1000.0       | Pagamento à vista, Pagamento em 2 vezes com juros                                                                 |
****************************************************Inicializando cenário

  Scenario Outline: Cliente elegível para múltiplas formas de pagamento                                                                                                                        # ../features/verificar_formas_pagamento.feature:3
    Given um cliente com <tempo_cliente> meses de cadastro, situação de crédito "<situacao_credito>", valor do último pedido <valor_ultimo_pedido> e último pagamento <valor_ultimo_pagamento> # <autogenerated>:1 -> *TestContext
    And um pedido com valor de entrada <valor_entrada> e valor total <valor_pedido>                                                                                                            # <autogenerated>:1 -> *TestContext
    When verificar as formas de pagamento disponíveis                                                                                                                                          # <autogenerated>:1 -> *TestContext
    Then as formas de pagamento disponíveis devem ser "<formas_pagamento>"                                                                                                                     # <autogenerated>:1 -> *TestContext

    Examples:
      | tempo_cliente | situacao_credito | valor_ultimo_pedido | valor_ultimo_pagamento | valor_entrada | valor_pedido | formas_pagamento                                                                                                  |
      | 6             | boa              | 500.0               | 250.0                  | 100.0         | 500.0        | Pagamento à vista, Pagamento em 2 vezes com juros                                                                 |

3 scenarios (3 passed)
12 steps (12 passed)
2.277632ms

Randomized with seed: 1709638137043744375
testing: warning: no tests to run
PASS
ok      cleutonsampaio.com/golangdemo1/internal 0.010s
```



