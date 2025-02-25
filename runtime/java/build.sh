#! /usr/bin/bash

mkdir /app/build_result
cd /app
cp source_code Main.java
javac Main.java -d /app/build_result
