CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    slot_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    conference_link TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT bookings_status_check
        CHECK (status IN ('active', 'cancelled')),

    CONSTRAINT fk_bookings_slot
        FOREIGN KEY (slot_id)
        REFERENCES slots(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_bookings_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_bookings_active_slot
    ON bookings(slot_id)
    WHERE status = 'active';

CREATE INDEX idx_bookings_user
    ON bookings(user_id);