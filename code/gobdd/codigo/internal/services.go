package internal

import "time"

type Parametros struct {
	ValorLimiteUltimoPedido    float64
	ValorLimiteUltimoPagamento float64
}

type Cliente struct {
	Id              int
	DataCadastro    time.Time
	SituacaoCredito string
}

type Pedido struct {
	Numero int
	Valor  float64
}

type Pagamento struct {
	Numero    int
	Valor     float64
	Liquidado bool
}

type PedidoServiceImpl struct{}

func (ps *PedidoServiceImpl) getPedidos(clienteID int) []Pedido {
	// Implementação fictícia de obtenção de pedidos
	return []Pedido{}
}

type PagamentoServiceImpl struct{}

func (ps *PagamentoServiceImpl) getPagamentosCliente(clienteID int) []Pagamento {
	// Implementação fictícia de obtenção de pagamentos de um cliente
	return []Pagamento{}
}
