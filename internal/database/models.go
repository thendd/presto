package database

type Guild struct {
	ID                                                  string `gorm:"primaryKey"`
	MaxWarningsPerUser                                  int8   `gorm:"default:3"`
	OnReachMaxWarningsPerUser                           int8   `gorm:"default:1"`
	SecondsToDeleteMessagesForOnReachMaxWarningsPerUser int    `gorm:"default:0"`
	RoleToGiveOnReachMaxWarningsPerUser                 string
	SecondsPunishedUserShouldKeepRoleFor                int `gorm:"default:0"`
	Members                                             []GuildMember
}

type GuildMember struct {
	UserId   string `gorm:"primaryKey"`
	GuildId  string
	Warnings int8 `gorm:"default:0"`
}
