#! /usr/bin/bash

mkdir /app/build_result
cp /app/source_code /app/source_code.f95
f95 -o /app/build_result/main.bin /app/source_code.f95
