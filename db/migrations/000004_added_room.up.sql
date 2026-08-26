CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    description TEXT,
    capacity INT CHECK (capacity > 0), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX rooms_name_unique_idx
ON rooms (name);