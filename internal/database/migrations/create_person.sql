CREATE SCHEMA IF NOT EXISTS %s;

CREATE TABLE IF NOT EXISTS %s.person (
    id TEXT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    age INTEGER,
    gender VARCHAR(10),
    nationality VARCHAR(255)
);