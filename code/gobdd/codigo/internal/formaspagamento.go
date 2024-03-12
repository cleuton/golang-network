package internal

import (
	"sort"
	"time"
)

type FormasPagamentoService struct {
	pedidoService    PedidoService
	pagamentoService PagamentoService
	parametros       Parametros
}

func NewFormasPagamentoService(pedidoService PedidoService, pagamentoService PagamentoService, parametros Parametros) *FormasPagamentoService {
	return &FormasPagamentoService{pedidoService, pagamentoService, parametros}
}

func (fps *FormasPagamentoService) VerificarFormasPagamento(cliente Cliente, valorEntrada float64, valorPedido float64) []string {
	if valorEntrada < 0 {
		panic("verificarFormasPagamento: Valor entrada inválido")
	}
	if valorPedido < 0 {
		panic("verificarFormasPagamento: Valor do pedido inválido")
	}
	if valorEntrada > valorPedido {
		panic("verificarFormasPagamento: Valor entrada maior que o valor do pedido")
	}

	formasPagamento := []string{}

	pedidosCliente := fps.pedidoService.GetPedidos(cliente.Id)
	sort.SliceStable(pedidosCliente, func(i, j int) bool {
		return pedidosCliente[i].Numero < pedidosCliente[j].Numero
	})
	pagamentos := fps.pagamentoService.GetPagamentosCliente(cliente.Id)

	sort.SliceStable(pagamentos, func(i, j int) bool {
		return pagamentos[i].Numero < pagamentos[j].Numero
	})

	valorUltimoPagamento := fps.getUltimoPagamento(pagamentos)
	valorUltimoPedido := fps.getValorUltimoPedido(pedidosCliente)
	tempoCliente := fps.getTempoCliente(cliente.DataCadastro)
	valorEntradaMaiorIgual20Porcento := valorEntrada >= valorPedido*0.2

	formasPagamento = append(formasPagamento, "Pagamento à vista")

	if tempoCliente >= 6 && cliente.SituacaoCredito != "ruim" {
		if cliente.SituacaoCredito == "boa" {
			formasPagamento = append(formasPagamento, "Pagamento em 2 vezes com juros")
			if tempoCliente >= 12 {
				formasPagamento = append(formasPagamento, "Pagamento em 3 vezes sem juros")
				if valorUltimoPedido >= fps.parametros.ValorLimiteUltimoPedido &&
					valorUltimoPagamento >= fps.parametros.ValorLimiteUltimoPagamento {
					formasPagamento = append(formasPagamento, "Pagamento em 6 vezes sem juros")
				}
			}
		} else if valorEntradaMaiorIgual20Porcento {
			formasPagamento = append(formasPagamento, "Pagamento em 2 vezes com juros")
		}
	}

	return formasPagamento
}

func (fps *FormasPagamentoService) getValorUltimoPedido(pedidosCliente []Pedido) float64 {
	valor := 0.0
	if len(pedidosCliente) > 0 {
		valor = pedidosCliente[len(pedidosCliente)-1].Valor
	}
	return valor
}

func (fps *FormasPagamentoService) getTempoCliente(dataCadastro time.Time) int {
	dataFinal := time.Now()
	anos := int(dataFinal.Year()) - int(dataCadastro.Year())
	mesesCalculo := 12
	if dataFinal.Month() == dataCadastro.Month() &&
		dataFinal.Day() < dataCadastro.Day() {
		mesesCalculo -= 1
	}
	diferenca := int((int(dataFinal.Month()) - int(dataCadastro.Month())) + mesesCalculo*anos)
	return diferenca
}

func (fps *FormasPagamentoService) getUltimoPagamento(pagamentos []Pagamento) float64 {
	valor := 0.0
	for i := len(pagamentos) - 1; i >= 0; i-- {
		if pagamentos[i].Liquidado {
			valor = pagamentos[i].Valor
			break
		}
	}
	return valor
}
