-- name: CreateSlots :exec
INSERT INTO slots (
    room_id,
    start_time
)
SELECT
    sqlc.arg(room_id),
    unnest(sqlc.arg(start_times)::timestamptz[])
ON CONFLICT (room_id, start_time) DO NOTHING;

-- name: GetUnbookedSlots :many
SELECT
    s.id,
    s.room_id,
    s.start_time,
    s.start_time + INTERVAL '30 minutes' AS end_time
FROM slots AS s
LEFT JOIN bookings AS b
    ON b.slot_id = s.id
    AND b.status = 'active'
WHERE b.id IS NULL
ORDER BY s.start_time;