#!/bin/sh

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir db/migration -seq init_schema

migrate -path ./internal/db/migrations -database 'postgres://postgres:postgres@db:5432/postgres?sslmode=disable' -verbose up

./api
