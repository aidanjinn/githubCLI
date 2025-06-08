package githubAPI

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Data []byte
}

func Ping(username string) (int, Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/users/aidanjinn/events/public?per_page=1", nil)

	if err != nil {
		return 0, Response{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)

	if err != nil {
		return 0, Response{}, err
	}
	// This will close the resp once the ping function ends
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, Response{}, err
	}

	var jsonResp interface{}
	err = json.Unmarshal(body, &jsonResp)
	retResp := Response{Data: body}

	return resp.StatusCode, retResp, nil
}
