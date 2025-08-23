package contapalavras

import (
	"fmt"
	"regexp"
	"strings"
)

// propositalmente ineficiente
func ContarPalavrasLento(texto string) map[string]int {
	contador := make(map[string]int)
	linhas := strings.Split(texto, "\n")
	for _, linha := range linhas {
		regexPalavras, _ := regexp.Compile(`[[:alpha:][:digit:]_]+`)
		palavras := regexPalavras.FindAllString(linha, -1)
		for _, p := range palavras {
			regexLimpeza, _ := regexp.Compile(`[^[:alpha:][:digit:]_]+`)
			limpa := regexLimpeza.ReplaceAllString(p, "")
			minusc := strings.ToLower(fmt.Sprintf("%s", limpa)) // alocação inútil
			contador[minusc]++
		}
	}
	return contador
}
