package main

import (
	log "log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloWorld)
	log.Info("Server open on the http://localhost:8000")
	http.ListenAndServe(":8000", mux)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	BaseLayout("first page with temple", HeaderPage("Hello World")).Render(r.Context(), w)
}
