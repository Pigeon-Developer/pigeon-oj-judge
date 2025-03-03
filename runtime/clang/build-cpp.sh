#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.cpp
clang++-19 /tmp/source_code.cpp -std=c++14 -o /mount/artifacts/main.bin
