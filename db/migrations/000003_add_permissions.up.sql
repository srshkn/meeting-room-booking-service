CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    code text NOT NULL UNIQUE
);

CREATE TABLE role_permissions (
    role_id SMALLINT NOT NULL,
    permission_id uuid NOT NULL,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_permissions
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id) 
        ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

INSERT INTO permissions (code) VALUES
    ('room:create'),
    ('room:list'),
    ('schedule:create'),
    ('slot:list'),
    ('booking:create'),
    ('booking:list'),
    ('booking:my'),
    ('booking:cancel');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'room:create',
    'room:list',
    'schedule:create',
    'slot:list',
    'booking:list'
)
WHERE r.name = 'admin';

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'room:list',
    'slot:list',
    'booking:create',
    'booking:my',
    'booking:cancel'
)
WHERE r.name = 'user';