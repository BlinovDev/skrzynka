package main

import (
	_ "embed"
	"log"
	"net/http"

	"anon-skrzynka/api"
	"anon-skrzynka/app/config"
	apphttp "anon-skrzynka/app/http"
	"anon-skrzynka/app/storage"
	"anon-skrzynka/app/ws"
)

//go:embed swagger.html
var swaggerHTML []byte

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
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(api.OpenAPIYAML)
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(swaggerHTML)
	})
	mux.HandleFunc("GET "+cfg.WSPath, func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWsHandler(hub, w, r)
	})
	addr := ":" + cfg.HTTPPort
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, apphttp.CORS(mux, cfg.AllowedOrigins)); err != nil {
		log.Fatal(err)
	}
}
