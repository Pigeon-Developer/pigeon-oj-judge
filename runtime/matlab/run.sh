#! /usr/bin/bash

cat /app/data.in | octave-cli -W -q -H /mount/artifacts/source_code.m > /app/data.out
