-- +goose Up
CREATE TABLE
    IF NOT EXISTS vet_clinics (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        about TEXT,
        address VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE vet_clinics;