-- name: CreateTask :one
INSERT INTO tasks (user_uuid, name)
VALUES (@user_uuid, @name)
RETURNING *;

-- name: UpdateTaskEndTime :one
UPDATE tasks
SET end_time = NOW()
WHERE user_uuid = @user_uuid
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE user_uuid = @user_uuid;
