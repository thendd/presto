package database

type OnReachMaxWarningsPerUserType int8

const (
	ON_REACH_MAX_WARNINGS_PER_USER_BAN       OnReachMaxWarningsPerUserType = 0
	ON_REACH_MAX_WARNINGS_PER_USER_KICK      OnReachMaxWarningsPerUserType = 1
	ON_REACH_MAX_WARNINGS_PER_USER_GIVE_ROLE OnReachMaxWarningsPerUserType = 2
	ON_REACH_MAX_WARNINGS_PER_USER_NOTHING   OnReachMaxWarningsPerUserType = 3
)
