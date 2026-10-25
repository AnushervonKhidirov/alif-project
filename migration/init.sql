CREATE TYPE "role" AS ENUM ('ADMIN', 'USER');

CREATE TYPE "gender" AS ENUM ('Male', 'Female');

CREATE TABLE
    IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        first_name VARCHAR(25) NOT NULL,
        last_name VARCHAR(25),
        role role NOT NULL DEFAULT 'USER',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS pets (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users (id) ON DELETE CASCADE,
        name VARCHAR(30) NOT NULL,
        about TEXT,
        gender gender,
        birthday DATE,
        microchip_id string,
        pet_type_id INT NOT NULL REFERENCES pet_types (id),
        breed_id INT REFERENCES breeds (id),
        lost bool DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS pet_types (
        id SERIAL PRIMARY KEY,
        value VARCHAR(50) UNIQUE NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS breeds (
        id SERIAL PRIMARY KEY,
        value VARCHAR(50) UNIQUE NOT NULL,
        pet_type_id INT REFERENCES pet_types (id) ON DELETE CASCADE
    );

CREATE INDEX idx_user_email ON users (email);