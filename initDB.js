// Use GOIODI DB
use goiodi;

// Remove everything (all collections documents) in the database
db.dropDatabase();

// Make some fields unique
db.users.ensureIndex({username:1},{unique:true});
db.users.ensureIndex({email:1},{unique:true});

db.words.ensureIndex({word:1},{unique:true});
