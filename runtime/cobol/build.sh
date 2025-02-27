#! /usr/bin/bash

mkdir /app/build_result
cp /app/source_code /app/source_code.cobol
cobc -free -x -o /app/build_result/main.bin source_code.cobol
