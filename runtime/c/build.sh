#! /usr/bin/bash

mkdir /app/build_result
gcc /app/source_code -std=c17 -lm -DONLINE_JUDG -fmax-errors=10 -o /app/build_result/main.bin
