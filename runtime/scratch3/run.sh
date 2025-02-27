#! /usr/bin/bash

cat /app/data.in | scratch-run /app/build_result/source_code.sb3 > /app/data.out
