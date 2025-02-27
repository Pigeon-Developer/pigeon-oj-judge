#! /usr/bin/bash

mkdir /app/build_result
cp /app/source_code /app/build_result/source_code.pl
perl -c /app/build_result/source_code.pl
