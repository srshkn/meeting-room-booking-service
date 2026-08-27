ALTER TABLE slots
    ADD COLUMN  end_time TIMESTAMPTZ NOT NULL;

ALTER TABLE slots
    ADD CONSTRAINT slots_duration_check
        CHECK (end_time = start_time + INTERVAL '30 minutes');