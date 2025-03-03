#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.pas
fpc /tmp/source_code.pas -Cs32000000 -Sh -O2 -Co -Ct -Ci -o/mount/artifacts/main.bin

