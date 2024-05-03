package main

import (
	"github.com/todo-enjoers/backend_v1/config"
	controller "github.com/todo-enjoers/backend_v1/internal/controller/http"
	"github.com/todo-enjoers/backend_v1/internal/storage"
)

func main() {
	store := storage.NewStorage()
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	app := controller.New(store, cfg)
	if err = app.Run(); err != nil {
		panic(err)
	}
}
