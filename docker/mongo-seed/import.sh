#!/bin/bash

mongoimport \
    --db budget-tracker-v2 \
    --host "mongodb:27017" \
    --username root \
    --password example \
    --collection users \
    --type json \
    --file /mongo-seed/init.json \
    --authenticationDatabase=admin \
    --jsonArray
