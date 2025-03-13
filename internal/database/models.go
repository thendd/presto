package database

type Guild struct {
	ID                                                  string `gorm:"primaryKey"`
	MaxWarningsPerUser                                  int    `gorm:"default:3"`
	OnReachMaxWarningsPerUser                           int    `gorm:"default:1"`
	SecondsToDeleteMessagesForOnReachMaxWarningsPerUser int    `gorm:"default:0"`
	RoleToGiveOnReachMaxWarningsPerUser                 string
	SecondsPunishedUserShouldKeepRoleFor                int `gorm:"default:0"`
	Members                                             []GuildMember
}

type GuildMember struct {
	UserId   string `gorm:"primaryKey"`
	GuildId  string
	Warnings int `gorm:"default:0"`
}
