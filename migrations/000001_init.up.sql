CREATE TABLE IF NOT EXISTS users
(
    user_id    SERIAL PRIMARY KEY,
    username   varchar(255) NOT NULL unique,
    password   varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS coins
(
    wallet_id SERIAL PRIMARY KEY,
    user_id   INT REFERENCES users (user_id) NOT NULL,
    amount    INT                            NOT NULL CHECK (amount >= 0)
);

CREATE TABLE IF NOT EXISTS merch
(
    merch_id SERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    price    INT          NOT NULL
);

CREATE TABLE IF NOT EXISTS purchases
(
    purchase_id   SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users (user_id)  NOT NULL,
    merch_id      INT REFERENCES merch (merch_id) NOT NULL,
    purchase_date TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS operations
(
    transaction_id   SERIAL PRIMARY KEY,
    from_user_id     INT REFERENCES users (user_id),
    to_user_id       INT REFERENCES users (user_id),
    amount           INT NOT NULL CHECK (amount > 0),
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO merch (name, price)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);