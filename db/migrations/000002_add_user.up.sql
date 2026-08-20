INSERT INTO users (
    id,
    username,
    email,
    password_hash,
    role_id
)
VALUES (
    '01a02084-0c80-7a11-8a11-111111111111',
    'dummy_user',
    'dummy-user@example.invalid',
    'dummy-disabled',
    (SELECT id FROM roles WHERE name = 'user')
);

INSERT INTO users (
    id,
    username,
    email,
    password_hash,
    role_id
)
VALUES (
    '01a02084-0c80-7a22-8a22-222222222222',
    'dummy_admin',
    'dummy-admin@example.invalid',
    'dummy-disabled',
    (SELECT id FROM roles WHERE name = 'admin')
);
