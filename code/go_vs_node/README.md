<img src="./go_vs_node.png" height=400>

---

<img src="../../golangnetwork-logo.png" height=300>

---

[**Cleuton Sampaio**](https://linkedin.com/in/cleutonsampaio)

# Go vs NodeJS

Muita gente está descobrindo **Golang** como alternativa de alto desempenho, baixo consumo de memória e baixo **footprint** para aplicações web e APIs **REST**. E realmente é! Vejo muitos posts e relatos de pessoas convertendo apps de **NodeJS** e **Python** para **Go**. Mas... Será que realmente há algum ganho ao migrar de NodeJS para Go? Não é apenas "achismo"? 

> **Footprint** em software refere-se ao consumo de recursos do sistema, como memória RAM, CPU e espaço em disco, que uma aplicação utiliza durante sua execução. Em ambientes de nuvem, isso é crucial porque os custos são tipicamente baseados no consumo desses recursos. Um menor footprint significa menor uso de memória, processamento e armazenamento, o que se traduz diretamente em menores custos operacionais. Além disso, aplicações com menor footprint tendem a ser mais rápidas, escaláveis e eficientes, permitindo que mais serviços rodem na mesma infraestrutura e respondam melhor sob carga. Isso é especialmente importante em nuvem onde se paga por recurso consumido e a eficiência impacta diretamente no orçamento e performance da solução.

**Vejamos o que dizem o benchmark**

O site [**Benchmarks game**](https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go-node.html) tem uma estatística comparativa, que eu resumirei aqui: 

1.  **Desempenho:** Go demonstra desempenho significativamente superior na maioria dos benchmarks computacionais (ex: `fannkuch-redux`, `n-body`, `mandelbrot`, `k-nucleotide`, `regex-redux`). Os programas em Go geralmente têm tempos de execução (`secs`) muito menores comparados aos programas em Node.js. Node.js aparece consistentemente entre os piores colocados em termos de velocidade de execução para esses tipos de tarefas.
2.  **Uso de Memória:** Go utiliza muito menos memória RAM do que Node.js para os mesmos benchmarks. Os valores de memória (`mem`) para Node.js estão frequentemente na casa de centenas de megabytes ou até gigabytes (ex: 111,055 KB, 448,524 KB, 1,897,566 KB), enquanto os programas Go operam com dezenas ou algumas centenas de megabytes (ex: 3,846 KB, 168,186 KB, 416,195 KB).
3.  **Eficiência de CPU:** Benchmarks computacionalmente intensivos mostram que programas Go otimizados podem utilizar múltiplos núcleos de CPU de forma mais eficaz, muitas vezes com `cpu secs` próximos ou superiores ao tempo real (`secs`), indicando bom paralelismo. Node.js, sendo single-threaded com um loop de eventos, não escala da mesma forma nesses cenários.
4.  **Implantação (Tamanho do Artefato):**
    *   **Go:** Produz um único **binário executável estático** após a compilação. O tamanho desse binário é relativamente pequeno.
    *   **Node.js:** Requer a distribuição do **código-fonte da aplicação** e da pasta **`node_modules`** que contém todas as dependências de terceiros. O tamanho da pasta `node_modules` pode variar significativamente. 


Este teste de performance é um comparativo que usei para dar subsídios a um cliente que está indeciso quanto à migração de NodeJS para Golang. Na verdade, este é um teste de geração de IDs com **Sony Snowflake**, implementado em **Golang** (usando `FastHTTP`) e em **Javascript** (usando `NodeJS` e `Express`).

Os códigos-fonte estão na [**pasta `servers`**](./performance-test/servers/).

Usei o [**K6**](https://k6.io/) para executar os testes e criei scripts em `shell` e `python` para coletar os resultados e métricas. 

Os parâmetros de teste que usei foram: 

- Ramp Up: 30 segundos até 50 VUs.
- Teste de 1 minuto.
- Ramp Down: 30 segundos até 0. 
- Thinking time de 0.05s.

**Workers / processos**

1.  **Go (`go_fasthttp`):**
    *   O código Go usando `fasthttp` não cria explicitamente múltiplos processos ou threads para lidar com requisições concorrentes no código fornecido. A função `fasthttp.ListenAndServe` internamente lida com concorrência, mas o número exato de goroutines ou threads de trabalho (workers) não é fixado no código e é gerenciado pela própria biblioteca `fasthttp` de forma otimizada. Geralmente, bibliotecas Go como `fasthttp` usam um modelo de concorrência leve (goroutines) e um pool de threads do sistema operacional, escalando conforme a demanda. O número de "workers" efetivos é dinâmico e não especificado como um número fixo no código.

2.  **Node.js (`nodejs_express`):**
    *   O código Node.js foi modificado para usar o módulo `cluster`. A linha `const numCPUs = require('os').cpus().length;` determina o número de CPUs disponíveis no sistema.
    *   O laço `for (let i = 0; i < numCPUs; i++) { cluster.fork(); }` cria explicitamente **um processo worker filho para cada núcleo de CPU** disponível.
    *   Portanto, se o sistema onde o teste foi executado tem, por exemplo, 4 núcleos de CPU, a versão Node.js iniciou **4 workers** (além de 1 processo master).

*   **Go (`go_fasthttp`):** O número de workers não é explicitamente definido no código; é gerenciado internamente pela biblioteca `fasthttp`.
*   **Node.js (`nodejs_express`):** O número de workers é igual ao número de núcleos de CPU da máquina onde o teste foi executado (por exemplo, 4 workers em uma máquina de 4 núcleos).

Só roda em Linux, por enquanto!

## Para rodar os testes: 

1) Instale as dependências **Python** (aconselho a criar um ambiente virtual com `venv`):

```
pip install -r requirements.txt
```
2) Instale o k6: https://k6.io/docs/getting-started/installation/

3) Dê permissão para os scripts

```
chmod +x run_tests.sh
chmod +x monitoring/resource_monitor.sh
```

## Executar testes

```
./run_tests.sh
```

## Analisar resultados

```
python3 analysis/analyze_results.py
```

## Resultados finais

Após executar várias vezes, cheguei a esses resultados: 

Com base na análise do benchmark fornecido e nos resultados específicos do seu teste de desempenho entre Go (`go_fasthttp`) e Node.js (`nodejs_express`), aqui está um resumo das comparações mais importantes:

**Resumo da Comparação: Go vs. Node.js**

Este teste específico comparou `go_fasthttp` com `nodejs_express` (usando clustering). Os resultados mostraram consistentemente que, para este caso de uso (geração de IDs sequenciais via API REST):

1.  **Desempenho (Tempo de Resposta):** Go (`go_fasthttp`) foi significativamente mais rápido, com tempos de resposta médios quase 3 vezes menores (484 ms vs. 1299 ms).
2.  **Eficiência de CPU:** Go utilizou menos CPU em média (3.63% vs. 7.47%) e teve picos máximos muito menores (7.90% vs. 25.70%).
3.  **Eficiência de Memória:** Go consumiu menos memória RAM em média (6682 MB vs. 7327 MB) e demonstrou uso de memória muito mais estável (desvio padrão de 22 MB vs. 244 MB).
4.  **Implantação:** Go produz um único executável binário compacto (7.28 MB). Node.js requer o ambiente de execução (`node`) e as dependências (`node_modules`, que somaram 2.14 MB neste caso, mas requerem o runtime Node.js para funcionar).

Então, **Go** provou ser uma opção de menor consumo, menor footprint e melhor desempenho do que Javascript. Note que escolhemos uma aplicação que nem é **CPU intensive** para esse teste, cujo consumo de memória deveria ser pequeno. 

