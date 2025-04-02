CREATE SCHEMA IF NOT EXISTS %s;

CREATE TABLE IF NOT EXISTS %s.users (
    id TEXT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    age INTEGER CHECK,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'unknown')),
    nationalize VARCHAR(255) NOT NULL
);