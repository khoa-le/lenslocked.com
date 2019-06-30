CREATE TABLE users(
    id SERIAL PRIMARY key,
    name TEXT,
    email TEXT NOT NULL
);

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount  INT,
    description TEXT
);