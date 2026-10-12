-- name: GetLocations :many
SELECT * FROM facility.locations;

-- name: RegistLocation :one
INSERT INTO facility.locations (name,address,timezone,open_time,close_time)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;
