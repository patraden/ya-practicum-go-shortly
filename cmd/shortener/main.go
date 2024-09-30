package main

import (
	"log"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/handlers"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func main() {
	repo := repository.NewBasicLinkRepository()
	r := handlers.NewRouter(repo)
	log.Fatal(http.ListenAndServe(":8080", r))
}
