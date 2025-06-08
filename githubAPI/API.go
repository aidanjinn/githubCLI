package githubAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Output interface {
	Str() string
}

func ToString(o Output) {
	fmt.Println(o.Str())
}

type GithubUser struct {
	UserActions []Action
	ActionStats map[string]int
}

func (g *GithubUser) generateStats() {
	g.ActionStats = make(map[string]int)
	// We have counters for all github actions types we encounter
	// When looking through our users github actions
	for _, action := range g.UserActions {

		_, exists := g.ActionStats[action.Type]
		if !exists {
			g.ActionStats[action.Type] = 1
		} else {
			g.ActionStats[action.Type]++
		}

	}
}

func (g *GithubUser) Str() string {
	var sb strings.Builder

	// String Building for total stats held in users ActionStats
	sb.WriteString("Total Stats: \n")
	for key, value := range g.ActionStats {
		sb.WriteString(key + " : " + strconv.Itoa(value) + "\n")
	}

	// String Building for each individual action
	for _, action := range g.UserActions {
		sb.WriteString(action.Str())
	}

	return sb.String()
}

type Action struct {
	Type    string  `json:"type"`
	Actor   Actor   `json:"actor"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

func (a *Action) Str() string {

	var output string

	if len(a.Payload.Commits) <= 0 {
		output += fmt.Sprintf("- %s performed %s on %s\n", a.Actor.DisplayName, a.Type, a.Repo.RepoName)
	} else {
		output += fmt.Sprintf("- %s performed %s on %s -- Message: %s\n", a.Actor.DisplayName, a.Type, a.Repo.RepoName, a.Payload.Commits[0].Message)
	}
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

	//First we load the data for the user
	ret := GithubUser{UserActions: data}
	//Then we generate the total stats based on that loaded data
	ret.generateStats()
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
