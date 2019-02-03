#!/usr/bin/env bash
P=$(ping -c 2 -w 2 -q $1 | sed -Esn 's/^rtt.*=[\s]*(.*)\sms$/\1/p') > /dev/null
if [[ "$P" != "" ]]; then
    IFS='/' read -r -a Array <<< "$P"
    echo "OK: Successfull Ping $1|min=${Array[0]};avg=${Array[1]};max=${Array[2]};mdev=${Array[3]}"
else
    echo "CRITICAL: Ping error $1"
    exit "1"
fi
