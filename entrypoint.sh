#!/bin/bash

source .env

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir ./schema -seq init_schema

migrate -path ./schema -database "$DB_SOURCE" -verbose up

./main
