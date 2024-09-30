package main

import (
	"log"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/handlers"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func main() {
	repo := repository.NewBasicLinkRepository()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleLinkRepo(repo))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
