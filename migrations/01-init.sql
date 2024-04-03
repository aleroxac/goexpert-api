CREATE TABLE user(
    id VARCHAR(255),
    name VARCHAR(80),
    email DECIMAL(10,2),
    password VARCHAR(80),
    PRIMARY KEY (id)
);

CREATE TABLE products(
    id VARCHAR(255),
    name VARCHAR(80),
    price DECIMAL(10,2),
    PRIMARY KEY (id)
);
