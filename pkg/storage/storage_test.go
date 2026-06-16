package storage

import (
	"os"
	"testing"
	"time"
)

func getGitHubEnv(t *testing.T) (owner, repo, path, token string) {
	owner = os.Getenv("GITHUB_OWNER")
	repo = os.Getenv("GITHUB_REPO")
	path = os.Getenv("GITHUB_FILE_PATH")
	token = os.Getenv("GITHUB_TOKEN")

	if owner == "" || repo == "" || path == "" || token == "" {
		t.Skip("GitHub storage env vars not set")
	}
	return
}

func TestGitHubStorageAppendAndRead(t *testing.T) {
	owner, repo, path, token := getGitHubEnv(t)

	marker := time.Now().Unix()
	item := map[string]interface{}{
		"id":       marker,
		"hostname": "test-github-storage",
		"ip":       "127.0.0.1",
		"port":     1234,
		"key":      "127.0.0.1:1234",
		"sessions": 1,
		"info":     "1 SESSIONS 0.00",
		"info2":    "TEST 0.00",
		"location": map[string]interface{}{
			"country": "TEST",
			"short":   "TT",
			"name":    "TEST STORAGE",
		},
	}

	if err := AppendItems(owner, repo, path, token, []map[string]interface{}{item}); err != nil {
		t.Fatalf("AppendItems failed: %v", err)
	}

	items, err := ReadAllItems(owner, repo, path, token)
	if err != nil {
		t.Fatalf("ReadAllItems failed: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected at least one item after append")
	}

	found := false
	for _, entry := range items {
		idValue, ok := entry["id"]
		if !ok {
			continue
		}
		if float64(marker) == idValue.(float64) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("appended item with id %d not found", marker)
	}
}
