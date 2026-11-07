-include .env

DB_USER ?= postgres
DB_PORT ?= 5432
DB_NAME ?= pet_db

DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable
GOOSE ?= goose

.PHONY: up down check-db-password migrate-up migrate-down

up:
	docker compose up

up-b:
	docker compose up --build

up-d:
	docker compose up -d

down:
	docker compose down

down-v:
	docker compose down -v

migrate-up: check-db-password
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

migrate-down: check-db-password
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down
