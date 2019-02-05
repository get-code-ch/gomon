#!/usr/bin/env bash
R=$(( RANDOM % 3 ))
if [[ "$R" == "1" ]]; then
    echo "OK: Value $R"
elif [[ "$R" == "2" ]]; then
    echo "WARNING: Value $R"
else
    echo "CRITICAL: Value $R"
    exit "1"
fi