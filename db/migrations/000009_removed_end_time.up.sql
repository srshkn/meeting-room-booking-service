ALTER TABLE slots
    DROP CONSTRAINT IF EXISTS slots_duration_check;

ALTER TABLE slots
    DROP COLUMN IF EXISTS end_time;