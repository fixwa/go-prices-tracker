db.createUser(
    {
        user: "root",
        pwd: "root",
        roles: [
            {
                role: "readWrite",
                db: "pricesTracker"
            }
        ]
    }
);

db.createCollection("products");
