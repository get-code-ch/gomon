#!/usr/bin/env bash

# docker exec -it gomon_mongo_1 mongodump  --username gomon --password 60M0n --authenticationDatabase admin --collection command --db gomondb
# docker exec -it gomon_mongo_1 mongorestore  --username gomon --password 60M0n --authenticationDatabase admin --collection command --db gomondb ./dump/gomondb/command.bson

mongo localhost:27017/gomondb --username gomon --password --authenticationDatabase admin  gomondb.js
