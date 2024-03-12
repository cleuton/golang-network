package internal_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	internal "cleutonsampaio.com/golangdemo1/internal"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/mock"
)

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

func NewTestContext() *TestContext {
	return &TestContext{}
}

// Mock para o serviço de Pedido
type MockPedidoService struct {
	mock.Mock
}

func (m *MockPedidoService) GetPedidos(clienteID int) []internal.Pedido {
	args := m.Called(clienteID)
	return args.Get(0).([]internal.Pedido)
}

// Mock para o serviço de Pagamento
type MockPagamentoService struct {
	mock.Mock
}

func (m *MockPagamentoService) GetPagamentosCliente(clienteID int) []internal.Pagamento {
	args := m.Called(clienteID)
	return args.Get(0).([]internal.Pagamento)
}

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

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	// Aqui você pode inicializar coisas necessárias antes da suite de teste ser executada
}

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
