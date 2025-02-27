#! /usr/bin/bash

cat /app/data.in | Rscript /app/build_result/source_code.R > /app/data.out
