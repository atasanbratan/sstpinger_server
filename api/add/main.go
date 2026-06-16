package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"

    "github.com/example/sstpinger/pkg/storage"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }

    var items []map[string]interface{}
    if err := json.Unmarshal(body, &items); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    if err := storage.AppendItems(os.Getenv("GITHUB_OWNER"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_FILE_PATH"), os.Getenv("GITHUB_TOKEN"), items); err != nil {
        http.Error(w, "failed to store data: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{"added": len(items)})
}

func main() {
    http.HandleFunc("/", Handler)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
