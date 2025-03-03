#! /usr/bin/bash


cp /mount/source-code/user_code /mount/artifacts/source_code.bas
cd /mount/artifacts
fbc -lang qb source_code.bas -x build_result/main.bin
