package tests

import (
	"githubCLI/githubAPI"
	"io"
	"net/http"
	"strings"
	"testing"
)

const jsonData = `[{
		"id": "50586368742",
		"type": "PushEvent",
		"actor": {
			"id": 96208746,
			"login": "aidanjinn",
			"display_login": "aidanjinn",
			"gravatar_id": "",
			"url": "https://api.github.com/users/aidanjinn",
			"avatar_url": "https://avatars.githubusercontent.com/u/96208746?"
		},
		"repo": {
			"id": 438780222,
			"name": "aidanjinn/aidanjinn",
			"url": "https://api.github.com/repos/aidanjinn/aidanjinn"
		},
		"payload": {
			"repository_id": 438780222,
			"push_id": 24735860700,
			"size": 1,
			"distinct_size": 1,
			"ref": "refs/heads/main",
			"head": "6401fba434db35583b566910f7373100fce823ac",
			"before": "6c212deafbaea9a94c29f62bbc83a858b0e4267c",
			"commits": [
			{
				"sha": "6401fba434db35583b566910f7373100fce823ac",
				"author": {
				"email": "96208746+aidanjinn@users.noreply.github.com",
				"name": "Aidan Wendorf"
			},
			"message": "Update README.md",
			"distinct": true,
			"url": "https://api.github.com/repos/aidanjinn/aidanjinn/commits/6401fba434db35583b566910f7373100fce823ac"
		}]},
		"public": true,
		"created_at": "2025-06-06T05:47:49Z"
	}]`

func errOut(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("ping() returned an error: %v", err)
	}
}

func TestDriver(t *testing.T) {
	TestPing(t)
	TestJSON(t)
	TestToString(t)
}

// Pinging the Github API and receiving a proper response
// Function Name we are testing is ping
func TestPing(t *testing.T) {

	events := 10

	// This should return a proper response for my github account activity
	statusCode, userData, err := githubAPI.Ping("aidanjinn", events)

	errOut(t, err)

	// Failure for not getting a response at all
	if statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if len(userData.UserActions) == 0 {
		t.Errorf("Expected userActions, got 0")
	}
}

func TestJSON(t *testing.T) {

	resp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(jsonData)),
		StatusCode: http.StatusOK,
	}

	userData, err := githubAPI.LoadJSON(resp)
	errOut(t, err)

	actions := userData.UserActions
	if len(actions) != 1 {
		t.Fatalf("Expected 1 action, got %d", len(actions))
	}

	action := actions[0]
	if action.Type != "PushEvent" {
		t.Errorf("Expected action type 'PushEvent', got '%s'", action.Type)
	}
	if action.Actor.DisplayName != "aidanjinn" {
		t.Errorf("Expected actor 'aidanjinn', got '%s'", action.Actor.DisplayName)
	}
	if action.Repo.RepoName != "aidanjinn/aidanjinn" {
		t.Errorf("Expected repo name 'aidanjinn/test-repo', got '%s'", action.Repo.RepoName)
	}
	if action.Payload.Commits[0].Author.Name != "Aidan Wendorf" {
		t.Errorf("Expected author name 'aidanjinn', got '%s'", action.Payload.Commits[0].Author.Name)
	}
	if action.Payload.Commits[0].Message != "Update README.md" {
		t.Errorf("Expected message 'Update README.md', got '%s'", action.Payload.Commits[0].Message)
	}

}

func TestToString(t *testing.T) {
	resp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(jsonData)),
		StatusCode: http.StatusOK,
	}

	userData, err := githubAPI.LoadJSON(resp)
	errOut(t, err)

	actions := userData.UserActions

	userDataCorrectString := "Total Stats: \n"
	userDataCorrectString += "PushEvent : 1\n"
	userDataCorrectString += "- aidanjinn performed PushEvent on aidanjinn/aidanjinn -- Message: Update README.md\n"
	if userData.Str() != userDataCorrectString {
		t.Errorf("Expected actions to be 'PushEvent', got '%v'", actions)
	}

	correctString := "- aidanjinn performed PushEvent on aidanjinn/aidanjinn -- Message: Update README.md\n"
	if actions[0].Str() != correctString {
		t.Errorf("Expected actions to be 'PushEvent', got '%v'", actions)
	}
}
