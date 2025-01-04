package cdn

import "presto/internal/discord/config"

func GetGuildIconURL(guildId string, icon string) string {
	return config.CDN_BASE_URL + "/icons/" + guildId + "/" + icon + ".png"
}
