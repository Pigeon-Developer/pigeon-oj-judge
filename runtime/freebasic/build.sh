#! /usr/bin/bash


cp /mount/source-code/user_code /tmp/source_code.bas
cd /tmp
fbc -lang qb source_code.bas -x /mount/artifacts/main.bin
