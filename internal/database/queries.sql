-- name: CreateGuild :one
INSERT INTO guilds (id) VALUES ($1) RETURNING *;

-- name: GetGuild :one
SELECT
    *
FROM
    guilds
WHERE
    id = $1
LIMIT 1;

-- name: CreateWarnedUser :one
INSERT INTO warned_users (guild_id, user_id) VALUES ($1, $2) RETURNING warnings;

-- name: UpdateWarnedUserWarnings :exec
UPDATE warned_users SET warnings = $1 WHERE guild_id = $2 AND user_id = $3;

-- name: GetWarningsFromWarnedUser :one
SELECT warnings FROM warned_users WHERE guild_id = $1 AND user_id = $2 LIMIT 1;
