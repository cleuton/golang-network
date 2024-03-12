package internal

type PedidoService interface {
	GetPedidos(clienteID int) []Pedido
}

type PagamentoService interface {
	GetPagamentosCliente(clienteID int) []Pagamento
}
