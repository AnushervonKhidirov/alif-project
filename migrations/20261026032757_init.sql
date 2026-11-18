-- +goose Up
CREATE TYPE "role" AS ENUM ('ADMIN', 'USER');

CREATE TYPE "gender" AS ENUM ('MALE', 'FEMALE');

CREATE TABLE
    IF NOT EXISTS pet_types (
        id SERIAL PRIMARY KEY,
        value VARCHAR(50) UNIQUE NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS breeds (
        id SERIAL PRIMARY KEY,
        value VARCHAR(50) NOT NULL,
        pet_type_id INT NOT NULL REFERENCES pet_types (id) ON DELETE CASCADE,
        UNIQUE (pet_type_id, value),
        UNIQUE (id, pet_type_id)
    );

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
        microchip_id VARCHAR(255),
        pet_type_id INT NOT NULL REFERENCES pet_types (id),
        breed_id INT REFERENCES breeds (id),
        CONSTRAINT fk_pet_breed_pet_type
            FOREIGN KEY (breed_id, pet_type_id)
            REFERENCES breeds (id, pet_type_id),
        lost bool DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_user_email ON users (email);

-- +goose Down
DROP TABLE IF EXISTS pets;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS breeds;

DROP TABLE IF EXISTS pet_types;

DROP TYPE IF EXISTS "role";

DROP TYPE IF EXISTS "gender";