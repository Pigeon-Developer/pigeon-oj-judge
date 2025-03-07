#! /usr/bin/bash

export DOTNET_EnableWriteXorExecute=0
cat /app/data.in | /mount/artifacts/app > /app/data.out
