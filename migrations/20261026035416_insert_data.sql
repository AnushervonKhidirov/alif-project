-- +goose Up
INSERT INTO
    users (email, password_hash, first_name, role)
VALUES
    (
        'ankhid@gmail.com',
        '$2a$10$0W9dT.CTLb6FbiLJfVWQpuFWCVfKoF51sMClP2uL7cx/mCiCWxMb2', -- password: mypass
        'Anushervon',
        'ADMIN'
    ) ON CONFLICT (email) DO NOTHING;

-- +goose Down
DELETE FROM users
WHERE
    email = 'ankhid@gmail.com';