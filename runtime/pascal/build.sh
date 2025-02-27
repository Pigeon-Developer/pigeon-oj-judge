#! /usr/bin/bash

mkdir -p /app/build_result
fpc /app/source_code -Cs32000000 -Sh -O2 -Co -Ct -Ci -o/app/build_result/main.bin

