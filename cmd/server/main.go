package main

import (
	"log"
	"net/http"

	"anon-skrzynka/app/config"
	apphttp "anon-skrzynka/app/http"
	"anon-skrzynka/app/storage"
	"anon-skrzynka/app/ws"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	repo := storage.NewMemory()
	hub := ws.NewHub()
	handlers := apphttp.NewHandlers(repo)
	handlers.OnMessageCreated = hub.BroadcastMessage
	router := apphttp.Router(handlers)
	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.HandleFunc("GET "+cfg.WSPath, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWsHandler(hub, w, r)
	})
	addr := ":" + cfg.HTTPPort
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
