package main

import (
	"log"
	"my_chatbot/infrastructures"
	"my_chatbot/messagehandlers"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

func initRoutes(mux *chi.Mux) {
	mux.Post("/send_message", messagehandlers.SendMessageHandler)
	mux.Get("/get_message", messagehandlers.GetMessageHandler)
	mux.Post("/ws", messagehandlers.WebsocketHandleConnection)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
}

func main() {
	infrastructures.InitializeDatabase()
	routingSystem := chi.NewRouter()
	initRoutes(routingSystem)

	// Start the server
	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", routingSystem)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
