-- name: GetExercise :one
SELECT * FROM exercises WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListExercises :many
SELECT * FROM exercises
WHERE deleted_at IS NULL
  AND (sqlc.narg('muscle_group')::muscle_group IS NULL OR muscle_group = sqlc.narg('muscle_group'))
  AND (sqlc.narg('equipment')::equipment IS NULL OR equipment = sqlc.narg('equipment'))
ORDER BY name;

-- name: CreateExercise :one
INSERT INTO exercises (id, name, description, muscle_group, equipment)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateExercise :one
UPDATE exercises
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    muscle_group = COALESCE(sqlc.narg('muscle_group'), muscle_group),
    equipment = COALESCE(sqlc.narg('equipment'), equipment),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteExercise :exec
UPDATE exercises SET deleted_at = NOW() WHERE id = $1;