#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.m
clang-19 -framework Foundation /tmp/source_code.m -o /mount/artifacts/main.bin
