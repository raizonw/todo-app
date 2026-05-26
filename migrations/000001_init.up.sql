CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id SERIAL PRIMARY KEY,
    version INT NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15)
);
