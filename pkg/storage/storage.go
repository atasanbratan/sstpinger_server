package storage

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
)

func apiURL(owner, repo, path string) string {
    return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)
}

// ReadAllItems fetches the JSON array stored in the configured GitHub file.
func ReadAllItems(owner, repo, path, token string) ([]map[string]interface{}, error) {
    if owner == "" || repo == "" || path == "" || token == "" {
        return nil, errors.New("missing GitHub storage environment variables")
    }

    client := &http.Client{Timeout: 15 * time.Second}
    req, _ := http.NewRequest("GET", apiURL(owner, repo, path), nil)
    req.Header.Set("Authorization", "token "+token)
    req.Header.Set("Accept", "application/vnd.github.v3.raw")
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return []map[string]interface{}{}, nil
    }
    if resp.StatusCode >= 400 {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("github api error: %s", string(b))
    }

    var items []map[string]interface{}
    dec := json.NewDecoder(resp.Body)
    if err := dec.Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}

// AppendItems appends new items to the stored JSON array using the GitHub contents API.
func AppendItems(owner, repo, path, token string, newItems []map[string]interface{}) error {
    items, err := ReadAllItems(owner, repo, path, token)
    if err != nil {
        return err
    }
    items = append(items, newItems...)

    raw, err := json.MarshalIndent(items, "", "  ")
    if err != nil {
        return err
    }

    // Need current sha if file exists
    sha, _ := getFileSHA(owner, repo, path, token)

    payload := map[string]interface{}{
        "message": "update sstp servers",
        "content": base64.StdEncoding.EncodeToString(raw),
    }
    if sha != "" {
        payload["sha"] = sha
    }

    body, _ := json.Marshal(payload)
    client := &http.Client{Timeout: 15 * time.Second}
    req, _ := http.NewRequest("PUT", apiURL(owner, repo, path), bytes.NewReader(body))
    req.Header.Set("Authorization", "token "+token)
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("github put error: %s", string(b))
    }
    return nil
}

func getFileSHA(owner, repo, path, token string) (string, error) {
    if owner == "" || repo == "" || path == "" || token == "" {
        return "", errors.New("missing GitHub storage environment variables")
    }
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", apiURL(owner, repo, path), nil)
    req.Header.Set("Authorization", "token "+token)
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusNotFound {
        return "", nil
    }
    if resp.StatusCode >= 400 {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("github api error: %s", string(b))
    }
    var info struct{ Sha string `json:"sha"` }
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        return "", err
    }
    return info.Sha, nil
}
