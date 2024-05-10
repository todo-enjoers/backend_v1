package main

import (
	"context"
	"fmt"
	"github.com/todo-enjoers/backend_v1/config"
	controller "github.com/todo-enjoers/backend_v1/internal/controller/http"
	"github.com/todo-enjoers/backend_v1/internal/storage/postgres"
	"github.com/todo-enjoers/backend_v1/pkg/client"
	"log"
)

func main() {
	//init config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	//init pool
	pool, err := client.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to connect to database: %v\n"))
	}
	defer pool.Close()

	//init storage
	store, _ := postgres.New(pool.P())

	//init app
	app := controller.New(store, cfg)
	if err = app.Run(); err != nil {
		panic(err)
	}
}
