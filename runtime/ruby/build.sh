#! /usr/bin/bash

mkdir /app/build_result
ruby -c /app/source_code
cp /app/source_code /app/build_result/source_code.rb
