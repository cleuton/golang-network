#!/bin/bash

# Monitor de recursos do sistema
# Uso: ./resource_monitor.sh <nome_do_teste> <duracao_segundos>

TEST_NAME=$1
DURATION=$2
OUTPUT_FILE="./results/${TEST_NAME}_monitoring.csv"

echo "timestamp,cpu_percent,memory_mb" > $OUTPUT_FILE

for i in $(seq 1 $DURATION); do
    # CPU usage (%)
    CPU=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    
    # Memory usage (MB)
    MEM=$(free -m | grep "Mem:" | awk '{print $3}')
    
    TIMESTAMP=$(date +%s)
    echo "$TIMESTAMP,$CPU,$MEM" >> $OUTPUT_FILE
    
    sleep 1
done