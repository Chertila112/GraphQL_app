-- name: GetSchedule :one
SELECT * FROM schedules WHERE id = $1 LIMIT 1;

-- name: GetActiveScheduleByUser :one
SELECT * FROM schedules WHERE user_id = $1 AND is_active = TRUE LIMIT 1;

-- name: ListSchedules :many
SELECT * FROM schedules ORDER BY created_at DESC;

-- name: CreateSchedule :one
INSERT INTO schedules (id, user_id, name, is_active)
VALUES ($1, $2, $3, FALSE)
RETURNING *;

-- name: UpdateSchedule :one
UPDATE schedules
SET name = COALESCE(sqlc.narg('name'), name),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateOtherSchedules :exec
UPDATE schedules SET is_active = FALSE WHERE user_id = $1 AND id != $2;
--
-- name: GetScheduleDay :one
SELECT * FROM schedule_days WHERE id = $1 LIMIT 1;

-- name: DeleteSchedule :exec
DELETE FROM schedules WHERE id = $1;

-- name: GetScheduleDaysBySchedule :many
SELECT * FROM schedule_days WHERE schedule_id = $1 ORDER BY day_of_week;

-- name: CreateScheduleDay :one
INSERT INTO schedule_days (id, schedule_id, day_of_week, is_rest_day, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateScheduleDay :one
UPDATE schedule_days
SET day_of_week = COALESCE(sqlc.narg('day_of_week'), day_of_week),
    is_rest_day = COALESCE(sqlc.narg('is_rest_day'), is_rest_day),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1
RETURNING *;

-- name: AddExerciseToDay :one
INSERT INTO schedule_exercises (id, schedule_day_id, exercise_id, position, sets, reps, weight, rest, duration, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetScheduleExercisesByDay :many
SELECT se.*, e.name as exercise_name, e.muscle_group, e.equipment
FROM schedule_exercises se
JOIN exercises e ON e.id = se.exercise_id
WHERE se.schedule_day_id = $1
ORDER BY se.position;

-- name: UpdateScheduleExercise :one
UPDATE schedule_exercises
SET position = COALESCE(sqlc.narg('position'), position),
    sets = COALESCE(sqlc.narg('sets'), sets),
    reps = COALESCE(sqlc.narg('reps'), reps),
    weight = COALESCE(sqlc.narg('weight'), weight),
    rest = COALESCE(sqlc.narg('rest'), rest),
    duration = COALESCE(sqlc.narg('duration'), duration),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1
RETURNING *;

-- name: RemoveExerciseFromDay :exec
DELETE FROM schedule_exercises WHERE id = $1;
