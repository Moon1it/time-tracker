-- name: CreateUser :one
INSERT INTO users (passport_number, surname, name, patronymic, address)
VALUES (@passport_number, @surname, @name, @patronymic, @address)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users
WHERE
    (passport_number = sqlc.narg('passport_number') OR sqlc.narg('passport_number') IS NULL)
    AND (surname = sqlc.narg('surname') OR sqlc.narg('surname') IS NULL)
    AND (name = sqlc.narg('name') OR sqlc.narg('name') IS NULL)
    AND (patronymic = sqlc.narg('patronymic') OR sqlc.narg('patronymic') IS NULL)
    AND (address = sqlc.narg('address') OR sqlc.narg('address') IS NULL)
LIMIT @user_limit OFFSET @user_offset;

-- name: GetUserByUUID :one
SELECT * FROM users
WHERE uuid = @user_uuid;

-- name: GetUserByPassportNumber :one
SELECT * FROM users
WHERE passport_number = @passport_number;

-- name: UpdateUserByUUID :one
UPDATE users
SET surname = coalesce(sqlc.narg('surname'), surname),
    name = coalesce(sqlc.narg('name'), name),
    patronymic = coalesce(sqlc.narg('patronymic'), patronymic),
    address = coalesce(sqlc.narg('address'), address),
    passport_number = coalesce(sqlc.narg('passport_number'), passport_number)
WHERE uuid = @user_uuid
RETURNING *;

-- name: DeleteUserByUUID :exec
DELETE FROM users
WHERE uuid = @user_uuid;
