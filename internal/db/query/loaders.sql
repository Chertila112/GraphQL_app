-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id = ANY($1::uuid[]);

-- name: GetExercisesByIDs :many
SELECT * FROM exercises WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL;

-- name: GetScheduleDaysByScheduleIDs :many
SELECT * FROM schedule_days
WHERE schedule_id = ANY($1::uuid[])
ORDER BY schedule_id, day_of_week;

-- name: GetScheduleExercisesByDayIDs :many
SELECT se.*
FROM schedule_exercises se
WHERE se.schedule_day_id = ANY($1::uuid[])
ORDER BY se.schedule_day_id, se.position;

-- name: GetSchedulesByUserIDs :many
SELECT DISTINCT ON (user_id) *
FROM schedules
WHERE user_id = ANY($1::uuid[]) AND is_active = TRUE;

-- name: GetAllSchedulesByUserIDs :many
SELECT * FROM schedules WHERE user_id = ANY($1::uuid[]);
