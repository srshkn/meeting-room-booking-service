CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    room_id UUID NOT NULL UNIQUE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_schedules_rooms
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,
    
    CONSTRAINT schedules_time_check
        CHECK (start_time < end_time)
);