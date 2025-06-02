// vault_example_test.go

/*
Teste básico da biblioteca em um arquivo local
*/
package vault_test

import (
	"os"
	"testing"

	"github.com/cleutonsampaio/senhas/vault"
)

func TestVaultLocal(t *testing.T) {
	path := "test_vault.json"
	os.Remove(path)
	storage := vault.FileStorage{}

	// 1) cria cofre v1
	v1, err := vault.CreateVault(path, "Senha123", storage)
	if err != nil {
		t.Fatal(err)
	}

	// 2) adiciona locais
	if err := v1.AddLocal("siteA", "userA", "passA"); err != nil {
		t.Fatal(err)
	}
	if err := v1.AddLocal("siteB", "userB", "passB"); err != nil {
		t.Fatal(err)
	}

	// 3) lista e valida quantos locais existem
	locs, err := v1.ListLocais()
	if err != nil {
		t.Fatal(err)
	}
	if len(locs) != 2 {
		t.Fatalf("esperado 2 locais, achou %d", len(locs))
	}

	// 4) recupera dados de um local
	u, p, err := v1.GetCredenciais("siteA")
	if err != nil {
		t.Fatal(err)
	}
	if u != "userA" || p != "passA" {
		t.Fatal("credenciais incorretas")
	}

	// 5) exporta cleartext
	all, err := v1.ExportClear()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatal("export falhou")
	}

	// 6) inicializa cofre v2 vazio
	path2 := path + "2"
	os.Remove(path2)
	v2, err := vault.CreateVault(path2, "Senha123", storage)
	if err != nil {
		t.Fatal(err)
	}

	// 7) sincroniza v1 -> v2
	if err := v1.Sync(v2); err != nil {
		t.Fatal(err)
	}

	// 8) valida sincronia
	locs2, err := v2.ListLocais()
	if err != nil {
		t.Fatal(err)
	}
	if len(locs2) != 2 {
		t.Fatal("sync falhou")
	}
}
