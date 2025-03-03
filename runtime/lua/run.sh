#! /usr/bin/bash

cat /app/data.in | lua /mount/artifacts/source_code.luac > /app/data.out
