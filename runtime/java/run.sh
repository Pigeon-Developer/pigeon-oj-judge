#! /usr/bin/bash

cd /app/build_result
cat /app/data.in | java Main > /app/data.out
