#! /usr/bin/bash


cp /mount/source-code/user_code /mount/artifacts/source_code.py
python -O -m py_compile /mount/artifacts/source_code.py
