CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL PRIMARY KEY,
    "username" TEXT UNIQUE NOT NULL,
    "password_hash" TEXT NOT NULL,
    "coins" INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Merch (
    "type" VARCHAR(128) PRIMARY KEY,
    "cost" INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inventory (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES users("id"),
    "item_name" VARCHAR(128) NOT NULL REFERENCES Merch("type"),
    "amount" INTEGER NOT NULL DEFAULT 1,
    UNIQUE (user_id, item_name)
);

CREATE TABLE IF NOT EXISTS wallet_operations (
    "id" SERIAL PRIMARY KEY,
    "sender_id" INTEGER NOT NULL REFERENCES users("id"),
    "receiver_id" INTEGER NOT NULL REFERENCES users("id"),
    "value" INT DEFAULT 0
);


INSERT INTO Merch("type", "cost") VALUES ('t-shirt', 80);
INSERT INTO Merch("type", "cost") VALUES ('cup', 20);
INSERT INTO Merch("type", "cost") VALUES ('book', 50);
INSERT INTO Merch("type", "cost") VALUES ('pen', 10);
INSERT INTO Merch("type", "cost") VALUES ('powerbank', 200);
INSERT INTO Merch("type", "cost") VALUES ('hoody', 300);
INSERT INTO Merch("type", "cost") VALUES ('umbrella', 200);
INSERT INTO Merch("type", "cost") VALUES ('socks', 10);
INSERT INTO Merch("type", "cost") VALUES ('wallet', 50);
INSERT INTO Merch("type", "cost") VALUES ('pink-hoody', 500);

INSERT INTO users("username", "password_hash", "coins")
VALUES ('rich_user', "hash", 1000);

INSERT INTO users("username", "password_hash", "coins")
VALUES ('poor_user', "hash2", 50);