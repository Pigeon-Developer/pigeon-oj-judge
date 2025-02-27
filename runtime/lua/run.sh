#! /usr/bin/bash

cat /app/data.in | lua /app/build_result/source_code.luac > /app/data.out
