#! /usr/bin/bash


cp /mount/source-code/user_code /tmp/source_code.f95
f95 -o /mount/artifacts/main.bin /tmp/source_code.f95
