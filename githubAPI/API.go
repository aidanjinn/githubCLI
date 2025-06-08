package githubAPI

import (
	"encoding/json"
	"io"
	"net/http"
)

type GithubUser struct {
	UserActions []Action
}

type Action struct {
	Type    string  `json:"type"`
	Actor   Actor   `json:"actor"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

type Actor struct {
	DisplayName string `json:"display_login"`
	Url         string `json:"url"`
}

type Repo struct {
	RepoName string `json:"name"`
	Url      string `json:"url"`
}

type Payload struct {
	Commits []Commits `json:"commits"`
}

type Commits struct {
	Author  Author `json:"author"`
	Message string `json:"message"`
}

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func LoadJSON(resp *http.Response) ([]Action, error) {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []Action{}, err
	}

	var data []Action
	if err := json.Unmarshal(body, &data); err != nil {
		return []Action{}, err
	}

	return data, nil
}

func Ping(username string) (int, error) {
	client := &http.Client{}
	url := "https://api.github.com/users/" + username + "/events/public?per_page=1"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	// This will close the resp once the ping function ends
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
