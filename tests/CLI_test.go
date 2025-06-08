package tests

import (
	"githubCLI/githubAPI"
	"net/http"
	"testing"
)

func errOut(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("ping() returned an error: %v", err)
	}
}

// Pinging the Github API and receiving a proper response
// Function Name we are testing is ping
func TestPing(t *testing.T) {
	statusCode, resp, err := githubAPI.Ping("aidanjinn")

	errOut(t, err)

	if statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if len(resp.Data) == 0 {
		t.Errorf("Expected a response but recived an empty response buffer")
	}
}
