#! /usr/bin/bash


cd /app
cp /mount/source-code/user_code /tmp/main.go
go build -o /mount/artifacts/main.bin /tmp/main.go
