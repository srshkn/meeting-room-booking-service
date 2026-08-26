-- name: CreateRoom :one
INSERT INTO rooms (
    name,
    description,
    capacity
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, name, description, capacity, created_at;