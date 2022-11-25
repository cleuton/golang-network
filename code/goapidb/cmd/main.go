package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// Tipos utilizados:

type Candidato struct {
	Id        int
	Nome      string
	CreatedAt time.Time
}

type CriarCandidato struct {
	Nome string
}

// Conexão com o database:

func connectDB() (*sql.DB, error) {
	// coloquei os dados de conexão em variáveis de ambiente
	// lembre-se do host, port, user, password e database name que você usou
	host := getEnv("DEMO_HOST", "localhost")
	dbPort := getEnv("DEMO_DBPORT", "5432")
	dbUser := getEnv("DEMO_USER", "postgres")
	password := getEnv("DEMO_DATABASE_PASSWORD", "mysecretpassword")
	dbName := getEnv("DEMO_DBNAME", "postgres")
	sslMode := getEnv("DEMO_SSLMODE", "disable")
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, dbPort, dbUser, password, dbName, sslMode)
	db, err := sql.Open("postgres", connectString)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}

// Tipos e funções auxiliares:

type Handlers struct {
	Db *sql.DB
}

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if body != nil {
		payload, _ := json.Marshal(body)
		w.Write(payload)
	}
}

// Handlers das rotas REST:

// Deleta um candidato:
func (h *Handlers) DeleteCandidatoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	sql := `DELETE FROM candidatos WHERE id = $1`
	res, err := h.Db.Exec(sql, id)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	contagem, err := res.RowsAffected()
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	if contagem == 0 {
		WriteResponse(http.StatusNotFound, map[string]string{"Erro": "Nenhum registro encontrado"}, w)
		return
	}

	WriteResponse(http.StatusNoContent, nil, w)
}

// Atualiza candidato (só o nome):
func (h *Handlers) UpdateCandidatoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var novo CriarCandidato
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&novo); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": "invalido"}, w)
		return
	}
	sql := `UPDATE candidatos SET nome = $2 WHERE id = $1`
	res, err := h.Db.Exec(sql, id, novo.Nome)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	contagem, err := res.RowsAffected()
	if err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	if contagem == 0 {
		WriteResponse(http.StatusNotFound, map[string]string{"Erro": "Nenhum registro encontrado"}, w)
		return
	}

	WriteResponse(http.StatusOK, map[string]string{"alterado": novo.Nome}, w)
}

// Criar novo candidato:
func (h *Handlers) CreateCandidatoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var novo CriarCandidato
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&novo); err != nil {
		WriteResponse(http.StatusBadRequest, map[string]string{"error": "invalido"}, w)
		return
	}
	sql := `INSERT INTO candidatos (nome, created_at) VALUES($1,$2)`
	if _, err := h.Db.Exec(sql, novo.Nome, time.Now()); err != nil {
		WriteResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}, w)
		return
	}
	WriteResponse(http.StatusCreated, nil, w)
}

// Listar os candidatos:
func (h *Handlers) CandidatosHandlerFunc(w http.ResponseWriter, r *http.Request) {
	candidatos, err := h.Db.Query(`SELECT * FROM candidatos`)
	codigo := http.StatusOK
	if err != nil {
		fmt.Println(err)
		message := map[string]string{"error": "erro ao acessar banco de dados"}
		switch err {
		case sql.ErrNoRows:
			codigo = http.StatusNotFound
			message["cause"] = "tabela vazia"
		default:
			codigo = http.StatusInternalServerError
			message["cause"] = "erro geral no banco de dados"
		}
		WriteResponse(codigo, message, w)
		return
	}
	var lista []Candidato
	for candidatos.Next() {
		var candidato Candidato
		if err := candidatos.Scan(&candidato.Id, &candidato.Nome, &candidato.CreatedAt); err != nil {
			panic(err)
		}
		lista = append(lista, candidato)
	}

	WriteResponse(codigo, lista, w)
}

func main() {
	porta := getEnv("DEMO_PORT", "8080")
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	h := Handlers{Db: db}
	router := mux.NewRouter()
	router.HandleFunc("/candidatos", h.CandidatosHandlerFunc).Methods("GET")
	router.HandleFunc("/candidato", h.CreateCandidatoHandlerFunc).Methods("POST")
	router.HandleFunc("/candidato/{id}", h.UpdateCandidatoHandlerFunc).Methods("PUT")
	router.HandleFunc("/candidato/{id}", h.DeleteCandidatoHandlerFunc).Methods("DELETE")
	err = http.ListenAndServe(fmt.Sprintf(":%s", porta), router)
	fmt.Println(err)
}
