CREATE TABLE schedule_days (
    schedule_id UUID NOT NULL,
    day_of_week SMALLINT NOT NULL,

    PRIMARY KEY (schedule_id, day_of_week),

    CONSTRAINT fk_schedule_days_schedule
        FOREIGN KEY (schedule_id)
        REFERENCES schedules(id)
        ON DELETE CASCADE,

    CONSTRAINT schedule_days_day_check
        CHECK (day_of_week BETWEEN 1 AND 7)
);