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

-- name: CreateGuildMember :one
INSERT INTO guild_members (guild_id, user_id) VALUES ($1, $2) RETURNING warnings;

-- name: UpdateGuildMemberWarnings :exec
UPDATE guild_members SET warnings = $1 WHERE guild_id = $2 AND user_id = $3;

-- name: GetWarningsFromGuildMember :one
SELECT warnings FROM guild_members WHERE guild_id = $1 AND user_id = $2 LIMIT 1;
