#!/bin/sh
# export env variables from .env file
export $(grep -v '^#' .env | xargs)

echo "----- Seeding data..."
go run script.go data_seeding local
echo "----- ...Data seeded"

air