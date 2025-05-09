# 🎟️ Ticket Score Engine

A Go-based ticket scoring service with gRPC support, backed by SQLite. It aggregates and computes weighted category scores over time based on customer ratings.

---

## 🚀 Getting Started

### ✅ Prerequisites

- **Go** (v1.20+ recommended)
- **SQLite**
- **Protocol Buffers Compiler (`protoc`)**
- Required Go tools:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

---

## 📦 Installation

1. **Initialize Go Module**
   ```bash
   go mod init <your-module-name>
   ```

2. **Install SQLite Driver**
   ```bash
   go get github.com/mattn/go-sqlite3
   ```

3. **Install `protoc`**
   - Download and install `protoc` from: https://github.com/protocolbuffers/protobuf/releases
   - Add it to your system `PATH`.

4. **Generate gRPC Code**
   ```bash
   protoc --go_out=generated --go-grpc_out=generated \
     --go_opt=paths=source_relative \
     --go-grpc_opt=paths=source_relative \
     --proto_path=api/proto api/proto/scoring.proto
   ```

---

## 🧪 Running the Project

Start the gRPC server:
```bash
go run cmd/server/main.go
```

---

## 🗂️ Project Structure

```
ticket-score-engine/
├── cmd/
│   └── server/               # Entry point for the gRPC server
│       └── main.go
├── api/
│   └── proto/                # Protobuf definitions
│       └── scoring.proto
├── generated/                # gRPC generated Go files
├── internal/
│   ├── db/                   # SQLite setup and helpers
│   ├── rating/               # Scoring logic and aggregation
│   └── server/               # gRPC server implementation
├── pkg/                      # Reusable/shared packages (optional)
├── database.db               # SQLite file (dev/test only)
├── go.mod
├── go.sum
└── Dockerfile                # Containerization (optional)
```

---

> ⚠️ **Note**: This README will be improved and expanded as the project evolves (e.g., Docker support, CI/CD, extended API usage).
