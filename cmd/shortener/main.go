package main

import (
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/services"
	"github.com/patraden/ya-practicum-go-shortly/internal/web"
)

func main() {
	ls := services.NewSimpleLinkStore()
	hl := web.NewLSHandlers(ls)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.HandleRequests)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
