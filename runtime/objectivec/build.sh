#! /usr/bin/bash

. /usr/share/GNUstep/Makefiles/GNUstep.sh
cp /mount/source-code/user_code /tmp/source_code.m
clang-19 `gnustep-config --objc-flags` /tmp/source_code.m -o /mount/artifacts/main.bin -lobjc -lgnustep-base -I /usr/include/GNUstep/ -I/usr/lib/gcc/x86_64-linux-gnu/12/include -L /usr/lib/GNUstep/Libraries/ -fconstant-string-class=NSConstantString
