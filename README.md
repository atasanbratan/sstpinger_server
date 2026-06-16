# SSTPinger Minimal Go API for Vercel

This project provides two Vercel Go serverless functions:

- `api/add` - POST a JSON array of server objects to append to storage.
- `api/list` - GET with optional `offset` and `limit` query params to read latest entries.

Storage is implemented by storing a JSON file in a GitHub repository using the GitHub Contents API. Set the following environment variables in Vercel:

- `GITHUB_OWNER` — repository owner
- `GITHUB_REPO` — repository name
- `GITHUB_FILE_PATH` — path to store JSON (e.g. `data/servers.json`)
- `GITHUB_TOKEN` — GitHub Personal Access Token with `repo` scope

Example POST:

```bash
curl -X POST https://your-vercel.app/api/add \
  -H 'Content-Type: application/json' \
  -d '[{"id":1,"hostname":"h","ip":"1.2.3.4"}]'
```

Example GET:

```bash
curl 'https://your-vercel.app/api/list?offset=0&limit=50'
```
