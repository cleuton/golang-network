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
