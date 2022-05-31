db.createUser(
    {
        user: "root",
        pwd: "root",
        roles: [
            {
                role: "readWrite",
                db: "news"
            }
        ]
    }
);

db.createCollection("products");
