CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    room_id UUID NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_slots_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT slots_duration_check
        CHECK (end_time = start_time + INTERVAL '30 minutes'),

    CONSTRAINT slots_room_start_unique
        UNIQUE (room_id, start_time)
);