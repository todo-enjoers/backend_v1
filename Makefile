.PHONY: help
help:
	@echo "Welcome to helper of Makefile!"
	@echo "Use 'make <target>' where <target> is one of:"
	@echo
	@echo "  all    		run: build -> run-prepare -> run"
	@echo "  build			building a binary file of project"
	@echo "  run-prepare		run: key-generation -> database-up"
	@echo "  key-generation 	creating a couples of keys in secret directory"
	@echo "  database-up		rise up a database with docker files"
	@echo "  run			run"
	@echo ""
	@echo "You should run <all> to fully build and run the project"

.PHONY: all
all: build run-prepare run

.PHONY: build
build:
	@go build --o server.o ./cmd/server/

.PHONY: run-prepare
run-prepare: key-generation database-up

key-generation: creating-dir gen-pub-key gen-pri-key
creating-dir:
	@mkdir -p certs
gen-pub-key:
	@openssl genrsa -out certs/private.pem 2048
gen-pri-key:
	@openssl rsa -in certs/private.pem -outform PEM -pubout -out certs/public.pem

database-up: docker-run docker-compose
docker-run:
	@service docker run
docker-compose:
	2docker compose -f infra/postgres.yaml up -d

.PHONY: run
run:
	@./server.o
