#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.c
gcc /tmp/source_code.c -std=c17 -lm -DONLINE_JUDG -fmax-errors=10 -o /mount/artifacts/main.bin
