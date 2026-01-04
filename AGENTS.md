# BigQuery Emulator — Agent Guide
 
## Project overview
 
This repository implements a BigQuery emulator server in Go.
 
It provides:
 
- A standalone `bigquery-emulator` CLI that serves:
  - BigQuery REST API (default port `9050`)
  - BigQuery Storage API over gRPC (default port `9060`)
- A Go library API (via `github.com/goccy/bigquery-emulator/server`) for embedding in tests.
 
Storage is backed by SQLite via `go-zetasqlite` (see `README.md` “How it works”).
 
## Key configuration files
 
- `go.mod`: module definition (`github.com/goccy/bigquery-emulator`) and Go version.
- `Makefile`: canonical build targets (notably `emulator/build`).
- `Dockerfile`: container build (uses `ghcr.io/goccy/go-zetasql`).
- `.github/workflows/*.yml`: CI build/test and release workflows.
 
## Repo layout
 
- `cmd/bigquery-emulator/`: CLI entrypoint.
- `server/`: core server implementation (HTTP + gRPC, middleware, lifecycle).
- `internal/`: internal packages (metadata/content repositories, request handling, etc.).
- `types/`: types used to load/seed project/dataset/table data.
 
## Build and test commands
 
### Build (recommended)
 
The repo’s `Makefile` builds the CLI and sets version metadata.
 
```bash
make emulator/build
```
 
Notes from the repo:
 
- CI builds with `CC=clang` and `CXX=clang++`.
- The emulator depends on `go-zetasql`, which requires `CGO_ENABLED=1` during build.
 
### Run tests
 
```bash
go test -v ./... -count=1
```
 
(This is also what CI runs in `.github/workflows/test.yml`.)
 
### Docker
 
```bash
docker build -t bigquery-emulator .
```
 
Or via Make:
 
```bash
make docker/build
```
 
## Running the emulator
 
The CLI requires `--project`.
 
```bash
./bigquery-emulator --project=test
```
 
Useful flags (see `README.md` / `cmd/bigquery-emulator/main.go`):
 
- `--dataset`
- `--port` (REST)
- `--grpc-port` (gRPC)
- `--database` (SQLite file; if unset, uses temp/in-memory storage)
- `--data-from-yaml` (seed data)
- `--log-level`, `--log-format`
 
## Runtime architecture (high level)
 
- `cmd/bigquery-emulator/main.go` parses flags and starts the server.
- `server.New(storage)` constructs a server backed by SQLite (`zetasqlite` driver).
- `Server.Serve(ctx, httpAddr, grpcAddr)` starts:
  - an HTTP server (REST API)
  - a gRPC server (BigQuery Storage API)
 
## Development conventions
 
- Prefer `make emulator/build` for reproducible builds (it matches CI).
- Be aware builds may take a long time due to ZetaSQL compilation (called out in `README.md`).
 
## Security considerations
 
- This is an emulator intended for local testing/development.
- Avoid exposing it to untrusted networks; the default host is `0.0.0.0`.
