package contapalavras

import "testing"

func TestContarPalavrasLento(t *testing.T) {
	texto := "Go é ótimo para performance.\nGo é rápido!"
	r := ContarPalavrasLento(texto)
	if r["go"] < 2 {
		t.Fatalf("esperava pelo menos 2 ocorrências de 'go', obtive %d", r["go"])
	}
}
