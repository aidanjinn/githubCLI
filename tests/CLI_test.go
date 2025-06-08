package tests

import (
	"githubCLI/githubAPI"
	"io"
	"net/http"
	"strings"
	"testing"
)

func errOut(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("ping() returned an error: %v", err)
	}
}

func TestDriver(t *testing.T) {
	TestPing(t)
	TestJSON(t)
}

// Pinging the Github API and receiving a proper response
// Function Name we are testing is ping
func TestPing(t *testing.T) {

	// This should return a proper response for my github account activity
	statusCode, err := githubAPI.Ping("aidanjinn")

	errOut(t, err)

	// Failure for not getting a response at all
	if statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, statusCode)
	}
}

func TestJSON(t *testing.T) {
	jsonData := `[{
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
	resp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(jsonData)),
		StatusCode: http.StatusOK,
	}

	actions, err := githubAPI.LoadJSON(resp)
	errOut(t, err)

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
