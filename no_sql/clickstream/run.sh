#!/bin/bash

docker-compose --env-file ./config/credintials.env up -d

sleep 5

docker exec mongodb /scripts/rs-init.sh