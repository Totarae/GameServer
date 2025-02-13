package main

import (
	"awesomeProject5/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/ws", handlers.HandleWebSocket)
	go handlers.GameLoop()
	go handlers.HeartbeatLoop()

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// TODO: Добавить graceful shutdown
}
