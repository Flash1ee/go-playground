rs.initiate({
    _id: "rs-shard-01",
    version: 1,
    members: [
        {_id: 0, host: "mongo-shard-01:27017"},
        {_id: 1, host: "mongo-shard-02:27017"},
        {_id: 2, host: "mongo-shard-03:27017"}
    ]
})
