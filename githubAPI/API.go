package githubAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Output interface {
	str() string
}

func ToString(o Output) {
	fmt.Println(o.str())
}

type GithubUser struct {
	UserActions []Action
}

func (g *GithubUser) str() string {

	var output string

	for _, action := range g.UserActions {
		output += action.str()
	}

	return output
}

type Action struct {
	Type    string  `json:"type"`
	Actor   Actor   `json:"actor"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

func (a *Action) str() string {

	var output string

	output += "- " + a.Actor.DisplayName + " " + a.Type + " : " + a.Repo.RepoName + "\n"

	return output
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

func LoadJSON(resp *http.Response) (GithubUser, error) {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return GithubUser{}, err
	}

	var data []Action
	if err := json.Unmarshal(body, &data); err != nil {
		return GithubUser{}, err
	}

	ret := GithubUser{UserActions: data}
	return ret, nil
}

func Ping(username string, events int) (int, GithubUser, error) {
	client := &http.Client{}
	perPage := strconv.Itoa(events)
	url := "https://api.github.com/users/" + username + "/events/public?per_page=" + perPage
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, GithubUser{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)

	if err != nil {
		return 0, GithubUser{}, err
	}
	// This will close the resp once the ping function ends
	defer resp.Body.Close()

	userData, err := LoadJSON(resp)
	if err != nil {
		return 0, GithubUser{}, err
	}

	return resp.StatusCode, userData, nil
}
