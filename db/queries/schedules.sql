-- name: CreateSchedule :one
WITH created_schedule AS (
    INSERT INTO schedules (
        room_id,
        start_time,
        end_time
    )
    VALUES (
        sqlc.arg(room_id),
        sqlc.arg(start_time),
        sqlc.arg(end_time)
    )
    RETURNING *
),
created_days AS (
    INSERT INTO schedule_days (
        schedule_id,
        day_of_week
    )
    SELECT
        created_schedule.id,
        unnest(sqlc.arg(days_of_week)::smallint[])
    FROM created_schedule
    RETURNING day_of_week
)
SELECT
    created_schedule.id,
    created_schedule.room_id,
    created_schedule.start_time,
    created_schedule.end_time,
    created_schedule.created_at,
    array_agg(created_days.day_of_week ORDER BY created_days.day_of_week)::smallint[] AS days_of_week
FROM created_schedule
CROSS JOIN created_days
GROUP BY
    created_schedule.id,
    created_schedule.room_id,
    created_schedule.start_time,
    created_schedule.end_time,
    created_schedule.created_at;