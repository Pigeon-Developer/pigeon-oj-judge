#! /usr/bin/bash

mkdir /app/build_result
source /app/cangjie/envsetup.sh
cjc --diagnostic-format noColor -o /app/build_result/main.bin /app/source_code
