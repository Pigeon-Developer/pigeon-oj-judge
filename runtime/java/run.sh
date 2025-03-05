#! /usr/bin/bash

cd /mount/artifacts
cat /app/data.in | java Main > /app/data.out
