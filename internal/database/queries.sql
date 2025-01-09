-- name: CreateGuild :one
INSERT INTO guilds (id) VALUES ($1) RETURNING *;

-- name: UpdateMaxWarningsPerUserFromGuild :exec
UPDATE guilds SET max_warnings_per_user = $1 WHERE id = $2;

-- name: UpdateOnReachMaxWarningsPerUserFromGuild :exec
UPDATE guilds SET on_reach_max_warnings_per_user = $1 WHERE id = $2;

-- name: UpdateSecondsToDeleteUserMessagesForOnReachMaxWarningsPerUserFromGuild :exec
UPDATE guilds SET seconds_to_delete_messages_for_on_reach_max_warnings_per_user = $1 WHERE id = $2;

-- name: UpdateRoletoGiveOnReachMaxWarningsPerUserFromGuild :exec
UPDATE guilds SET role_to_give_on_reach_max_warnings_per_user = $1 WHERE id = $2;

-- name: UpdateSecondsUserShouldKeepRoleForFromGuild :exec
UPDATE guilds SET seconds_warned_user_should_keep_role_for = $1 WHERE id = $2;

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
