#!/bin/bash

set -e

RESULTS_DIR="results"
mkdir -p $RESULTS_DIR

# Função para executar teste
run_test() {
    local SERVER_NAME=$1
    local BUILD_CMD=$2
    local EXECUTABLE_PATH=$3 # Novo parâmetro para o caminho do executável

    echo "Iniciando teste para $SERVER_NAME"

    # Construir o servidor
    echo "Construindo $SERVER_NAME..."
    cd servers/$SERVER_NAME
    if [ -n "$BUILD_CMD" ]; then
        eval $BUILD_CMD
    fi
    cd ../..

    # Determinar o caminho do executável
    if [ -n "$EXECUTABLE_PATH" ]; then
        EXECUTABLE="$EXECUTABLE_PATH"
    elif [ "$SERVER_NAME" == "rust_axum" ]; then
        EXECUTABLE="./servers/$SERVER_NAME/target/release/zaptidgen-rust"
    elif [ "$SERVER_NAME" == "nodejs_express" ]; then
         # Para Node.js, o executável é o próprio comando node
         EXECUTABLE="node"
    else
        EXECUTABLE="./servers/$SERVER_NAME/zaptidgen"
    fi

    # Verificar se o executável existe (exceto para Node.js onde o comando é 'node')
    if [ "$SERVER_NAME" != "nodejs_express" ] && [ ! -f "$EXECUTABLE" ]; then
        echo "ERRO: Executável não encontrado: $EXECUTABLE"
        return 1
    fi

    # Iniciar servidor em background
    echo "Iniciando servidor..."
    if [ "$SERVER_NAME" == "nodejs_express" ]; then
        # Para Node.js, executar o cluster.js
        $EXECUTABLE ./servers/$SERVER_NAME/cluster.js > /dev/null 2>&1 &
    else
        $EXECUTABLE > /dev/null 2>&1 &
    fi
    SERVER_PID=$!

    # Esperar servidor iniciar
    sleep 5

    # Testar se o servidor está respondendo
    if ! curl -s --max-time 10 http://localhost:8888/nextid > /dev/null; then
        echo "ERRO: Servidor $SERVER_NAME não está respondendo"
        kill $SERVER_PID 2>/dev/null || true
        return 1
    fi

    # Iniciar monitoramento de recursos (do diretório raiz)
    echo "Iniciando monitoramento de recursos..."
    ./monitoring/resource_monitor.sh "$SERVER_NAME" 120 > /dev/null 2>&1 &
    MONITOR_PID=$!

    # Executar teste K6 com limite de tempo
    echo "Executando teste K6..."
    timeout 130s k6 run --out csv=$RESULTS_DIR/${SERVER_NAME}_k6_results.csv k6_tests/load_test.js || true

    # Parar servidor e monitor
    echo "Parando servidor e monitor..."
    kill $SERVER_PID 2>/dev/null || true
    kill $MONITOR_PID 2>/dev/null || true

    # Esperar processos terminarem
    sleep 2

    # Salvar tamanho do executável ou da pasta node_modules
    if [ "$SERVER_NAME" == "nodejs_express" ]; then
        # Calcular o tamanho da pasta node_modules para Node.js
        # find ./servers/$SERVER_NAME -name "node_modules" -prune -o -type f -exec du -b {} + 2>/dev/null | awk '{sum += $1} END {print sum}' > $RESULTS_DIR/${SERVER_NAME}_binary_size.txt
        # Opção mais direta: calcular o tamanho total da pasta node_modules
        if [ -d "./servers/$SERVER_NAME/node_modules" ]; then
            NODE_MODULES_SIZE=$(du -sb ./servers/$SERVER_NAME/node_modules | cut -f1)
            echo "$NODE_MODULES_SIZE" > $RESULTS_DIR/${SERVER_NAME}_binary_size.txt
        else
             # Se node_modules não existir (erro na instalação?), salva 0
             echo "0" > $RESULTS_DIR/${SERVER_NAME}_binary_size.txt
        fi
    else
        BINARY_SIZE=$(stat -f%z $EXECUTABLE 2>/dev/null || stat -c%s $EXECUTABLE)
        echo "$BINARY_SIZE" > $RESULTS_DIR/${SERVER_NAME}_binary_size.txt
    fi

    echo "Teste concluído para $SERVER_NAME"
    echo "----------------------------------------"
}

# Compilar e executar testes
echo "Iniciando bateria de testes..."

run_test "go_fasthttp" "go build -o zaptidgen ." ""
# run_test "rust_axum" "cargo build --release" "" # Comentado conforme solicitado
run_test "nodejs_express" "npm install" "node" # Passa 'node' como o executável

echo "Todos os testes concluídos!"
