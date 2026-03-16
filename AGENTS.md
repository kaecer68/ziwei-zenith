# AGENTS.md - Ziwei Zenith Development Guide

## Build & Test Commands

```bash
# Build all packages
go build ./...

# Run all tests (none currently exist)
go test ./...

# Run single test file
go test -v ./pkg/engine/...

# Lint with go vet
go vet ./...

# Format code
go fmt ./...

# Run all checks
go build ./... && go vet ./... && go fmt ./...
```

## Development Commands

```bash
# Run CLI
go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10
go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10 -gender female -json

# Start server (REST :8081 + gRPC :50051)
go run ./cmd/ziwei-server/main.go

# Test REST API
curl -X POST -H "Content-Type: application/json" \
  -d '{"year":1972,"month":6,"day":8,"hour":2,"gender":"male"}' \
  http://localhost:8081/api/v1/calculate

# Test gRPC (requires grpcurl)
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"year":1972,"month":6,"day":8,"hour":1,"gender":"male"}' \
  localhost:50051 ziwei.v1.ZiweiService/Calculate
grpcurl -plaintext -d '{}' localhost:50051 ziwei.v1.ZiweiService/ListTags

# Regenerate gRPC code (after modifying proto/ziwei.proto)
protoc --go_out=pkg/api/grpc/v1 --go_opt=paths=source_relative \
  --go-grpc_out=pkg/api/grpc/v1 --go-grpc_opt=paths=source_relative \
  -I proto proto/ziwei.proto
```

---

## Code Style Guidelines

### Project Structure

```
ziwei-zenith/
├── cmd/
│   ├── ziwei-cli/          # CLI application
│   └── ziwei-server/       # REST + gRPC server
├── proto/
│   └── ziwei.proto         # Protobuf service definition
├── pkg/
│   ├── api/v1/             # REST JSON API types
│   ├── api/grpc/v1/        # Generated gRPC Go code
│   ├── basis/              # Core definitions (stars, palaces, wuxing, etc.)
│   ├── engine/             # Calculation engine
│   └── service/            # Shared service layer (calculate + gRPC server)
├── web/                    # React + TypeScript frontend
└── go.mod
```

### Package Organization

- **One file per concern**: `stars.go`, `palaces.go`, `brightness.go`
- **basis package**: All type definitions, constants, lookup tables
- **engine package**: Algorithms and business logic
- **service package**: Shared calculation logic, gRPC server implementation

### Types & Constants

```go
// Use iota for enum-like constants
type Star int

const (
    StarZiwei Star = iota
    StarTianfu
    StarTianji
    // ...
)

// Add Chinese comment for clarity
const StarZiwei Star = iota // 紫微星
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Types | PascalCase | `ZiweiEngine`, `Star`, `Palace` |
| Functions | PascalCase | `CalcLifePalace()`, `PlaceMainStars()` |
| Variables | camelCase | `lifePalace`, `starList` |
| Constants | PascalCase or camelCase | `StarZiwei`, `birthInfo` |
| Package | lowercase | `basis`, `engine` |

### Import Style

```go
import (
    "fmt"

    "github.com/kaecer68/ziwei-zenith/pkg/basis"
)
```

- Group stdlib first, then external
- No unused imports (will fail build)

### Error Handling

- Return `error` as last return value
- Use descriptive errors: `fmt.Errorf("invalid solar date: %w", err)`
- Handle errors at call site or propagate

```go
func (e *ZiweiEngine) BuildChart(birth BirthInfo) (*ZiweiChart, error) {
    // ... validation
    // ... calculations
    return chart, nil
}
```

### String Representation

- Implement `String() string` for all types for display
- Use `fmt.Sprintf` for structured output
- Return Traditional Chinese for user-facing strings

```go
func (s Star) String() string {
    names := []string{
        "紫微", "天府", "天機", // ...
    }
    return names[s]
}
```

### Map & Slice Usage

- Use maps for lookups: `var StarBrightnessTable = map[Star][]Brightness{...}`
- Use slices for ordered collections
- Check existence before map access: `if v, ok := m[key]; ok {...}`

### Type Conversions

- Use explicit conversions: `basis.Branch(value)`
- Avoid `any` or empty interface unless necessary

### Logic Patterns

```go
// Switch over type for interface handling
switch v := value.(type) {
case basis.AuspiciousStar:
    starStr += v.String() + " "
case basis.MaleficStar:
    starStr += v.String() + " "
}

// Modulo with proper handling
index := (base + offset) % 12
```

---

## Feature Implementation Workflow

1. Add types to `pkg/basis/` (definitions, lookup tables)
2. Add algorithms to `pkg/engine/` (calculation logic)
3. Update `ZiweiChart` struct in `engine.go`
4. Update `String()` method for output
5. If adding new API fields: update `proto/ziwei.proto` and regenerate gRPC code
6. If adding new API fields: update `pkg/service/grpc_server.go` conversion logic
7. Test with CLI: `go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10`
8. Test with server: `go run ./cmd/ziwei-server/main.go` then verify REST + gRPC

---

## Important Notes

- **Chinese characters**: Use Traditional Chinese (繁體中文) for all user-facing output
- **lunar-zenith**: External library for accurate stem-branch calculations
- **Go version**: Requires Go 1.25+
- **No type suppression**: Never use `as any`, `@ts-ignore`, or similar
- **Test-first**: Write tests for new algorithms before implementation (TDD preferred)
