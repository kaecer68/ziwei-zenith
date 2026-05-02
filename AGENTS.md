<!-- AGENTS.md - Ziwei Zenith Development Guide -->

**專案**: Ziwei Zenith（紫微斗數排盤引擎）  
**更新日期**: 2026-04-12  
**分支**: main

---

## 1. 專案概述

Ziwei Zenith 是一個開源的紫微斗數排盤引擎，使用 **Go 1.25+** 實現核心計算邏輯，並提供 **REST API** 與 **gRPC** 雙協定服務。前端採用 **React 19 + TypeScript 5 + Vite 8**，提供互動式命盤顯示界面。

### 核心功能
- **星曜系統**：十四主星、六吉星、六煞星、祿存天馬、丙級星、小星神煞、十二長生、博士十二神
- **運限系統**：大運（十年）、流年、流月、流日
- **分析功能**：五行局、納音五行、星曜亮度（七級制）、格局判斷、來因宮偵測、三方四正診斷、時空感應（Temporal Resonance）

---

## 2. 技術棧與運行時架構

### 後端技術棧
| 項目 | 版本 / 套件 |
|------|------------|
| Go | 1.25.6 |
| gRPC | `google.golang.org/grpc` |
| Protobuf | `google.golang.org/protobuf` |
| 曆法計算 | `github.com/kaecer68/lunar-zenith` v0.1.1 |
| HTTP 框架 | 標準函式庫 `net/http` |

### 前端技術棧
| 項目 | 版本 / 套件 |
|------|------------|
| React | 19.2.4 |
| TypeScript | ~5.9.3 |
| Vite | 8.0.0 |
| UI 動畫 | Framer Motion 12.36.0 |
| 樣式 | TailwindCSS（透過 CDN 或專案配置） |
| HTTP 客戶端 | Axios 1.13.6 |
| 圖標 | Lucide React |

### 運行時架構

本專案採用**單體雙協定伺服器**設計：

```
┌─────────────────────────────────────────────────────────────┐
│                    cmd/ziwei-server                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ REST Server  │  │ gRPC Server  │  │ CORS / Health│      │
│  │   (:8083)    │  │   (:50053)   │  │   Handlers   │      │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘      │
│         │                 │                                 │
│         └────────┬────────┘                                 │
│                  ▼                                          │
│         pkg/service/calculate.go                            │
│                  ▼                                          │
│           pkg/engine/ (核心排盤算法)                         │
│                  ▼                                          │
│           pkg/basis/ (類型與常數定義)                        │
└─────────────────────────────────────────────────────────────┘
```

- **REST 與 gRPC 共用同一進程**：`cmd/ziwei-server/main.go` 同時啟動 REST (`:8083`) 與 gRPC (`:50053`) 服務。
- **狀態持久化**：紀錄 (`records.json`) 與標籤 (`tags.json`) 以本地 JSON 檔案儲存，透過 `sync.RWMutex` 保護。
- **曆法計算**：所有陰陽曆轉換、四柱計算均委託給外部函式庫 `lunar-zenith`。
- **契約驅動 Port**：所有服務埠號由 `destiny-cloud/contracts/runtime/ports.env` 統一定義，經 `make sync-contracts` 生成 `.env.ports` 後載入。**嚴禁手動修改 `.env.ports`**。

---

## 3. 專案結構與模組劃分

```
ziwei-zenith/
├── cmd/
│   ├── ziwei-cli/          # CLI 命令列工具
│   └── ziwei-server/       # REST + gRPC 雙協定伺服器入口
├── contracts/              # ← symlink → destiny-cloud/contracts (single source of truth)
│   ├── openapi/            # OpenAPI 規格（如 ziwei-zenith.yaml）
│   ├── runtime/ports.env   # 單一真相：所有服務埠號定義
│   └── proto/              # 跨服務共享的 Protobuf 定義
├── proto/
│   └── ziwei.proto         # 本服務專用的 Protobuf 定義
├── pkg/
│   ├── api/
│   │   ├── v1/types.go     # REST JSON API 類型定義
│   │   └── grpc/v1/        # protoc 生成的 gRPC Go 程式碼（*.pb.go）
│   ├── basis/              # 核心定義：星曜、宮位、五行、亮度、四化等
│   ├── engine/             # 排盤計算引擎：命宮、主星、輔星、大運、格局
│   └── service/            # 服務層：統一計算入口 + gRPC 服務器實作
├── web/                    # React + TypeScript 前端
│   ├── src/
│   │   ├── App.tsx
│   │   ├── components/     # Palace chart, interpretation panel 等
│   │   └── styles/         # Tailwind / 傳統風格 CSS
│   ├── package.json
│   └── vite.config.ts      # Vite 配置（含 API Proxy）
├── scripts/
│   ├── sync-contracts.sh   # 契約同步腳本
│   ├── dev-clean.sh        # 清理被佔用的開發埠號
│   └── dev-watchdog.js     # Vite 開發伺服器看門狗
├── .github/workflows/
│   └── verify-contracts.yml # CI：驗證契約同步狀態
├── ecosystem.config.js     # PM2 進程管理配置
├── Makefile                # 常用開發命令
├── go.mod                  # Go 模組定義（Go 1.25.6）
├── .env.example            # 本地環境變數範例
└── AGENTS.md               # 本文件
```

### 模組職責

| 模組 | 職責 | 關鍵文件 |
|------|------|---------|
| `pkg/basis` | **純數據定義**。所有類型、常數、查找表（Lookup Tables）。不可包含業務邏輯。 | `stars.go`, `palaces.go`, `brightness.go`, `wuxing.go` |
| `pkg/engine` | **核心算法**。紫微斗數排盤的所有數學邏輯，必須是純函數、無副作用。 | `engine.go`, `lifepalace.go`, `starplacement.go`, `dayun.go` |
| `pkg/service` | **協議適配層**。統一計算入口（時區、夏令時、陰曆轉換）+ gRPC Server 實作。 | `calculate.go`, `grpc_server.go` |
| `pkg/api/v1` | **API 類型**。REST 響應結構定義，需與 `proto/ziwei.proto` 及契約保持一致。 | `types.go` |
| `cmd/ziwei-server` | **運行時入口**。HTTP handler、CORS、狀態管理（records/tags）、REST ↔ Engine 轉換。 | `main.go` |
| `web/` | **前端界面**。命盤渲染、互動、API 調用。 | `App.tsx`, `components/ZiweiChart.tsx` |

---

## 4. 關鍵配置文件說明

| 文件 | 用途 | 注意事項 |
|------|------|---------|
| `go.mod` | Go 模組定義。依賴 `lunar-zenith`、gRPC、Protobuf。 | Go 版本要求 1.25+ |
| `web/package.json` | 前端依賴與腳本。 | 使用 `type: "module"` |
| `web/vite.config.ts` | Vite 構建與開發伺服器配置。會讀取 `.env.ports` 中的 `VITE_API_TARGET` 來設定 proxy。 | 若 `.env.ports` 不存在會直接報錯 |
| `web/go.mod` | **假 go.mod**，用於防止 `go test ./...` 掃描 `node_modules`。 | 無實際 Go 代碼 |
| `proto/ziwei.proto` | gRPC 服務與訊息定義。 | 修改後需重新執行 `protoc` 生成 `.pb.go` |
| `ecosystem.config.js` | PM2 進程管理配置。定義 `ziwei-server`、`ziwei-web`、`ziwei-watchdog` 三個應用。 | 用於本地長期運行或部署 |
| `Makefile` | 標準化開發命令：run、sync-contracts、verify-contracts、web-dev、dev-all 等。 | 所有命令執行前會先跑 `sync-contracts` |
| `.env.example` | 環境變數範例。目前僅定義 `GRPC_PORT` 與 `REST_PORT`。 | 實際運行時優先讀 `.env.ports` 或環境變數 |
| `contracts/runtime/ports.env` | **單一真相檔**。所有服務埠號的源頭。 | 修改需透過 `destiny-cloud/contracts` |

---

## 5. Build & Test Commands

```bash
# ─── Go 後端 ───
# 編譯所有套件
go build ./...

# 執行所有測試（自動排除 web/ 的 node_modules 因 web/go.mod 存在）
go test ./...

# 僅執行引擎層測試（最詳盡）
go test -v ./pkg/engine/...

# 程式碼檢查與格式化
go vet ./...
go fmt ./...

# 完整檢查
go build ./... && go vet ./... && go fmt ./...

# ─── 前端 ───
cd web && npm install
cd web && npm run build    # tsc -b && vite build
cd web && npm run lint     # eslint .

# ─── 契約同步 ───
make sync-contracts        # 開發前必跑：生成 .env.ports
make verify-contracts      # 提交/CI 必跑：驗證 .env.ports 與契約一致
make dev-clean             # 若 port 被佔用，先清理再重跑
```

### 測試策略說明
- **核心測試**：集中於 `pkg/engine/`（約 9 個 `*_test.go` 檔案），覆蓋命宮計算、主星配置、輔星、大運、格局、解盤等關鍵算法。
- **運行時測試**：`cmd/ziwei-server/main_test.go` 測試環境變數與 `.env.ports` 的載入邏輯。
- **CI Gate**：`.github/workflows/verify-contracts.yml` 會檢查契約同步狀態；未同步會導致 CI fail。
- **前端**：目前無自動化測試框架配置，以手動驗證與 `eslint` 靜態檢查為主。

---

## 6. 開發與部署流程

### 本地開發步驟

1. **安裝依賴**
   ```bash
   go mod download
   cd web && npm install
   ```

2. **契約同步**（每次開發前必做）
   ```bash
   make sync-contracts
   ```

3. **啟動後端**
   ```bash
   make run
   # 或
   go run ./cmd/ziwei-server/main.go
   ```

4. **啟動前端**
   ```bash
   make web-dev          # 一般模式
   make web-dev-safe     # 帶看門狗（推薦長期開發）
   make dev-all          # 同時啟動後端 + 前端看門狗
   ```

### REST API 端點

| 路徑 | 方法 | 說明 |
|------|------|------|
| `/health` | GET | 健康檢查（新版） |
| `/api/v1/health` | GET | 健康檢查（向下兼容） |
| `/api/v1/calculate` | POST | 基本排盤 |
| `/v1/ziwei/calculate` | POST | 基本排盤（舊版兼容） |
| `/api/v1/calculate/temporal` | POST | 動態運限計算（大限/流年/流月/流日） |
| `/api/v1/records` | GET/POST | 紀錄列表 / 新增 |
| `/api/v1/records/{id}` | PUT/DELETE | 更新 / 刪除紀錄 |
| `/api/v1/tags` | GET/PUT | 標籤列表 / 更新 |

### 部署方式

- **PM2 部署**：使用根目錄 `ecosystem.config.js` 管理進程，包含後端伺服器、Vite 開發伺服器、以及看門狗進程。
- **無 Dockerfile**：目前專案根目錄未提供 Docker 配置，部署以原生二進位檔或 PM2 為主。
- **前端生產構建**：`make web-build` 會在 `web/dist/` 產出靜態檔案。

---

## 7. 程式碼風格規範

### Go 風格

| 元素 | 命名規範 | 範例 |
|------|---------|------|
| 類型 | PascalCase | `ZiweiEngine`, `Star`, `Palace` |
| 函數 | PascalCase（匯出）/ camelCase（內部） | `CalcLifePalace()`, `placeMainStars()` |
| 變數 | camelCase | `lifePalace`, `starList` |
| 常數 | PascalCase | `StarZiwei`, `SexMale` |
| 套件 | lowercase | `basis`, `engine`, `service` |

### 編碼約定
- **iota 枚舉**：所有類似枚舉的常數使用 `iota` 定義，並附帶繁體中文註解。
  ```go
  const (
      StarZiwei Star = iota // 紫微
      StarTianfu            // 天府
  )
  ```
- **String() 方法**：所有 `basis` 類型都必須實作 `String() string`，回傳繁體中文。
- **錯誤處理**：回傳 `error` 作為最後一個返回值，使用 `fmt.Errorf("...: %w", err)` 包裝。
- **匯入分組**：標準函式庫在前，外部套件在後，中間空一行。禁止未使用匯入。
- **顯式轉型**：使用 `basis.Branch(value)` 等顯式轉型，避免 `any` / `interface{}`。
- **Map 查找**：存取 map 前檢查 `ok`。
  ```go
  if v, ok := m[key]; ok { ... }
  ```

### TypeScript / React 風格
- **嚴禁型別壓制**：禁止 `as any`、禁止 `@ts-ignore`、禁止 `// @ts-expect-error`。
- **組件職責分離**：UI 組件中禁止混入後端計算邏輯。
- **Tab 設計**：使用圓角方形（`border-radius: 0.5rem`），文字不換行（`white-space: nowrap`）。

---

## 8. 契約優先開發流程（Contract-First）

> ⚠️ **這是本项目最重要的開發原則。**

### 規則
- **契約是唯一真相來源**：所有 API 欄位與類型必須對齊 `contracts/openapi/ziwei-zenith.yaml`。
- **先改契約，再改代碼**。
- **禁止**添加契約未定義的欄位。
- **禁止**修改契約已定義的類型。
- `.env.ports` 由 `scripts/sync-contracts.sh` 生成，**嚴禁手動編輯**。

### AI 開發流程
1. 讀取 `contracts/TASK-BOARD.md` 了解當前任務。
2. 讀取契約文件確認欄位定義。
3. 修改 `proto/ziwei.proto` 後執行 `protoc` 重新生成 gRPC 代碼。
4. 實現業務邏輯（`pkg/engine/` 或 `pkg/service/`）。
5. 執行 `make verify-contracts` 驗證。
6. 更新 `contracts/HANDOFF.md` 回報結果。

### 完成檢查清單
```markdown
- [ ] 已讀取最新契約文件
- [ ] 已運行 make sync-contracts
- [ ] 已運行 make verify-contracts（或 openapi-generator validate）
- [ ] 單元測試通過（go test ./...）
- [ ] 新增欄位已出現在契約中
- [ ] API 響應範例與契約一致
- [ ] REST & gRPC 轉換邏輯已同步更新
- [ ] 已更新 HANDOFF.md
```

---

## 9. 功能實作工作流程

當需要新增功能或 API 欄位時，請依以下順序進行：

1. **新增類型與常數** → `pkg/basis/`（定義、查找表）
2. **新增算法** → `pkg/engine/`（計算邏輯）
3. **更新核心結構體** → `pkg/engine/engine.go`（`ZiweiChart` 等）
4. **更新輸出格式化** → `engine.go` 的 `String()` 方法
5. **更新協議定義** → `proto/ziwei.proto`，然後重新生成 `pkg/api/grpc/v1/*.pb.go`
6. **更新 REST 類型** → `pkg/api/v1/types.go`
7. **更新數據轉換**（必須同步兩處）：
   - `pkg/service/grpc_server.go`（gRPC 轉換）
   - `cmd/ziwei-server/main.go`（REST 轉換，函數 `mapChartToResponse`）
8. **測試驗證**：
   - CLI：`go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10`
   - Server：`make run` 後驗證 REST + gRPC

---

## 10. 重要實踐經驗與禁止事項

### ANTI-PATTERNS（絕對禁止）

| 禁止項目 | 說明 |
|---------|------|
| **型別壓制** | 禁止 `as any`、禁止 `@ts-ignore`、禁止無理由的 `interface{}` |
| **違反契約** | 禁止添加契約未定義的 API 欄位 |
| **跨專案污染** | 禁止將 `lunar-zenith`、`bazi-zenith` 的邏輯混入本專案 |
| **手改生成檔** | 禁止手動編輯 `.pb.go` 檔案，必須使用 `protoc` 重新生成 |
| **未使用匯入** | Go 編譯會因此失敗，請保持匯入清潔 |
| **手改 .env.ports** | 任何埠號變更必須從 `destiny-cloud/contracts/runtime/ports.env` 同步 |

### 專案隔離原則
- **純離線可用**：`pkg/` 目錄下不得進行任何外部網路連線。
- **純邏輯層**：`pkg/engine/` 必須保持純函數，無副作用、無 I/O。
- **獨立邊界**：所有紫微斗數算法必須限制在 `ziwei-zenith` 倉庫內。

### 關鍵 UI 規範（後端需配合輸出正確數據）

#### 1. 宮位變化映射公式
```typescript
// ✅ 正確：從新基準宮位看當前宮位
const shifted = palaceCycle[(baseIndex - sourceIndex + 12) % 12];

// ❌ 錯誤：方向相反
const shifted = palaceCycle[(sourceIndex - baseIndex + 12) % 12];
```

#### 2. 星曜亮度規範
- **所有主要星曜都有亮度**：十四主星、六吉星、祿存天馬、六煞星。
- 後端通過 `assistant_star_details` 欄位傳遞輔星與煞星的亮度資訊。
- 亮度查詢函數：
  - 主星：`basis.BrightnessLevel(star, branch)`
  - 六吉星：`basis.AuspiciousBrightnessLevel(star, branch)`
  - 祿存天馬：`basis.LuCunBrightnessLevel(star, branch)`
  - 六煞星：`basis.MaleficBrightnessLevel(star, branch)`

#### 3. 星曜層級顯示
- **六吉星**（左輔、右弼、文昌、文曲、天魁、天鉞）必須單獨列顯示，與主星同級，不可與火星、鈴星等雜星混排。

#### 4. 四化顯示規範
- **生年四化**（紅色）：本命盤固定不變
- **本宮飛化**（金色）：該宮天干對本宮星曜
- **選中宮位飛化**（紫色）：動態計算，影響所有宮位

---

## 11. 安全與隱私考量

- **資料儲存**：紀錄與標籤僅儲存在本地 JSON 檔案（`records.json`、`tags.json`），無資料庫相依。
- **CORS 設定**：REST Server 目前設定 `Access-Control-Allow-Origin: *`，允許所有來源。若部署至生產環境，建議根據前端域名縮小範圍。
- **無認證層**：本服務不處理身份驗證，預設由上游 Gateway 或 destiny-cloud 負責。
- **Fail Fast**：若缺少 `.env.ports` 或必要 runtime 環境變數，伺服器會立即 `panic` 並提示錯誤訊息。

---

## 12. 子目錄專屬指南

更詳細的套件級指引請參閱：

- `pkg/basis/AGENTS.md` — 核心定義、類型設計、查找表規範
- `pkg/engine/AGENTS.md` — 計算算法、測試策略、純函數原則
- `pkg/service/AGENTS.md` — gRPC 服務實作、數據轉換規範、亮度查詢範例
- `web/AGENTS.md` — React 組件結構、UI/UX 規範、命盤顯示細則
- `.ziwei-skills/` — 紫微斗數領域知識文件（如宮位映射、亮度系統、四化飛星、三方四正）

---

## 13. 常用命令速查表

```bash
# 後端
make run                    # 契約同步 + 啟動 REST/gRPC 伺服器
go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10 -gender female -json
go test -v ./pkg/engine/...

# 前端
make web-dev                # 啟動 Vite 開發伺服器
make web-dev-safe           # 啟動 Vite + 看門狗
make web-build              # 生產構建

# 契約與環境
make sync-contracts         # 開發前必跑
make verify-contracts       # 提交前必跑
make dev-clean              # 清理佔用埠號
make dev-all                # 同時啟動後端 + 前端看門狗

# gRPC 測試（需 grpcurl）
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext -d '{"year":1972,"month":6,"day":8,"hour":1,"gender":"male"}' \
  localhost:50053 ziwei.v1.ZiweiService/Calculate
```

## graphify

This project has a graphify knowledge graph at graphify-out/.

Rules:
- Before answering architecture or codebase questions, read graphify-out/GRAPH_REPORT.md for god nodes and community structure
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files
- After modifying code files in this session, run `graphify update .` to keep the graph current (AST-only, no API cost)
