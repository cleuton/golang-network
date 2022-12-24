package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	ws = websocket.Upgrader{}
)

func WriteResponse(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(body)
	w.Write(payload)
}

// HTTP Handlers

func StartController(w http.ResponseWriter, r *http.Request) {
	ws.CheckOrigin = func(r *http.Request) bool { return true }
	wsocket, err := ws.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade failed: ", err)
		return
	}

	defer wsocket.Close()

	commandsLoop(wsocket)
}

// Process websocket commands

func commandsLoop(wsocket *websocket.Conn) {
	for {
		msgType, message, err := wsocket.ReadMessage()
		if err != nil {
			fmt.Println("Error or connection closed:", err)
			break
		}
		command := string(message)
		response := sendToDrone(command)

		returnMessage := []byte(response)
		err = wsocket.WriteMessage(msgType, returnMessage)
		if err != nil {
			fmt.Println("Error or connection closed:", err)
			break
		}
	}
}

// Send command to drone (fake)

func sendToDrone(command string) string {
	fmt.Println("Sending comand to drone...")
	return "- Drone OK (" + command + ")!"
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/controller", StartController).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	fmt.Println(err)

}
