#! /usr/bin/bash

mkdir /app/build_result
cd /app
cp source_code main.go
go build -o /app/build_result/main.bin main.go
