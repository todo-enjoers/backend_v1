DIR_NAME := "certs"
DOCKER_IMAGE_TAG=todoer
CONFIG_PATH=./config.toml

export CONFIG_PATH

.PHONY: help
help:
	@echo "Welcome to helper of Makefile!"
	@echo "Use 'make <target>' where <target> is one of:"
	@echo
	@echo "  all    		run: build -> run-prepare -> run"
	@echo "  build			building a binary file of project"
	@echo "  run-prepare		run: key-generation"
	@echo "  key-generation 	creating a couples of keys in secret directory"
#	@echo "  database-up		rise up a database with docker files"
	@echo "  run			run"
	@echo ""
	@echo "To start a db connection use:"
	@echo "	1. service docker run"
	@echo "	2. sudo docker compose -f infra/postgres.yaml up -d"
	@echo "You should run <all> to fully build and run the project"

.PHONY: start
start: build key-generation run

.PHONY: build
build:
	go build --o server.o ./cmd/server/

.PHONY: run
run:
	@./server.o

.PHONY: key-generation
key-generation:
	@if [ -d "$(DIR_NAME)" ]; then \
  		echo "Директория '$(DIR_NAME)' существует."; \
  	else \
    	mkdir -p certs &&\
    	openssl genrsa -out certs/private.pem 2048 &&\
    	openssl rsa -in certs/private.pem -outform PEM -pubout -out certs/public.pem; \
    fi


.PHONY: dock/build
dock/build:
	docker build \
		--file=infra/Dockerfile \
        --tag=$(DOCKER_IMAGE_TAG) .

.PHONY: dock/run
dock/run: dock/build key-generation
	cd infra && docker-compose up --build -d

.PHONY: diogram
diogram:
	go tool goplantuml -recursive ./ > ./diogram.puml