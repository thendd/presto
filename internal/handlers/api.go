package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"presto/configs"
)

// TODO: using a client in the way I'm using right now might lead into some problems in the future
// so I'll look further into that later
var client = http.Client{}

// Creates an `http.Request` object and does a request using a global `http.Client`
func MakeHTTPRequestToDiscord(endpoint string, method configs.HTTPMethod, body any) string {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		log.Println(err)
	}

	// Ignore error just because (documentatio doesn't say anything about the cases in
	// which an error is returned).
	request, _ := http.NewRequest(string(method), configs.API_URL+"/"+endpoint, &buffer)

	request.Header.Add("Authorization", "Bot "+configs.DISCORD_BOT_TOKEN)
	request.Header.Add("Content-Type", "application/json")

	// Ignore error since non-2xx status code doesn't cause any errors.
	// However, if Discord's API is down, the request won't work. In the future, this will
	// most likely cause some unwanted errors.
	res, err := client.Do(request)

	// Ignore error since Discord's response body should always be appropriate
	rawResponse, _ := io.ReadAll(res.Body)

	return string(rawResponse)
}
