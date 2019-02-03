#!/usr/bin/env bash
if [ $# -eq 1 ]
    then
        echo "OK: Hello $1"
        exit 0
    else
        echo "CRITICAL: Error wrong number of arguments"
        exit 1
fi
