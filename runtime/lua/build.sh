#! /usr/bin/bash

cp /mount/source-code/user_code /tmp/source_code.lua
luac -o /mount/artifacts/source_code.luac /tmp/source_code.lua
