-- name: CreateTaskHistory :one
INSERT INTO task_histories (user_uuid, name, start_time, end_time)
VALUES (@user_uuid, @name, @start_time, @end_time)
RETURNING name,
    CONCAT(
        FLOOR(EXTRACT(EPOCH FROM (end_time - start_time)) / 3600), ' hours ',
        FLOOR(EXTRACT(EPOCH FROM (end_time - start_time)) / 60 % 60), ' minutes'
    ) AS duration;

-- name: GetTasksResultByPeriod :many
WITH task_durations AS (
    SELECT
        th.name AS task_name,
        EXTRACT(EPOCH FROM (th.end_time - th.start_time)) AS duration_seconds
    FROM
        task_histories th
    WHERE
        th.end_time >= NOW() - CAST($1 AS INTERVAL) AND user_uuid = $2
)
SELECT
    td.task_name,
    CONCAT(
        FLOOR(td.duration_seconds / 3600), ' hours ',
        FLOOR((td.duration_seconds / 60) % 60), ' minutes'
    ) AS duration,
    CONCAT(
        FLOOR(SUM(td.duration_seconds) OVER () / 3600), ' hours ',
        FLOOR((SUM(td.duration_seconds) OVER () / 60) % 60), ' minutes'
    ) AS total_duration
FROM
    task_durations td
ORDER BY
    td.duration_seconds desc;
