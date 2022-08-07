#!/bin/bash

mongo <<EOF
// Enable sharding for database `clickstream`
sh.enableSharding("clickstream")

// Setup shardingKey for collection `MyCollection`**
db.adminCommand( { shardCollection: "clickstream.events", key: { created_at:  "hashed" } } )
EOF