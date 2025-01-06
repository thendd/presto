CREATE TABLE IF NOT EXISTS guilds (
    id TEXT NOT NULL PRIMARY KEY,
    max_warnings_per_user INTEGER DEFAULT 3
);

CREATE TABLE IF NOT EXISTS guild_members (
    user_id TEXT NOT NULL PRIMARY KEY,
    guild_id TEXT NOT NULL,
    warnings INTEGER
);
