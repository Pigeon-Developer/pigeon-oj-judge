#! /usr/bin/bash


cp /mount/source-code/user_code /tmp/source_code.cobol
cobc -free -x -o /mount/artifacts/main.bin /tmp/source_code.cobol
