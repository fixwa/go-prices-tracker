db.createUser(
    {
        user: "root",
        pwd: "root",
        roles: [
            {
                role: "readWrite",
                db: "products"
            }
        ]
    }
);

db.createCollection("products");
