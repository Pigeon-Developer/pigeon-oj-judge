#! /usr/bin/bash

mkdir /app/build_result
cp /app/source_code /app/build_result/source_code.bas
cd /app/build_result
fbc -lang qb source_code.bas -x build_result/main.bin
