#!/bin/sh

echo "----- Seeding data..."
go run script.go data_seeding local
echo "----- ...Data seeded"

air