#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.oc
clang-19 /tmp/source_code.oc -ObjC -o /mount/artifacts/main.bin
