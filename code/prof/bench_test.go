package main

import (
	"strings"
	"testing"

	"flamegraph-go/contapalavras"
)

func BenchmarkContarPalavrasLento(b *testing.B) {
	texto := strings.Repeat("Go é ótimo para performance e Go é rápido!\n", 2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = contapalavras.ContarPalavrasLento(texto)
	}
}
