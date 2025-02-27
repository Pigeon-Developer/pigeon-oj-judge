#! /usr/bin/bash

mkdir /app/build_result
g++ /app/source_code -std=c++14 -lm -DONLINE_JUDG -fmax-errors=10 -o /app/build_result/main.bin
