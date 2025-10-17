CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL,
                                     balance NUMERIC(12,2) DEFAULT 0
    );

INSERT INTO users (name, email, balance) VALUES
                                             ('Alice', 'alice@example.com', 1000),
                                             ('Bob', 'bob@example.com', 500),
                                             ('Charlie', 'charlie@example.com', 300);
