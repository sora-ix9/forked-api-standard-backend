# Fakduai Standard API

Fakduai Standard Project for API

## Table of Contents

- [Project Structure](#project-structure)
- [Setup](#setup)
- [Usage](#usage)
- [Best Pracites](#best-pracites)
- [Testing](#testing)
- [FAQ](#faq)

## Project Structure

```bash
.
├── api-collection                  # Postman/Rest Client collection files
│   ├── helpcheck.http
│   ├── role.http
│   └── user.http
├── cmd
│   ├── main.go                     # Application entry point
│   └── tools.go
├── credentials                     # Service account credentials (ignored by git)
├── docs                            # Project documentation
│   ├── database-schema.html
│   ├── database-schema.mermaid
│   └── swagger.yaml
├── internal
│   ├── dto                         # Data Transfer Objects
│   │   ├── jwt_data.go
│   │   ├── role_data.go
│   │   ├── standard_response_data.go
│   │   └── user_data.go
│   ├── echo
│   │   └── server.go               # Echo server configuration
│   ├── handler                     # HTTP Request Handlers
│   │   ├── role_handler.go
│   │   ├── user_handler.go
│   │   └── websocket_handler.go
│   ├── middlewares                 # Custom Middlewares
│   │   ├── jwtauth_middleware.go
│   │   ├── request_log_middleware.go
│   │   └── role_middleware.go
│   ├── models                      # Domain Models
│   │   ├── role.go
│   │   └── user.go
│   ├── repositories                # Data Access Layer
│   │   ├── role_repository.go
│   │   └── user_repository.go
│   ├── routes                      # Route Registrations
│   │   ├── role_routes.go
│   │   └── user_routes.go
│   ├── services                    # Business Logic Layer
│   │   ├── role_service.go
│   │   ├── user_service.go
│   │   └── websocket_service.go
│   ├── tests                       # Tests (Unit & Integration)
│   │   ├── handlers                # Handler Unit Tests
│   │   │   ├── role_handler_test.go
│   │   │   ├── user_handler_test.go
│   │   │   └── websocket_handler_test.go
│   │   ├── integration             # Integration Tests
│   │   │   ├── role_routes_test.go
│   │   │   ├── setup_test.go
│   │   │   └── user_routes_test.go
│   │   ├── mock                    # Mocks for Testing
│   │   │   ├── repositories
│   │   │   │   ├── role_repository_mock.go
│   │   │   │   └── user_repository_mock.go
│   │   │   └── services
│   │   │       ├── role_service_mock.go
│   │   │       ├── user_service_mock.go
│   │   │       └── websocket_service_mock.go
│   │   └── services                # Service Unit Tests
│   │       ├── role_service_test.go
│   │       ├── user_service_test.go
│   │       └── websocket_service_test.go
│   ├── types                       # Type Definitions
│   │   └── uuid.go
│   └── utils                       # Utilities
│       ├── response.go
│       └── timeutil.go
├── k6_loadtest                     # Load Testing Scripts
│   ├── README.md
│   └── k6_loadtest.js
├── migrations                      # Database Migrations
│   ├── 20240427172736.sql
│   ├── 20240427173848.sql
│   ├── 20240502084112.sql
│   └── atlas.sum
├── pkg                             # Public Library Code
│   ├── config                      # Configuration Loading
│   │   └── config.go
│   ├── db                          # Database Connection
│   │   └── db.go
│   ├── providers                   # External Service Providers
│   │   ├── email
│   │   │   └── email.go
│   │   ├── firebase
│   │   │   └── firebase.go
│   │   ├── httpclient
│   │   │   └── httpclient.go
│   │   ├── logger
│   │   │   └── logger.go
│   │   ├── mongodb
│   │   │   └── mongodb.go
│   │   ├── s3
│   │   │   └── s3.go
│   │   ├── stripe
│   │   │   └── stripe.go
│   │   └── taskqueue
│   │       └── taskqueue.go
│   └── redisclient                 # Redis Client
│       └── redisclient.go
├── Dockerfile.dev
├── Dockerfile.prod
├── LICENSE
├── README.md
├── atlas.hcl
├── coverage.out
├── go.mod
└── go.sum
```

## Setup
### Load package
```bash
go mod tidy
```

### Database setup
MongoDB is used. Ensure you have a MongoDB instance running and update `MONGODB_URI` in `.env`.

## Usage
### Run project
```bash
go run ./cmd/main.go
```

### Run Live Reload
```bash
air
```

## Best Practise
### File Naming
- Descriptive Names: Files should be named clearly and descriptively. For example, user_service.go for services related to users.
- Use Snake Case: File names should use snake case (e.g., user_handler.go).
### Variable Naming
- CamelCase Usage: Use CamelCase for naming variables. Local variables should start with a lowercase letter (e.g., localVariable), and exported variables should start with an uppercase letter (e.g., ExportedVariable).
- Descriptive and Concise: Names should be both descriptive and concise. Prefer userID over id.
- Acronyms and Initialisms: Keep acronyms uppercase, e.g., HTTPServer, userID.

## Testing

ในโปรเจคนี้เรามีการทำ Test อยู่ 2 รูปแบบ คือ Unit Test และ Integration Test ซึ่งมีความสำคัญและวัตถุประสงค์ที่ต่างกัน

### Unit Test vs Integration Test (สำหรับน้องๆ ฝึกงาน)

*   **Unit Test (การทดสอบหน่วยย่อย):**
    *   **คืออะไร?**: การทดสอบโค้ดเฉพาะจุดเล็กๆ (เช่น Function เดียว หรือ Method เดียว) ว่าทำงานถูก logic ภายในตัวมันเองหรือไม่
    *   **Concept**: *ตัดขาดภายนอกทั้งหมด* (Isolate) เช่น ถ้า Service ต้องต่อ Database เราจะ "จำลอง" (Mock) Database ขึ้นมาแทน เพื่อให้เราเทสแค่ Logic ของ Service จริงๆ ไม่ต้องสนว่า Database จะล่มไหม
    *   **จุดเด่น**: รันไวมาก (ระดับ ms), หาจุดบั๊กง่ายเพราะเทสจุดเล็กๆ
    *   **ไฟล์อยู่ที่**: `internal/tests/services`, `internal/tests/handlers`

*   **Integration Test (การทดสอบการทำงานร่วมกัน):**
    *   **คืออะไร?**: การทดสอบว่าหลายๆ components ทำงานร่วมกันได้ถูกต้องหรือไม่ (เช่น API -> Service -> Database -> Response)
    *   **Concept**: *ใช้งานจริง* พยายามใช้ของจริงให้มากที่สุด เช่นใช้ Database จริง (ในที่นี้เราใช้ SQLite In-Memory จำลองให้เหมือนจริง) ยิง Request เข้ามาจริงๆ ผ่าน Router เพื่อดูว่า Flow ตั้งแต่ต้นจนจบทำงานถูกต้องไหม
    *   **จุดเด่น**: มั่นใจได้ว่าระบบภาพรวมทำงานได้จริง (User สมัครสมาชิกได้จริงๆ login ได้จริงๆ)
    *   **ไฟล์อยู่ที่**: `internal/tests/integration`

---

### คำสั่งสำหรับรัน Test (Run Commands)

#### 1. การรัน Unit Test (Run Unit Tests)
รัน Unit Test ทั้งหมดในโปรเจค (ทั้ง Handlers และ Services)

```bash
go test -v ./internal/tests/services/... ./internal/tests/handlers/...
```

#### 2. การรัน Integration Test (Run Integration Tests)
รัน Integration Test เพื่อดู Flow การทำงานจริง

```bash
go test -v ./internal/tests/integration/...
```

#### 3. การรัน Test ทั้งหมด (Run All Tests)
รันทั้ง Unit Test และ Integration Test พร้อมกัน

```bash
go test -v ./internal/tests/...
```

---

### การดู Test Coverage (Coverage Report)

หากต้องการดูว่า Test ครอบคลุม code เราไปกี่ % แล้ว ให้ใช้คำสั่งนี้:

**1. ดูตารางสรุป % Coverage**
```bash
go test -coverpkg=./internal/services/...,./internal/handler/... -cover ./internal/tests/...
```

**2. สร้าง Report แบบ HTML (ดูง่าย เห็นบรรทัดที่ยังไม่ได้เทส)**
```bash
# สร้างไฟล์ coverage profile
go test -coverpkg=./internal/services/...,./internal/handler/... -coverprofile=coverage.out ./internal/tests/...

# เปิดหน้าเว็บ report
go tool cover -html=coverage.out
```

## FAQ
## [Redis] Run on docker
```bash
docker run --name local-redis -d -e REDIS_PASSWORD='redispassword' redis redis-server --requirepass redispassword
```

## [Air] Setup Air (Live Reload)
```bash
# Open zshrc config
nano ~/.zshrc

# Add to the end of file
export PATH=$PATH:$(go env GOPATH)/bin

# Update zshrc config
source ~/.zshrc
```
# install air
```
go install github.com/air-verse/air@latest
```
-  close terminal and new teminal -> run air