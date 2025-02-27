#! /usr/bin/bash

cat /app/data.in | guile /app/build_result/source_code.scm > /app/data.out
