#! /usr/bin/bash

cat /app/data.in | octave-cli -W -q -H /app/build_result/source_code.m > /app/data.out
