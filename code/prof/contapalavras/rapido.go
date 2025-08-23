package contapalavras

import (
	"regexp"
	"strings"
)

var (
	rxPalavras = regexp.MustCompile(`[[:alpha:][:digit:]_]+`)
	rxLimpeza  = regexp.MustCompile(`[^[:alpha:][:digit:]_]+`)
)

// versão corrigida: compila regex uma vez só
func ContarPalavrasRapido(texto string) map[string]int {
	contador := make(map[string]int)
	for _, linha := range strings.Split(texto, "\n") {
		for _, p := range rxPalavras.FindAllString(linha, -1) {
			limpa := rxLimpeza.ReplaceAllString(p, "")
			minusc := strings.ToLower(limpa)
			contador[minusc]++
		}
	}
	return contador
}
