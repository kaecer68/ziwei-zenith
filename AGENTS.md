# AGENTS.md - Ziwei Zenith Development Guide

---

## 用戶溝通語言設定

**主要語言**: 繁體中文 (Traditional Chinese)

### 規則
- **永遠使用繁體中文回覆用戶**，除非用戶明確要求使用其他語言
- 技術術語可保留英文（如：API、gRPC、JSON），但解釋必須用中文
- 代碼註解可使用中文
- 錯誤訊息和日誌建議使用中文

---

## 📋 契約優先開發流程 (Contract-First)

### 契約是唯一真相

**所有 API 變更必須遵循**:
```
destiny-contracts/openapi/ziwei-zenith.yaml
```

**規則**:
- ✅ 先更新契約，再修改代碼
- ✅ 不允許添加契約未定義的欄位
- ✅ 不允許修改契約定義的類型
- ✅ 如發現契約有問題，先更新契約

### 契約文件位置

```
ziwei-zenith/
├── contracts/              # ← symlink 指向 destiny-contracts
│   ├── openapi/
│   │   └── ziwei-zenith.yaml   # ← 契約源
│   ├── TASK-BOARD.md       # 跨服務任務看板
│   └── HANDOFF.md          # AI 交接報告
```

### AI 任務執行流程

1. **檢查 TASK-BOARD.md** → 了解當前任務
2. **讀取契約文件** → 確認欄位定義
3. **生成代碼** → `make generate`
4. **實現業務邏輯** → `pkg/service/`（依專案實際結構）
5. **驗證** → `openapi-generator validate`
6. **填寫 HANDOFF.md** → 回報結果

### 完成檢查清單

```markdown
- [ ] 已讀取最新契約文件
- [ ] 已運行 make generate
- [ ] 已運行 openapi-generator validate
- [ ] 單元測試通過
- [ ] 新增欄位已出現在契約中
- [ ] API 響應範例與契約一致
- [ ] 已更新 HANDOFF.md
```

---

## Build & Test Commands

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run engine package tests
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

# Start server (REST :8083 + gRPC :50053)
go run ./cmd/ziwei-server/main.go

# Test REST API
curl -X POST -H "Content-Type: application/json" \
  -d '{"year":1972,"month":6,"day":8,"hour":2,"gender":"male"}' \
  http://localhost:8083/api/v1/calculate

# Test gRPC (requires grpcurl)
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext -d '{"year":1972,"month":6,"day":8,"hour":1,"gender":"male"}' \
  localhost:50053 ziwei.v1.ZiweiService/Calculate
grpcurl -plaintext -d '{}' localhost:50053 ziwei.v1.ZiweiService/ListTags

# Regenerate gRPC code (after modifying proto/ziwei.proto)
protoc --go_out=pkg/api/grpc/v1 --go_opt=paths=source_relative \
  --go-grpc_out=pkg/api/grpc/v1 --go-grpc_opt=paths=source_relative \
  -I proto proto/ziwei.proto

# Runtime Port 契約同步
# 單一真相檔：contracts/runtime/ports.env（由 destiny-contracts 維護）
make sync-contracts        # 每次開發前必跑，更新 .env.ports
make verify-contracts      # 提交/CI 必跑，未同步會 fail
make dev-clean             # 若 sync/verify 指出 port 被占用，先釋放再重跑
# 注意：.env.ports 由 scripts/sync-contracts.sh 生成，嚴禁手動修改
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
6. If adding new API fields: update **BOTH**:
   - `pkg/service/grpc_server.go` (gRPC conversion logic)
   - `cmd/ziwei-server/main.go` (REST conversion logic in `mapChartToResponse`)
7. Test with CLI: `go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10`
8. Test with server: `go run ./cmd/ziwei-server/main.go` then verify REST + gRPC

---

## Important Notes

- **Chinese characters**: Use Traditional Chinese (繁體中文) for all user-facing output
- **lunar-zenith**: External library for accurate stem-branch calculations
- **Go version**: Requires Go 1.25+
- **No type suppression**: Never use `as any`, `@ts-ignore`, or similar
- **Test-first**: Write tests for new algorithms before implementation (TDD preferred)

## ANTI-PATTERNS (THIS PROJECT)

### Critical Forbidden Patterns
- **Type suppression**: NEVER use `as any`, `@ts-ignore`, `interface{}` without justification
- **Contract violations**: NEVER add API fields not defined in `destiny-contracts/openapi/ziwei-zenith.yaml`
- **Cross-project contamination**: NEVER mix logic with `lunar-zenith` or `bazi-zenith` projects
- **Auto-generated edits**: NEVER manually edit `.pb.go` files (use `make generate` instead)
- **Unused imports**: Will cause build failure - keep imports clean

### Project Isolation Rules
- **Hermetic builds**: No external network connections from `pkg/` directory
- **Pure logic**: All紫微-specific algorithms must stay within `ziwei-zenith` boundaries  
- **Memory pollution prevention**: Regular cleanup of cross-project AI artifacts

## SUBDIRECTORY GUIDES

Detailed package-specific guidance available in:
- `pkg/basis/AGENTS.md` - Core definitions and constants
- `pkg/engine/AGENTS.md` - Calculation algorithms and business logic  
- `pkg/service/AGENTS.md` - Shared services and protocol handling
- `web/AGENTS.md` - React frontend components and styling, including **UI/UX Skills** for palace chart display

---

## LESSONS LEARNED - 重要實踐經驗

> **注意**：以下為簡要總結，詳細技能文檔請參見 `.ziwei-skills/` 目錄。

### 紫微斗數核心技能索引

修改代碼前，**必須**閱讀對應技能文件：

| 主題 | 技能文件 | 重要性 |
|------|---------|--------|
| **宮位映射計算** | `.ziwei-skills/03-CHARTING/02-PALACE_CYCLE.md` | ⭐⭐⭐ 極重要 |
| **星曜亮度系統** | `.ziwei-skills/02-STARS/04-BRIGHTNESS.md` | ⭐⭐⭐ 極重要 |
| **四化飛星** | `.ziwei-skills/03-CHARTING/04-FLY_HUA.md` | ⭐⭐⭐ 極重要 |
| **三方四正** | `.ziwei-skills/03-CHARTING/03-SANFANG_SIZHENG.md` | ⭐⭐⭐ 極重要 |
| **大限流年** | `.ziwei-skills/04-TEMPORAL/` | ⭐⭐ 重要 |

### Critical UI Patterns（不可違反）

#### 1. 宮位變化映射公式（必須正確）
```typescript
// ✅ 正確：從新基準宮位看當前宮位
const shifted = palaceCycle[(baseIndex - sourceIndex + 12) % 12];

// ❌ 錯誤：方向相反（會導致宮名映射錯誤）
const shifted = palaceCycle[(sourceIndex - baseIndex + 12) % 12];
```
> 📖 詳見：`.ziwei-skills/03-CHARTING/02-PALACE_CYCLE.md`

#### 2. 星曜亮度規範（所有星曜都有亮度）
- **十四主星**、**六吉星**、**祿存天馬**、**六煞星** 都有亮度
- 亮度必須緊跟星曜名稱顯示：`紫微星 廟`、`左輔星 旺`
- 後端通過 `assistant_star_details` 欄位傳遞亮度資訊
- 亮度查詢函數：
  - 主星：`BrightnessLevel(star, branch)`
  - 六吉星：`AuspiciousBrightnessLevel(star, branch)`
  - 祿存天馬：`LuCunBrightnessLevel(star, branch)`
  - 六煞星：`MaleficBrightnessLevel(star, branch)`
> 📖 詳見：`.ziwei-skills/02-STARS/04-BRIGHTNESS.md`

#### 3. 星曜層級與顯示
- **六吉星**（左輔、右弼、文昌、文曲、天魁、天鉞）必須單獨列顯示，與主星同級
- 不可與其他雜星（火星、鈴星等）混排成一列
> 📖 詳見：`.ziwei-skills/02-STARS/02-ASSISTANT_STARS.md`

#### 4. 四化顯示規範
- **生年四化**（紅色）：本命盤固定不變
- **本宮飛化**（金色）：該宮天干對本宮星曜
- **選中宮位飛化**（紫色）：動態計算，影響所有宮位
> 📖 詳見：`.ziwei-skills/03-CHARTING/04-FLY_HUA.md`

#### 5. API 欄位更新規範（REST & gRPC 必須同步）
新增欄位時必須同時更新：
1. `proto/ziwei.proto` - Protocol Buffer 定義
2. `pkg/api/v1/types.go` - Go 類型定義
3. `pkg/service/grpc_server.go` - gRPC 轉換邏輯
4. `cmd/ziwei-server/main.go` - REST 轉換邏輯

> 📖 詳見：`pkg/service/AGENTS.md`
  assistantStarDetails = append(assistantStarDetails, &pb.PalaceStar{
    Name: starName, Brightness: brightness,
  })
}
```

4. **REST API 處理器** (`cmd/ziwei-server/main.go`):
```go
assistantStarDetails := make([]v1.PalaceStar, 0)
for _, s := range chart.AssistantStars[b] {
  // 根據星曜類型獲取亮度（與 gRPC 相同邏輯）
  ...
}
palaces[pType.String()] = v1.PalaceData{
  ...
  AssistantStarDetails: assistantStarDetails,
}
```

**常見錯誤**：只更新 gRPC 而遺漏 REST API，導致網頁無法顯示新欄位。

參見 `web/AGENTS.md` 完整規範。
