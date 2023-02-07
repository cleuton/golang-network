package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"com.blocopad/blocopad_poc/internal/backend"
	"com.blocopad/blocopad_poc/internal/db"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

// Counters

var (
	NewNotes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "demo_new_notes_total",
			Help: "Total number of new notes requests",
		},
	 )
	 
	GetNotes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "demo_get_notes_total",
		Help: "Total number of get notes requests",
	},)

	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_requests_per_status",
			Help: "Total requests per http status code",
		},
		[]string{"code"},
	)
)

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

// HTTP Handlers

func ReadNote(w http.ResponseWriter, r *http.Request) {
	GetNotes.Inc()
	vars := mux.Vars(r)
	id := vars["id"]
	statusCode := 200
	if data, err := backend.GetKey(id); err == nil {
		WriteResponse(statusCode, data, w)
	} else {
		if err.Error() == "not found" {
			statusCode = 404
			WriteResponse(404, "Note not found", w)
		} else {
			statusCode = 500
			WriteResponse(500, "Error", w)
		}
	}
	totalRequests.WithLabelValues(strconv.Itoa(statusCode)).Inc()
}

func WriteNote(w http.ResponseWriter, r *http.Request) {
	NewNotes.Inc()
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var note db.Note
	statusCode := 200
	if err := decoder.Decode(&note); err != nil {
		statusCode = http.StatusBadRequest
		WriteResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}, w)
		return
	}
	uuidString, err := backend.SaveKey(note.Text, note.OneTime)
	if err != nil {
		fmt.Println(err)
		statusCode = http.StatusBadRequest
		WriteResponse(http.StatusBadRequest, map[string]string{"error": "invalid request"}, w)
	} else {
		WriteResponse(200, map[string]string{"code": uuidString}, w)
	}
	totalRequests.WithLabelValues(strconv.Itoa(statusCode)).Inc()
}

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  next.ServeHTTP(w, r)
	})
  }

func main() {
	serverPort := "8080"
	if port, hasValue := os.LookupEnv("API_PORT"); hasValue {
		serverPort = port
	}
	databaseUrl := "localhost:6379"
	if dbUrl, hasValue := os.LookupEnv("API_DB_URL"); hasValue {
		databaseUrl = dbUrl
	}
	databasePassword := ""
	if dbPassword, hasValue := os.LookupEnv("API_DB_PASSWORD"); hasValue {
		databasePassword = dbPassword
	}

	fmt.Printf("\nAPI_PORT: %s, API_DB_URL: %s\n", serverPort, databaseUrl)
	db.DatabaseUrl = databaseUrl
	db.DatabasePassword = databasePassword
	prometheus.MustRegister(NewNotes)
	prometheus.MustRegister(GetNotes)
	prometheus.MustRegister(totalRequests)
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)
	router.Path("/metrics").Handler(promhttp.Handler())
	router.HandleFunc("/api/note/{id}", ReadNote).Methods("GET")
	router.HandleFunc("/api/note", WriteNote).Methods("POST")
	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), router)
	fmt.Println(err)

}
