#! /usr/bin/bash

export DOTNET_EnableWriteXorExecute=0
rm /app/Program.cs
cp /mount/source-code/user_code /app/Program.cs
cd /app
dotnet build --property:Configuration=Release -o /mount/artifacts
