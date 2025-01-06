package api

import (
	"errors"
	"fmt"
	"net/http"
)

func BanUser(guildId string, userId string) error {
	_, statusCode := MakeRequest("/guilds/"+guildId+"/bans/"+userId, http.MethodPut, nil)
	if statusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Could not ban user %s from guild %s", userId, guildId))
	}

	return nil
}
