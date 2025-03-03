#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.c
clang-19 /tmp/source_code.c -std=c17 -o /mount/artifacts/main.bin
