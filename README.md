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
/src
|-- /cmd
|   |-- main.go                     # Entry point for the application.
|-- /internal
|   |-- /dto
|       |-- user_data.go            # Data transfer objects.
|   |-- /echo
|       |-- server.go               # Echo framework setup.
|   |-- /handler
|       |-- user_handler.go         # HTTP handlers for user operations.
|   |-- /middlewares
|       |-- role_middleware.go      # Middleware for role checking.
|   |-- /models
|       |-- user.go                 # User model definition.
|       |-- role.go                 # Role model definition.
|   |-- /repositories
|       |-- user_repository.go      # User repositories for DB operations.
|   |-- /services
|       |-- user_service.go         # Business logic for user management.
|   |-- /types
|       |-- uuid.go                 # Type definitions, such as UUID.
|   |-- /util
|       |-- response.go             # Utility function for response with standard format.
|   |-- /tests                      # Unit and Integration Tests
|       |-- /handlers               # Unit tests for handlers
|       |-- /services               # Unit tests for services
|       |-- /integration            # Integration tests
|       |-- /mock                   # Mocks for unit testing
|-- /pkg
|   |-- /config
|       |-- config.go               # Configuration setup.
|-- /db
|   |-- /db.go                      # Database connection and setup.
|-- /redisclient
|   |-- /redisclient.go             # Redis client configuration.
|-- go.mod                          # Go module dependencies.
|-- go.sum                          # Go module checksums.
|-- .env                            # Environment variables.
|-- README.md                       # This file.

```

## Setup
### Load package
```bash
go mod tidy
```

### Postgresql run on docker
```bash
# docker run --name local-postgres -e POSTGRES_PASSWORD=password -d postgres
docker run --name local-postgres -e POSTGRES_PASSWORD=password -p 54320:5432 -d postgres
```

## Usage
### Run project
```bash
go run ./cmd/main.go
```

### Run Live Reload
```bash
air
```

### [DB] Create Migration file
```bash
atlas migrate diff --env local
```

### [DB] Apply with migration
```bash
atlas migrate apply --env local --url "postgres://postgres:password@localhost:54320/fdlp-dev-db?search_path=public&sslmode=disable"
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
### [Postgresql] Fix UUID not found
```bash
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```
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