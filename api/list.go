package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/example/sstpinger/pkg/storage"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	qs := r.URL.Query()
	off, _ := strconv.Atoi(qs.Get("offset"))
	limit, _ := strconv.Atoi(qs.Get("limit"))
	if limit <= 0 {
		limit = 100
	}

	items, err := storage.ReadAllItems(os.Getenv("GITHUB_OWNER"), os.Getenv("GITHUB_REPO"), os.Getenv("GITHUB_FILE_PATH"), os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		http.Error(w, "failed to read data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	start := off
	if start < 0 {
		start = 0
	}
	if start > len(items) {
		start = len(items)
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}

	json.NewEncoder(w).Encode(items[start:end])
}
