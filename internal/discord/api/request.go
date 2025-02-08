package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"presto/internal/log"

	"presto/internal/discord/config"
)

// TODO: using a client in the way I'm using right now might lead into some problems in the future
// so I'll look further into that later
var client = http.Client{}

// Creates an `http.Request` object and does a request using a global `http.Client`
func MakeRequest(endpoint string, method string, body any) ([]byte, int) {
	var buffer bytes.Buffer

	if body != nil {
		err := json.NewEncoder(&buffer).Encode(body)
		if err != nil {
			log.Error(err.Error())
		}
	}

	// Ignore error just because (documentatio doesn't say anything about the cases in
	// which an error is returned).
	request, _ := http.NewRequest(method, config.API_BASE_URL+endpoint, &buffer)

	request.Header.Add("Authorization", "Bot "+config.DISCORD_BOT_TOKEN)
	request.Header.Add("Content-Type", "application/json")

	// Ignore error since non-2xx status code doesn't cause any errors.
	// However, if Discord's API is down, the request won't work. In the future, this will
	// most likely cause some unwanted errors.
	res, _ := client.Do(request)

	// Ignore error since Discord's response body should always be appropriate
	rawResponse, _ := io.ReadAll(res.Body)

	return rawResponse, res.StatusCode
}
