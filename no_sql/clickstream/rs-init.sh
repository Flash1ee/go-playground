#!/bin/bash

mongo <<EOF
var config = {
    "_id": "rs0",
    "version": 1,
    "members": [
        {
            "_id": 0,
            "host": "mongodb:27017",
            "priority": 3
        },
        {
            "_id": 1,
            "host": "mongodb2:27018",
            "priority": 1
        }
    ]
};
rs.initiate(config, { force: true });
rs.status();
EOF