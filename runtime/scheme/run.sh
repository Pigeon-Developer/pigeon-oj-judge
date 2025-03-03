#! /usr/bin/bash

cat /app/data.in | guile /mount/artifacts/source_code.scm > /app/data.out
