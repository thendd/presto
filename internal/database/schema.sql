CREATE TABLE IF NOT EXISTS guilds (
    id TEXT NOT NULL PRIMARY KEY,
    max_warnings_per_user INTEGER DEFAULT 3,
    on_reach_max_warnings_per_user INTEGER DEFAULT 1,
    seconds_to_delete_messages_for_on_reach_max_warnings_per_user INTEGER DEFAULT 0,
    role_to_give_on_reach_max_warnings_per_user TEXT,
    seconds_warned_user_should_keep_role_for INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS guild_members (
    user_id TEXT NOT NULL PRIMARY KEY,
    guild_id TEXT NOT NULL,
    warnings INTEGER DEFAULT 0
);
