#! /usr/bin/bash

mkdir /app/build_result
rm /app/Program.cs
cp /app/source_code /app/Program.cs
cd /app
dotnet build --property:Configuration=Release -o build_result
