#!/usr/bin/env bash
if [ $# -eq 2 ]
    then
        if [ "$2" = "" ]
            then
                echo "WARNING: Hello $1 no password"
                exit 0
            else
                echo "OK: Hello $1 your password is $2"
                exit 0
            fi
    else
        echo "CRITICAL: Error wrong number of arguments"
        exit 1
fi
