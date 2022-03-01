db = db.getSiblingDB('chatdb');
db.createCollection('users')
db.users.insertMany([
    { _id: '1221a07d-6ee9-4d60-92e7-46850e7565a2', username: 'Bob', password: '111111' },
    { _id: '3f68696f-cf68-4d8e-a440-05230f3e1a4f', thatfield: 'Alice', password: '222222' },
    { _id: '1cd8e021-1921-4d85-abe9-323acf326656', thatfield: 'Jane', password: '333333' },
    { _id: 'f3690232-27c0-425c-9962-8828d32ee4eb', thatfield: 'Justin', password: '444444' }
])

printjson(db.users.find())