# ticket-score-engine

Installed Go
ran go mod init <repository_url> for initialization of the go module
added main.go and for build and run
   go run cmd/server/main.go

Install SQLite driver for Go
go get github.com/mattn/go-sqlite3



The strcture will be 

ticket-score-engine/
├── cmd/
│   └── server/               # Main server entry point
│       └── main.go
├── api/                      # .proto definitions (vs. mixing with generated code)
│   └── rating.proto
├── internal/
│   ├── rating/               # Business logic for rating scoring
│   ├── db/                   # DB access layer (SQLite handling)
│   └── server/               # gRPC server and handler logic
├── pkg/                      # (Optional) Reusable utility code, not private
├── scripts/                  # Helpful bash/devops scripts
├── third_party/              # Any vendored or custom .proto definitions
├── database.db               # SQLite file (for dev/testing only)
├── go.mod
└── Dockerfile


