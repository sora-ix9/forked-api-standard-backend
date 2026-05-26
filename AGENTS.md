---
trigger: always_on
---

You are an expert in Go (Golang) programming and best practices.

- ก่อนจะเริ่มทำอะไรให้สร้างไฟล์ TODO.md ทุกครั้ง จากนั้นก็
- ออกแบบ unittest , integration ก่อนเพิ่ม implement กระบวนการ TDD ลงใน TODO.md ถ้ามี unittest ก่อนหน้าอยู่แล้วให้ไปใช้ไฟล์นั้นได้เลยแล้วเพิ่ม function test logic นั้นๆเอา
- implement unittest integratuin ก่อนเริ่ม implement code
- Smart people manage to come up with simple solutions for difficult problems, and dumb people do the opposite. Read TODO.md, find the first unchecked task [ ], complete it, update TODO.md to mark it [x], save the file, then repeat until all tasks are done. Use the least amount of comments possible.
- เรียบร้อยแล้วลบ TODO.md ออกไป ไม่ต้องลบ test ออกนะถ้าเป็น services/ ถ้าอยู่ที่อื่นลบออกให้หมดเลย
- update swagger.ymal ให้ด้วยหล่ะ แบบเขียนเองไม่ต้อง generate นะฉันไม่ชอบมันลกโค้ดไม่ชอบเอาด้วยที่นี่เราไม่ทำกันแบบนั้น

Skill Usage:
- Always invoke `/karpathy-guidelines` before writing, reviewing, or refactoring code. It is short and prevents overengineering.
- Use `/diagnose` only for unclear bugs, failing behavior, hard-to-reproduce issues, or performance regressions. Build a fast repro loop before reading too much code.
- Use `/golang-performance` only after profiling, benchmarks, or production metrics identify a bottleneck. Do not spend tokens on it for normal Go changes.
- Use `/scrutinize` for plan, PR, diff, design, or code-review second opinions. Use it before merge when risk is non-trivial; skip for tiny mechanical edits.
- Use `/improve-codebase-architecture` when user asks for architecture improvement, refactoring opportunities, testability, module seams, or AI-navigability. Do not use during normal feature work unless architecture blocks testing or clarity.
- Use `/to-issues` when converting a plan, spec, PRD, or large task into tracker issues. Prefer after plan is stable, not during early exploration.
- Token rule: choose the narrowest skill that matches the task; do not chain skills unless output from one clearly feeds the next.

Key Principles:
- Follow idiomatic Go (Effective Go)
- Keep it simple and readable
- Handle errors explicitly
- Prefer composition over inheritance
- Use goroutines for concurrency

Code Organization:
- Use standard project layout (cmd/, pkg/, internal/)
- Group related code in packages
- Keep packages small and focused
- Use meaningful package names
- Avoid circular dependencies

Naming Conventions:
- Use CamelCase for exported names
- Use camelCase for unexported names
- Keep names short and concise
- Use single-letter names for short loops/scopes
- Avoid stuttering (e.g., user.UserInfo -> user.Info)

Error Handling:
- Check errors immediately after function calls
- Return errors as the last return value
- Use custom error types for specific cases
- Wrap errors with context (fmt.Errorf("%w"))
- Don't panic unless truly unrecoverable

Functions and Methods:
- Keep functions short and focused
- Use named return values sparingly
- Use defer for cleanup
- Use interfaces for flexibility
- Accept interfaces, return structs

Data Structures:
- Use slices over arrays
- Use maps for key-value storage
- Use structs for grouping data
- Use pointers for large structs or mutability
- Initialize structs with field names

Architecture Patterns:
- Hexagonal Architecture (Ports and Adapters)
- Clean Architecture
- Domain-Driven Design (DDD)
- Event-Driven Architecture
- CQRS (Command Query Responsibility Segregation)

Concurrency:
- Use goroutines for concurrent tasks
- Use channels for communication
- Use sync.WaitGroup to wait for goroutines
- Use sync.Mutex for shared state
- Avoid sharing memory by communicating

Goroutines:
- Start goroutines with 'go' keyword
- Keep goroutines lightweight
- Manage goroutine lifecycle
- Avoid leaking goroutines
- Use WaitGroup to wait for completion

Channels:
- Use unbuffered channels for synchronization
- Use buffered channels for throughput
- Close channels from the sender side
- Use range to iterate over channels
- Use select for multiplexing channels

Synchronization:
- Use sync.Mutex for critical sections
- Use sync.RWMutex for read-heavy data
- Use sync.Once for one-time initialization
- Use sync.Cond for signaling
- Use atomic package for simple counters

Testing:
- Write unit tests in _test.go files
- Use the testing package
- Use table-driven tests
- Run tests with go test
- Use go test -race to check for race conditions

Dependency Management:
- Use Go Modules (go.mod)
- Keep dependencies minimal
- Vendor dependencies if necessary
- Use semantic versioning
- Audit dependencies regularly

Formatting and Linting:
- Always run gofmt
- Use go vet to catch common errors
- Use golangci-lint for comprehensive linting
- Follow community style guides
- Document exported names with comments

Best Practices:
- Handle all errors
- Avoid global state
- Use context for cancellation and timeouts
- Write benchmarks for performance-critical code
- Keep main() simple
- Use standard library when possible