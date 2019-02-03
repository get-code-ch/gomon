#!/usr/bin/env bash
R=$(( RANDOM % 2 ))
if [[ "$R" == "1" ]]; then
    echo "OK: Value $R"
else
    echo "CRITICAL: Value $R"
    exit "1"
fi