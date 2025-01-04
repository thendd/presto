package api

import (
	"encoding/json"
	"net/http"

	"presto/internal/discord"
)

func GetGuildById(guildId string) (guild discord.Guild) {
	response, statusCode := MakeRequest("/guilds/"+guildId, http.MethodGet, nil)
	if statusCode != 200 {
		return discord.Guild{}
	}

	json.Unmarshal(response, &guild)
	return guild
}
