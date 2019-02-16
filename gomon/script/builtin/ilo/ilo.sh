#!/usr/bin/env bash

 http_request=$(cat Get_EmHealth.xml | sed -e "s/~username~/$2/" -e "s/~password~/$3/")

 http_response=$(curl -X POST -s -w "%{http_code}" -d "$http_request" -k  "https://$1/ribcl")
 http_code=$(echo "$http_response" | tail -n 1)
 http_body=$(echo "$http_response" | head -n -1)

 echo "CODE: $http_code"
 echo "BODY: $http_body"

 # For new iLO version
 # curl -k  "https://$1/redfish/v1/systems/1" -i -L -u "$2:$3"
