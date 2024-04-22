CREATE TABLE IF NOT EXISTS allowances (
    id SERIAL PRIMARY KEY,
    type VARCHAR(25) NOT NULL UNIQUE,
    max_amount FLOAT NOT NULL
);

INSERT INTO allowances (type, max_amount) VALUES
('personal', 60000.0),
('donation', 100000.0),
('k-receipt', 50000.0);