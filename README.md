# Ziwei Zenith (紫微斗數排盤引擎)

[![Go Version](https://img.shields.io/github/go-mod/go-version/kaecer68/ziwei-zenith)](https://github.com/kaecer68/ziwei-zenith)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

紫微斗數 (Ziwei Douju) 免費開源排盤引擎，採用 Go 語言實現，支援 REST API 與 gRPC 雙協定服務調用。

## 功能特點

### 星曜系統
- **十四主星**：紫微、天府、天機、太陽、武曲、天同、廉貞、巨門、天相、天梁、七殺、破軍、太陰、貪狼
- **六吉星**：左輔、右弼、文昌、文曲、天魁、天鉅
- **六煢星**：擎羊、陀羅、火星、鈴星、地空、地劫
- **祿存、天馬**
- **丙級星**：紅鹎、天喜、孤辰、寡宿、龍池、鳳閣、天刑、天姚、解神、天巫、三台、八座、咸池、天月、陰煢、台輔、封誥
- **小星神煢**：恩光、天貴、天官、天福、天哭、天虛、蓑廉、破碎、華蓋、天才、天壽、天德、月德、天傷、天使、天空、劫煢
- **十二長生**：長生、沐浴、冠帶、臨官、帝旺、衰、病、死、墓、絕、胎、養
- **博士十二神**：博士、力士、青龍、小耗、將軍、奎書、飛廉、喜神、病符、大耗、伏兵、官府
- **四化飛星**：化祿、化權、化科、化忌

### 運限系統
- **大運** (DaYun)：十年大運
- **流年** (LiuNian)：流年運勢
- **流月** (LiuYue)：流月運勢
- **流日** (LiuRi)：流日運勢

### 分析功能
- **命宮、身宮計算**
- **五行局**（金木水火土各局）
- **納音五行**
- **星曜亮度**：採用七級制（廟、旺、得、利、平、陷、不）
- **格局判斷**：收錄三奇加會、機月同梁、殺破狼等經典格局
- **來因宮偵測**：自動識別靈魂能量來源的「來因宮」
- **能量循環分析 (祿隨忌走)**：基於四化飛星的動態因果敘事
- **三方四正專業診斷**：提供本位(體)、對宮(用)、氣數位(氣)、財帛位(數)的深層解析
- **時空感應 (Temporal Resonance)**：自動偵測流年/流月/流日與本命盤的能量疊併（疊忌、祿忌交戰）

## 安裝

```bash
# Clone the repository
git clone https://github.com/kaecer68/ziwei-zenith.git
cd ziwei-zenith

# Build
go build ./...

# Or run directly
go run ./cmd/ziwei-cli/main.go -year 1990 -month 6 -day 15 -hour 10
```

## 使用方法

### CLI 命令列工具

```bash
# 基本用法
ziwei-cli -year 1990 -month 6 -day 15 -hour 10

# 指定性別
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -gender female

# JSON 輸出
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -json
```

### Running Tests

To ensure the accuracy of the calculation engine, you can run the suite of unit tests:

```bash
# Run all project tests
go test ./...

# Run all engine tests
go test -v ./pkg/engine/...

# Run specific functional tests
go test -v pkg/engine/lifepalace_test.go pkg/engine/lifepalace.go
go test -v pkg/engine/starplacement_test.go pkg/engine/starplacement.go
```

Current status: unit tests are maintained mainly in `pkg/engine`, while other packages are validated via build/vet checks.

The test suite covers:
- **Life/Body Palace** placement rules.
- **Main Stars** (14) placement.
- **Assistant/Malefic Stars** placement.
- **NaYin & Five Elements** logic.
- **Transformation (Si Hua)** logic.

Recommended quick verification:

```bash
go build ./...
go vet ./...
```

### 其他 CLI 選項

```bash
# 指定出生地經緯度（用於計算真太陽時）
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -lat 25.033 -lon 121.565
```

### REST API + gRPC Server

> REST 路徑約定：目前以 `/api/v1/*` 為主要路徑；`/v1/ziwei/calculate` 僅保留為契約中的舊版兼容標記。
>
> REST Port 同步規則：優先使用環境變數 `REST_PORT`，未設置時會從 `contracts/openapi/ziwei-zenith.yaml` 的 `servers.url` 解析 port。

```bash
# 契約同步 + 啟動伺服器 (REST :8083 + gRPC :50053)
make dev

# REST: 計算命盤
curl -X POST -H "Content-Type: application/json" \
  -d '{"year":1972, "month":6, "day":8, "hour":2, "minute":0, "gender":"male", "is_lunar":false, "is_leap":false, "is_dst":false, "longitude":121.565}' \
  http://localhost:8083/api/v1/calculate | jq .

# REST: 列出紀錄 / 標籤
curl http://localhost:8083/api/v1/records
curl http://localhost:8083/api/v1/tags
```

### Runtime Port 契約同步

- **單一真相檔**：`contracts/runtime/ports.env`（由 destiny-contracts 維護）
- **同步腳本**：每次開發前必跑 `make sync-contracts`
- **驗證腳本**：提交前、PR/CI 會跑 `make verify-contracts`，若 `.env.ports` 未與契約同步會直接 fail
- **本地 port 衝突**：若執行 sync/verify 時提示 port 已被占用，使用 `make dev-clean` 釋放資源後再試
- **禁止手改 `.env.ports`**：此檔由 `scripts/sync-contracts.sh` 產生，請勿人工編輯；任何變更需從契約 repo 更新後再同步
- **PR/CI 驗證 Gate**：Pipeline 內建 `make verify-contracts`，未同步 `.env.ports` 會直接 fail 並阻擋合併

```bash
# 同步契約 port（必跑）
make sync-contracts

# 驗證 .env.ports 是否與契約一致（CI 亦會執行）
make verify-contracts

# 清除契約定義埠號的佔用行程
make dev-clean
```

### gRPC 服務調用

```bash
# 列出服務
grpcurl -plaintext localhost:50051 list

# 計算命盤
grpcurl -plaintext -d '{"year":1972,"month":6,"day":8,"hour":1,"gender":"male"}' \
  localhost:50051 ziwei.v1.ZiweiService/Calculate

# 列出紀錄 / 標籤
grpcurl -plaintext -d '{}' localhost:50051 ziwei.v1.ZiweiService/ListRecords
grpcurl -plaintext -d '{}' localhost:50051 ziwei.v1.ZiweiService/ListTags
```

### 作為 Library 使用

```go
package main

import (
    "fmt"
    "github.com/kaecer68/ziwei-zenith/pkg/service"
)

func main() {
    chart, _ := service.Calculate(service.CalculateInput{
        Year:   1990,
        Month:  6,
        Day:    15,
        Hour:   10,
        Gender: "male",
    })
    fmt.Println(chart)
}
```

## 專案結構

```
ziwei-zenith/
├── cmd/
│   ├── ziwei-cli/          # CLI 應用程式
│   └── ziwei-server/       # REST + gRPC 伺服器
├── proto/
│   └── ziwei.proto         # Protobuf 服務定義
├── pkg/
│   ├── api/
│   │   ├── v1/             # REST JSON API 類型定義
│   │   └── grpc/v1/        # 生成的 gRPC Go 程式碼
│   ├── basis/              # 核心定義（星曜、宮位、五行等）
│   │   ├── definitions.go  # 干支、五行、柱定義
│   │   ├── stars.go        # 主星定義
│   │   ├── palaces.go      # 十二宮定義
│   │   ├── wuxing.go       # 五行、納音
│   │   ├── auspicious.go   # 吉星
│   │   ├── malefic.go      # 煢星
│   │   ├── secondary.go    # 丙級星 + 小星神煢
│   │   ├── transformation.go # 四化飛星
│   │   ├── lifecycle.go    # 十二長生 & 博士十二神
│   │   ├── dayun.go        # 大運類型
│   │   ├── pattern.go      # 格局定義
│   │   └── brightness.go   # 星曜亮度 (7級制)
│   ├── engine/             # 計算引擎
│   │   ├── engine.go       # 主入口 + 輸出格式化
│   │   ├── lifepalace.go   # 命宮計算
│   │   ├── starplacement.go # 主星配置算法
│   │   ├── assistant.go    # 輔助星 + 小星神煢
│   │   ├── lifecycle.go    # 十二長生 & 博士十二神算法
│   │   ├── dayun.go        # 大運算法
│   │   ├── pattern.go      # 格局判斷
│   │   └── brightness.go   # 亮度計算
│   └── service/            # 共用服務層
│       ├── calculate.go    # 統一計算邏輯
│       └── grpc_server.go  # gRPC 服務實作
├── web/                    # React + TypeScript 前端
├── LICENSE
└── README.md
```

## 技術規格

- **Language**: Go 1.25+
- **Dependencies**: lunar-zenith, google.golang.org/grpc, google.golang.org/protobuf
- **Protocols**: REST (HTTP/JSON) + gRPC (Protobuf)
- **Output Formats**: Text (CLI), JSON (REST API), Protobuf (gRPC)
- **Frontend**: React + TypeScript + Vite + Framer Motion + TailwindCSS

## 安全與隔離性規範

- **物理隔離設計**：本專案為完全離線環境可用。
- **專案獨立性**：與 `lunar-zenith` (曆法) 及 `bazi-zenith` (八字) 保持原子化解耦。
- **無污染原則**：代碼中嚴格禁止混入其他非相關項目的遺留邏輯。

## 參考文獻

- 《紫微斗數全書》
- 《紫微斗數命理學》
- 《現代紫微斗數論命寶典》

## 授權條款

本專案採用 MIT License - 詳見 [LICENSE](LICENSE) 檔案

## 聯繫

如有任何問題，歡迎提交 Issue 或 Pull Request