// vault.go
// Biblioteca para gerenciamento de cofre de senhas em Go.
// Independente de armazenamento, separa lógica de persistência.

package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

// Interface de persistência: salvar e carregar o JSON bruto do cofre.
type Storage interface {
	Save(path string, data []byte) error
	Load(path string) ([]byte, error)
}

// estrutura interna do cabeçalho
type header struct {
	Salt     string `json:"salt"`
	Iter     int    `json:"iteracoes"`
	TagCheck string `json:"tag_check"`
}

// estrutura de cada entrada cifrada
type entry struct {
	ID     string `json:"id"`
	IV     string `json:"iv"`
	Cipher string `json:"cipher"`
}

// representação do cofre em disco
type vaultFile struct {
	Cabecalho header  `json:"cabecalho"`
	Entradas  []entry `json:"entradas"`
}

// Vault é o objeto ativo em memória
// guarda o conteúdo e a chave de cifra
type Vault struct {
	file     vaultFile
	key      []byte // Kenc
	backend  Storage
	filePath string
}

// Deriva chaves de autenticação e cifra a partir da senha-mestre
func deriveKeys(senha []byte, salt []byte, iter int) (kauth, kenc []byte) {
	// PBKDF2-HMAC-SHA256
	mk := pbkdf2.Key(senha, salt, iter, 32, sha256.New)
	// HKDF-Extract
	extract := hkdf.New(sha256.New, mk, nil, nil)
	kauth = make([]byte, 32)
	io.ReadFull(extract, kauth)
	// Expand para kenc
	expander := hkdf.New(sha256.New, mk, nil, []byte("cifragem"))
	kenc = make([]byte, 32)
	expander.Read(kenc)
	return
}

// Cria um novo cofre e persiste
func CreateVault(path string, senha string, storage Storage) (*Vault, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	kauth, kenc := deriveKeys([]byte(senha), salt, 200000)
	// tag de verificação
	h := hmac.New(sha256.New, kauth)
	h.Write([]byte("CHECK_VAULT_V1"))
	tag := h.Sum(nil)
	vf := vaultFile{
		Cabecalho: header{
			Salt:     base64.StdEncoding.EncodeToString(salt),
			Iter:     200000,
			TagCheck: base64.StdEncoding.EncodeToString(tag),
		},
		Entradas: []entry{},
	}
	data, _ := json.MarshalIndent(vf, "", "  ")
	if err := storage.Save(path, data); err != nil {
		return nil, err
	}
	return &Vault{file: vf, key: kenc, backend: storage, filePath: path}, nil
}

// Abre um cofre existente, valida senha (sem decifrar entradas)
func OpenVault(path string, senha string, storage Storage) (*Vault, error) {
	raw, err := storage.Load(path)
	if err != nil {
		return nil, err
	}
	var vf vaultFile
	if err := json.Unmarshal(raw, &vf); err != nil {
		return nil, err
	}
	// decodifica salt e tag
	salt, _ := base64.StdEncoding.DecodeString(vf.Cabecalho.Salt)
	tagStored, _ := base64.StdEncoding.DecodeString(vf.Cabecalho.TagCheck)
	kauth, kenc := deriveKeys([]byte(senha), salt, vf.Cabecalho.Iter)
	// verifica tag
	h := hmac.New(sha256.New, kauth)
	h.Write([]byte("CHECK_VAULT_V1"))
	if !hmac.Equal(tagStored, h.Sum(nil)) {
		return nil, errors.New("senha-mestre incorreta")
	}
	return &Vault{file: vf, key: kenc, backend: storage, filePath: path}, nil
}

// lista todos os locais (nome decifrado)
func (v *Vault) ListLocais() ([]string, error) {
	res := []string{}
	for _, e := range v.file.Entradas {
		plain, err := v.decryptEntry(e)
		if err != nil {
			return nil, err
		}
		res = append(res, plain.Local)
	}
	return res, nil
}

// recupera usuario e senha de um local
func (v *Vault) GetCredenciais(local string) (user, pass string, err error) {
	for _, e := range v.file.Entradas {
		plain, errD := v.decryptEntry(e)
		if errD != nil {
			continue
		}
		if plain.Local == local {
			return plain.Usuario, plain.Senha, nil
		}
	}
	return "", "", errors.New("local não encontrado")
}

// cria nova entrada
func (v *Vault) AddLocal(local, usuario, senha string) error {
	eid := make([]byte, 16)
	rand.Read(eid)
	iv := make([]byte, 12)
	rand.Read(iv)
	data := map[string]string{"local": local, "usuario": usuario, "senha": senha}
	jsonB, _ := json.Marshal(data)
	block, _ := aes.NewCipher(v.key)
	aesgcm, _ := cipher.NewGCM(block)
	cipherText := aesgcm.Seal(nil, iv, jsonB, eid)
	v.file.Entradas = append(v.file.Entradas, entry{
		ID:     base64.StdEncoding.EncodeToString(eid),
		IV:     base64.StdEncoding.EncodeToString(iv),
		Cipher: base64.StdEncoding.EncodeToString(cipherText),
	})
	return v.persist()
}

// altera usuario/senha de um local existente
func (v *Vault) UpdateLocal(local, novoUser, novaSenha string) error {
	// remove e re-adiciona
	if err := v.DeleteLocal(local); err != nil {
		return err
	}
	return v.AddLocal(local, novoUser, novaSenha)
}

// apaga um local
func (v *Vault) DeleteLocal(local string) error {
	novo := []entry{}
	found := false
	for _, e := range v.file.Entradas {
		plain, err := v.decryptEntry(e)
		if err == nil && plain.Local == local {
			found = true
			continue
		}
		novo = append(novo, e)
	}
	if !found {
		return errors.New("local não encontrado")
	}
	v.file.Entradas = novo
	return v.persist()
}

// exporta todas as entradas em texto claro
func (v *Vault) ExportClear() ([]map[string]string, error) {
	out := []map[string]string{}
	for _, e := range v.file.Entradas {
		plain, err := v.decryptEntry(e)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]string{
			"local":   plain.Local,
			"usuario": plain.Usuario,
			"senha":   plain.Senha,
		})
	}
	return out, nil
}

// sincroniza entradas do vault v com outro target
func (v *Vault) Sync(target *Vault) error {
	locs, err := v.ListLocais()
	if err != nil {
		return err
	}
	for _, local := range locs {
		u, s, _ := v.GetCredenciais(local)
		target.DeleteLocal(local)
		target.AddLocal(local, u, s)
	}
	return nil
}

// descarrega e persiste no backend
func (v *Vault) persist() error {
	data, _ := json.MarshalIndent(v.file, "", "  ")
	return v.backend.Save(v.filePath, data)
}

// decryptEntry retorna estrutura interna
func (v *Vault) decryptEntry(e entry) (struct{ Local, Usuario, Senha string }, error) {
	id, _ := base64.StdEncoding.DecodeString(e.ID)
	iv, _ := base64.StdEncoding.DecodeString(e.IV)
	ct, _ := base64.StdEncoding.DecodeString(e.Cipher)
	block, _ := aes.NewCipher(v.key)
	aesgcm, _ := cipher.NewGCM(block)
	plain, err := aesgcm.Open(nil, iv, ct, id)
	if err != nil {
		return struct{ Local, Usuario, Senha string }{}, err
	}
	var d map[string]string
	json.Unmarshal(plain, &d)
	return struct{ Local, Usuario, Senha string }{d["local"], d["usuario"], d["senha"]}, nil
}

// ------------------ exemplo de storage local ------------------

type FileStorage struct{}

func (fs FileStorage) Save(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0600)
}
func (fs FileStorage) Load(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
