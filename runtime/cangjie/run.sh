#! /usr/bin/bash

source /app/cangjie/envsetup.sh
cat /app/data.in | /app/build_result/main.bin > /app/data.out
