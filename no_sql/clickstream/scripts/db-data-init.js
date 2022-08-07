use clickstream

db.events.insertOne({myfield: 'hello', thatfield: 'testing'})

// Выдача информации о пользователе по userID.
db.users.insertOne({
    "userID": 1,
    name: "Vladimir",
    "nickname": "flashie",
    "phone": "+79539005439",
    "email": "flashie@mail.com"
})
db.users.insertOne({
    "userID": 2,
    name: "Alexey",
    "nickname": "habrouser",
    "phone": "+79539005439",
    "email": "flashie@mail.com"
})
db.users.insertOne({
    "userID": 3,
    name: "Maria",
    "nickname": "cheburashka",
    "phone": "+79539005439",
    "email": "flashie@mail.com"
})
db.users.insertOne({
    "userID": 4,
    name: "Kate",
    "nickname": "stonny",
    "phone": "+79539005439",
    "email": "flashie@mail.com"
})

db.createUser(
    {
        user: "flashie",
        pwd: "project",
        roles: [
            {
                role: "dbAdmin",
                db: "clickstream"
            }
        ]
    }
)