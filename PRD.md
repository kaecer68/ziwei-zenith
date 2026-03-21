# Ziwei Zenith - 紫微斗數排盤引擎

## 1. Project Overview

**Project Name**: Ziwei Zenith  
**Project Type**: Go Library + CLI + REST API + gRPC Service  
**Core Functionality**: 紫微斗數 (Ziwei Douju) fate calculation engine - a complete Chinese astrology chart calculation system  
**Target Users**: Chinese astrology practitioners, developers building fortune-telling applications, API consumers

### Maintenance Snapshot (2026-03)
- Backend quality checks: `go test ./...`, `go build ./...`, `go vet ./...` pass.
- OpenAPI validation: `contracts/openapi/ziwei-zenith.yaml` passes validation.
- Frontend build: `web` package `npm run build` passes.
- Items in **Section 8 Future Enhancements** are product roadmap tasks, not current defects.

## 2. Technical Architecture

### Technology Stack
- **Language**: Go 1.25+
- **Core Dependency**: lunar-zenith (for accurate stem-branch calculations)
- **Frontend**: React + TypeScript + Vite + Framer Motion (Pro Max Refined)
- **UI/UX Strategy**: **Input-First Architecture** (Landing directly on the Archive/Input Hub)
- **Logic Baseline**: **K-Master Engine v4.0** (Aligns with verified JavaScript implementation)
- **Development Standards**: 
    - **Plugin**: **oh-my-opencode** (Sisyphus-based agentic workflow)
    - **Skillset**: **Superpowers** (Workflow Acceleration: TDD, Brainstorming, Writing-Plans)
    - **Methodology**: PRD-Driven Execution, Atomic Updates (PRD + Implementation + Skills)
- **Security & Isolation**:
    - **Hermetic Build**: 專案為物理隔離環境設計，禁止在 `pkg/` 中引入非授權的外部連線。
    - **Independent Skills**: 雖然與 `lunar-zenith`, `bazi-zenith` 有輸入輸出連結，但所有紫微專屬邏輯（安星、四化等）必須嚴格封裝在 `ziwei-zenith` 目錄內，禁止與姊妹專案混淆。
    - **Memory Pollution Prevention**: 定期檢查並清除 AI 產生的跨專案誤導信息，確保邏輯純淨。

### Project Structure
```
ziwei-zenith/
├── cmd/
│   ├── ziwei-cli/main.go       # CLI tool
│   └── ziwei-server/main.go    # REST + gRPC server
├── proto/
│   └── ziwei.proto             # Protobuf service definition
├── pkg/
│   ├── api/
│   │   ├── v1/types.go         # REST JSON API types
│   │   └── grpc/v1/            # Generated gRPC Go code
│   ├── basis/                  # Core data types & definitions
│   │   ├── definitions.go      # Stem, Branch, Element, Pillar
│   │   ├── stars.go            # 14 main stars (十四主星)
│   │   ├── palaces.go         # 12 palaces (十二宮)
│   │   ├── wuxing.go          # Five elements (五行), NaYin (納音)
│   │   ├── auspicious.go      # Auspicious stars (六吉星)
│   │   ├── malefic.go         # Malefic stars (六煢星)
│   │   ├── secondary.go       # Secondary stars (丙級星 + 小星神煢)
│   │   ├── transformation.go # Transformation stars (四化飛星)
│   │   ├── lifecycle.go       # 十二長生 & 博士十二神
│   │   ├── dayun.go           # DaYun/LiuNian/LiuYue/LiuRi types
│   │   ├── pattern.go         # Pattern definitions
│   │   └── brightness.go      # Star brightness (7-level system)
│   ├── engine/                 # Calculation engine
│   │   ├── engine.go          # Main ZiweiEngine
│   │   ├── lifepalace.go      # Life palace calculation
│   │   ├── starplacement.go   # Main star placement algorithm
│   │   ├── assistant.go       # Assistant + secondary stars placement
│   │   ├── lifecycle.go       # 十二長生 & 博士十二神 algorithms
│   │   ├── dayun.go          # DaYun/LiuNian/LiuYue/LiuRi algorithms
│   │   ├── pattern.go        # Pattern detection
│   │   └── brightness.go     # Brightness calculation
│   └── service/                # Shared service layer
│       ├── calculate.go       # Unified calculation logic
│       └── grpc_server.go     # gRPC server implementation
├── web/                        # React + TypeScript frontend
└── go.mod
```

## 3. Implemented Features

### 3.1 Core Calculation
- ✅ **十四主星 (14 Main Stars)**: 紫微、天府、天機、太陽、武曲、天同、廉貞、巨門、天相、天梁、七殺、破軍、太陰、貪狼
- ✅ **六吉星 (6 Auspicious Stars)**: 左輔、右弼、文昌、文曲、天魁、天鉅
- ✅ **六煢星 (6 Malefic Stars)**: 擎羊、陀羅、火星、鈴星、地空、地劫
- ✅ **祿存、天馬 (LuCun & Tianma)**
- ✅ **丙級星 (Secondary Stars)**: 紅鹎、天喜、孤辰、寡宿、龍池、鳳閣、天刑、天姚、解神、天巫、三台、八座、咸池、天月、陰煢、台輔、封誥
- ✅ **小星神煢 (Minor Stars)**: 恩光、天貴、天官、天福、天哭、天虛、蓑廉、破碎、華蓋、天才、天壽、天德、月德、天傷、天使、天空、劫煢
- ✅ **十二長生 (12 Life Stages)**: 長生、沐浴、冠帶、臨官、帝旺、衰、病、死、墓、絕、胎、養
- ✅ **博士十二神 (12 PhD Stars)**: 博士、力士、青龍、小耗、將軍、奎書、飛廉、喜神、病符、大耗、伏兵、官府
- ✅ **四化飛星 (Transformation Stars)**: 化祿、化權、化科、化忌

### 3.2 Luck Cycles
- ✅ **大運 (DaYun)**: 10-year luck periods
- ✅ **流年 (LiuNian)**: Yearly luck (Fixed at Birth Year Branch)
- ✅ **流月 (LiuYue)**: Monthly luck (Dou Jun method)
- ✅ **流日 (LiuRi)**: Daily luck (Branch-step method)
- ✅ **疊併感應 (Resonance)**: Detecting clashing 'Ji' or 'Lu' between cycles.

### 3.3 Chart Analysis & Interpretation
- ✅ **十四主星 (14 Main Stars)**: 含廟旺利陷判斷
- ✅ **格局判斷 (Pattern Detection)**: 12+ 種經典格局（機月同梁、石中隱玉等）
- ✅ **深度解盤 (Deep Interpretation)**:
    - **祿隨忌走 (Karmic Flow)**: 因果鏈條敘事
    - **多校並行 (Multi-School)**: 三合、四化、欽天三位一體解讀
    - **來因宮飛化 (Origin Fly-Hua)**: 靈魂發射源分析
    - **星曜演化 (Star Evolution)**: 文化溯源與修煉建議
- ✅ **時空感應 (Temporal Resonance)**: 疊併感應警告

### 3.4 Service Layer
- ✅ **REST API** (port 8083): Full CRUD for records, tags, and chart calculation
- ✅ **gRPC Service** (port 50053): Protobuf-defined service with 5 RPCs (Calculate, ListRecords, CreateRecord, DeleteRecord, ListTags)
- ✅ **gRPC Reflection**: 支援 grpcurl 動態探索
- ✅ **Shared Service Layer**: `pkg/service/` 統一封裝計算邏輯，REST 與 gRPC 共用

### 3.5 Input/Output
- ✅ **Solar Date Calculation**: Convert Gregorian date to lunar date using lunar-zenith
- ✅ **Time System**: 支持時辰起盤 (12時辰制)
- ✅ **Gender Support**: 男命/女命 (配合 乾造/坤造 專業標識)
- ✅ **DST Intelligence**: 內建 1945-1991 台灣/大陸夏令時歷史指南與自動修正邏輯
- ✅ **UI Input Hub**: 高密度單行輸入設計，微型標籤化 (Micro-labels)，極致縮減輸入層級
- ✅ **Record Management**: 命盤紀錄庫與分類標籤（家人/親戚/朋友/同事/客戶）篩選
- ✅ **Text Output**: Human-readable console output
- ✅ **JSON Output**: `-json` flag for API integration

## 4. Algorithm Details

### 4.1 Life Palace (命宮) & Body Palace (身宮)
- **命宮**: 從寅宮起正月，順數生月，再從該宮起子時，逆數生時。
- **身宮**: 從寅宮起正月，順數生月，再從該宮起子時，順數生時。

- **紫微星定位**: 根據五行局與農曆生日，沿用 Master Engine 的 Parity 算法 (Odd Forward/Backward 邏輯)。
- **來因宮 (Origin Palace)**: 採用欽天門算法，定位生年天干所在宮位。
- **閏月處理**: 遵循 15 日分隔線法則 (Master Engine 規範)。
- **主星分佈**: 
  - 紫微星群：紫微、天機、太陽、武曲、天同、廉貞（逆時鐘分佈）。
  - 天府星群：天府、太陰、貪狼、巨門、天相、天梁、七殺、破軍（順時鐘分佈）。
  - 對稱規則：tf = (4 - zw + 12) % 12。

- **安四化規則**:
  - 甲：廉破武陽 (廉貞化祿、破軍化權、武曲化科、太陽化忌)
  - 乙：機梁紫陰 (天機化祿、天梁化權、紫微化科、太陰化忌)
  - 丙：同機昌廉 (天同化祿、天機化權、文昌化科、廉貞化忌)
  - 丁：陰同機巨 (太陰化祿、天同化權、天機化科、巨門化忌)
  - 戊：貪陰右機 (貪狼化祿、太陰化權、右弼化科、天機化忌)
  - 己：武貪梁曲 (武曲化祿、貪狼化權、天梁化科、文曲化忌)
  - 庚：陽武府同 (太陽化祿、武曲化權、天府化科、天同化忌)
  - 辛：巨陽曲昌 (巨門化祿、太陽化權、文曲化科、文昌化忌)
  - 壬：梁紫左武 (天梁化祿、紫微化權、左輔化科、武曲化忌)
  - 癸：破巨陰貪 (破軍化祿、巨門化權、太陰化科、貪狼化忌)

- **大運起點**: 五行局數（如水二局則自2歲起大運）。
- **行運方向**: 陽男陰女順行，陰男陽女逆行。
- **流年**: 固定在地盤各宮位（如甲子年流年命宮在子）。
- **流月**: 採用「斗君起月」法。自流年命宮起正月，逆數生月，再順數生時（正月點），隨後依農曆月份順行。
- **時空感應**: 自動檢測「歲運疊併」。當流年/月/日三層四化與本命四化產生重疊（如雙重化忌）時，觸發警告邏輯。

### 4.5 Star Brightness (7-Level System)
| 等級 | Description |
|------|-------------|
| 廟 (Ming) | Highest influence (宮位與星曜五行極其契合) |
| 旺 (Wang) | Strong influence |
| 得 (De) | Good influence (次旺，能量充泛) |
| 利 (Li) | Moderate influence |
| 平 (Ping) | Neutral influence |
| 陷 (Xian) | Low influence (能量受挫) |
| 不 (Bu/None)| Weakest/Disharmonious |

## 5. Usage

### REST API Server

> REST 路徑約定：主路徑為 `/api/v1/*`；`/v1/ziwei/calculate` 僅作契約中的舊版兼容標記。

```bash
# Start server (REST :8083 + gRPC :50053)
go run ./cmd/ziwei-server/main.go

# REST: Calculate chart
curl -X POST -H "Content-Type: application/json" \
  -d '{"year":1972, "month":6, "day":8, "hour":2, "minute":0, "gender":"male", "is_lunar":false, "is_leap":false, "is_dst":false, "longitude":121.565}' \
  http://localhost:8083/api/v1/calculate

# REST: List records
curl http://localhost:8083/api/v1/records

# REST: List tags
curl http://localhost:8083/api/v1/tags
```

### gRPC Service
```bash
# List services
grpcurl -plaintext localhost:50053 list

# Calculate chart
grpcurl -plaintext -d '{"year":1972,"month":6,"day":8,"hour":1,"gender":"male"}' \
  localhost:50053 ziwei.v1.ZiweiService/Calculate

# List records
grpcurl -plaintext -d '{}' localhost:50053 ziwei.v1.ZiweiService/ListRecords

# List tags
grpcurl -plaintext -d '{}' localhost:50053 ziwei.v1.ZiweiService/ListTags
```

### Library
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

## 6. Data Models

### Core Types
- `Star`: 14 main stars + assistant + secondary + transformed
- `Palace`: 12 palaces (命宮、兄弟宮、夫妻宮、子女宮、財帛宮、疾厄宮、遷移宮、僕役宮、官祿宮、田宅宮、父母宮、福德宮)
- `Branch`: 12 earthly branches (子丑寅卯辰巳午未申酉戌亥)
- `Stem`: 10 heavenly stems (甲乙丙丁戊己庚辛壬癸)
- `Wuxing`: 五行 (木火土金水)
- `NaYin`: 納音五行


### Output Structure
```json
{
  "yearPillar": "庚午",
  "interpretation": {
    "summary": "...",
    "karmic_narrative": [...],
    "san_fang_diagnosis": [...],
    "star_details": [...],
    "origin_fly_hua": {...},
    "classic_patterns": [...]
  }
}
```

## 7. Acceptance Criteria (v0.1.0)
- [x] 輸入：農曆出生年、月、日、時 → 輸出：十四主星分布
- [x] 命宮、身宮定位正確
- [x] 十二宮位排列正確
- [x] 十四主星安星法正確實現
- [x] JSON API 可用 (cmd/ziwei-server)
- [x] CLI 工具可執行排盤 (cmd/ziwei-cli)
- [x] 單元測試覆蓋核心算法 (TDD Implementation Complete)
- [x] 深度解盤測試驗證通過
- [x] 安全性審計完成 (隔離性與無污染驗證)
- [x] **UI/UX Pro Max 重構完成**: 高密度輸入、玻璃擬態視覺、DST 精準修正引導
- [x] **落地策略優化**: App 啟動默認進入「檔案館/輸入頁」，實現「即刻輸入，即刻排盤」

## 8. Future Enhancements

> 說明：以下為規劃中能力，屬於增量開發範圍，不代表當前版本功能缺陷。

### Phase 2 (Planned)
- [ ] 小限 (Xiao Xian) calculation
- [ ] 鐵板神數 (Tieban Shenshu) integration
- [ ] 更多格局判斷
- [ ] 星曜組合分析
- [ ] gRPC-Gateway (REST 自動代理 gRPC)

### Phase 3 (Roadmap)
- [ ] Mobile SDK
- [ ] Birth chart visualization
- [ ] Multi-tenant SaaS deployment

## 9. Dependencies

- **lunar-zenith**: High-precision Chinese lunar calendar library
  - Provides accurate 天干地支 (Stem-Branch) calculations
  - Handles 閏月, 節氣, and exact solar-lunar conversion
- **google.golang.org/grpc**: gRPC framework for Go
- **google.golang.org/protobuf**: Protocol Buffers for Go

## 10. Calibration History

以下校正均基於專業盤面對標驗證（參考命例：1972-06-08 01:00 男命）：

| 項目 | 修正內容 |
|------|----------|
| 天傷/天使 | 落宮偏移修正（offset -8→-7, -6→-5） |
| 蓑廉 | 查表改公式 `(yearBranch+8)%12` |
| 破碎 | 修正三組映射表（子午卯酉→巳、丑辰未戌→丑、寅巳申亥→酉） |
| 亮度系統 | 6級→7級，新增「得」等級，重寫 14 主星完整亮度表 |
| 日柱計算 | 修正 UTC 偏移問題，改用本地時間 JD 計算日柱 |
| 新增星曜 | 天空、劫煢、恩光、天貴、天官、天福、天哭、天虛、華蓋、破碎、蓑廉、天才、天壽、天德、月德、天傷、天使 |
| 新增系統 | 十二長生、博士十二神 |
| 服務協定 | 新增 gRPC 服務（port 50053），與 REST 雙協定運行 |
| 前端標籤 | 記錄庫分類標籤功能完善（家人/親戚/朋友/同事/客戶） |

## 11. References

- 《紫微斗數全書》
- 《紫微斗數命理學》
- 《現代紫微斗數論命寶典》

## 12. License

MIT License
