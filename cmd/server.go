package main

import (
	"log"
	"net/http"
)

func main() {
	store := store.New()
	ctrl := http.New(store)
	if err := ctrl.Start(); err != nil {
		panic(err)
	}
	log.Fatal(http.ListenAndServe(":8080", r))
}
