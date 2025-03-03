#! /usr/bin/bash


source /app/cangjie/envsetup.sh
cp /mount/source-code/user_code /tmp/source_code.cj
cjc --diagnostic-format noColor -o /mount/artifacts/main.bin /tmp/source_code.cj
