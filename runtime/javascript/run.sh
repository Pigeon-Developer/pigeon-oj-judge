#! /usr/bin/bash

cat /app/data.in | node --experimental-transform-types /mount/artifacts/source_code.ts > /app/data.out
