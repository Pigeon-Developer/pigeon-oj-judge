#! /usr/bin/bash

mkdir /app/build_result
clang-19 /app/source_code -std=c17 -o /app/build_result/main.bin
