db.clickstream.insertOne({ myfield: 'hello', thatfield: 'testing' })

db.createUser(
    {
        user: "flashie",
        pwd: "project",
        roles: [
            {
                role: "readWrite",
                db: "clickstream"
            }
        ]
    }
)