# 🎟️ Ticket Score Engine

A Go-based ticket scoring service with gRPC support, backed by SQLite. It aggregates and computes weighted category scores over time based on customer ratings.

---

## 🚀 Getting Started

### ✅ Prerequisites

### Prerequisites
- Docker and Docker Compose installed

### Run the Application
```bash
# Start the service with SQLite (development)
docker-compose up --build
```

---
### Test endpoints

Install gRPCcurl (https://github.com/fullstorydev/grpcurl)

Navigate to the ```api/proto``` directory.

Then run below grpcurl commands

```bash
# Get category scores (plaintext)
grpcurl -plaintext \
  -import-path . \
  -proto scoring.proto \
  -d '{"start_date": "2020-01-01", "end_date": "2020-01-16"}' \
  localhost:50051 \
  scoring.ScoringService/GetCategoryScores

# Get ticket scores
grpcurl -plaintext \
  -import-path . \
  -proto scoring.proto \
  -d '{"start_date": "2020-01-01", "end_date": "2020-01-16"}' \
  localhost:50051 \
  scoring.ScoringService/GetTicketScores

# Get overall score
grpcurl -plaintext \
  -import-path . \
  -proto scoring.proto \
  -d '{"start_date": "2020-01-01", "end_date": "2020-01-16"}' \
  localhost:50051 \
  scoring.ScoringService/GetOverallScore

# Get period comparison
grpcurl -plaintext \
  -import-path . \
  -proto scoring.proto \
  -d '{
    "current_period": {
      "start_date": "2020-02-01",
      "end_date": "2020-02-28"
    },
    "previous_period": {
      "start_date": "2020-01-01",
      "end_date": "2020-01-31"
    }
  }' \
  localhost:50051 \
  scoring.ScoringService/GetPeriodComparison

```
There can be some issues with the quote and the formatting issues with above commands. 

Or else you can use Postman for invoking above endpoints. Simply import the          ```scoring.proto``` file in Postman

### For further improvments

For new APIs or changes, update ```scoring.proto``` and run below command
The generated gRPC code will be put in to the generated directory.

**Generate gRPC Code**
   ```bash
   protoc --go_out=generated --go-grpc_out=generated \
     --go_opt=paths=source_relative \
     --go-grpc_opt=paths=source_relative \
     --proto_path=api/proto api/proto/scoring.proto
   ```
---

## 🧪 Running the Project

If docker is not installed , use below command to start the gRPC server:
```bash
go run cmd/server/main.go
```

---

### gRPC Endpoints

| Service Method           | Request Type              | Response Type             | Description |
|--------------------------|---------------------------|---------------------------|-------------|
| `GetCategoryScores`      | `ScoreRequest`            | `ScoreResponse`           | Returns aggregated scores by category for a given time period (daily/weekly) |
| `GetTicketScores`        | `ScoreRequest`            | `TicketScoreResponse`     | Provides scores grouped by ticket ID with category breakdown |
| `GetOverallScore`        | `ScoreRequest`            | `OverallScoreResponse`    | Returns composite quality score across all categories |
| `GetPeriodComparison`    | `PeriodComparisonRequest` | `PeriodComparisonResponse`| Compares scores between two time periods |

View complete protocol buffer definition: ```api/proto/scoring.proto```

## 🗂️ Project Structure

```
ticket-score-engine/
├── api/
│   └── proto/                # Protocol Buffer definitions
│       ├── scoring.proto
│       └── generated/        # Auto-generated gRPC code
├── cmd/
│   └── server/               # Main application entrypoint
│       └── main.go
├── internal/
│   ├── domain/               # Core models
│   │   ├── category.go
│   │   ├── overall.go
│   │   └── ticket.go
│   ├── repository/           # Data access layer
│   │   ├── category_repo.go
│   │   ├── overall_repo.go
│   │   ├── ticket_repo.go
│   │   └── test/             # Repository unit tests => Data level testing
│   ├── scoring/              # Business logic/call to Data layer
│   │   ├── category_scores.go
│   │   ├── overall_scores.go
│   │   ├── ticket_scores.go
│   │   └── test/             # Business logic tests
│   └── server/               # gRPC server implementation
│       ├── grpc_server.go
├── generated                 # Auto-generated gRPC code from endpoints defined in proto
├── kubernetes/               # Kubernetes deployment files
│   ├── deployment.yml
│   └── service.yml
├── test                 # integration tests (grpc server tests)
│   ├── integration
├── Dockerfile                # Container build config
├── docker-compose.yml        # Dev environment
├── docker-compose.prod.yml   # Prod environment
├── go.mod
└── go.sum
```
## Kubernetes Deployment

Basic Kubernetes manifests are included in the `kubernetes/` directory to deploy the service in a cluster. These include:

- `deployment.yaml`: Deploys the gRPC service
- `service.yaml`: Exposes it internally

PS: I haven't tested the Kubernates deployment part and it is still incomplete. 

-----

For Questions please reach out ravindalakshan@gmail.com