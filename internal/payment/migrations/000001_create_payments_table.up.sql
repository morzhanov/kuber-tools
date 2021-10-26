CREATE TABLE payments (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount int NOT NULL,
    status VARCHAR(255) NOT NULL
);
