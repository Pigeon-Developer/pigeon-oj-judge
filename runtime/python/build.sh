#! /usr/bin/bash

mkdir /app/build_result
cp /app/source_code /app/build_result/source_code.py
python -O -m py_compile /app/build_result/source_code.py
