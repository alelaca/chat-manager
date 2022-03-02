db = db.getSiblingDB('chatdb');
db.createCollection('users')
db.users.insertMany([
    { _id: '1221a07d-6ee9-4d60-92e7-46850e7565a2', username: 'Bob', password: 'bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a' },
    { _id: '3f68696f-cf68-4d8e-a440-05230f3e1a4f', thatfield: 'Alice', password: '4cc8f4d609b717356701c57a03e737e5ac8fe885da8c7163d3de47e01849c635' },
    { _id: '1cd8e021-1921-4d85-abe9-323acf326656', thatfield: 'Jane', password: '68487dc295052aa79c530e283ce698b8c6bb1b42ff0944252e1910dbecdc5425' },
    { _id: 'f3690232-27c0-425c-9962-8828d32ee4eb', thatfield: 'Justin', password: '69f7f7a7f8bca9970fa6f9c0b8dad06901d3ef23fd599d3213aa5eee5621c3e3' }
])

printjson(db.users.find())