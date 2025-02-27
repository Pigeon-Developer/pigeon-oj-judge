#! /usr/bin/bash

cat /app/data.in | node --experimental-transform-types /app/build_result/source_code.ts > /app/data.out
