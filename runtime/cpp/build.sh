#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.cpp
g++ /tmp/source_code.cpp -std=c++14 -lm -DONLINE_JUDG -fmax-errors=10 -o /mount/artifacts/main.bin
