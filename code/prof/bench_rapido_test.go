package main

import (
	"strings"
	"testing"

	"flamegraph-go/contapalavras"
)

func BenchmarkContarPalavrasRapido(b *testing.B) {
	// workload maior para garantir amostras suficientes
	texto := strings.Repeat("Go é ótimo para performance e Go é rápido!\n", 10000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = contapalavras.ContarPalavrasRapido(texto)
	}
}
