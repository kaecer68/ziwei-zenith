# Ziwei Zenith (紫微斗數排盤引擎)

[![Go Version](https://img.shields.io/github/go-mod/go-version/kaecer68/ziwei-zenith)](https://github.com/kaecer68/ziwei-zenith)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

紫微斗數 (Ziwei Douju) 免費開源排盤引擎，採用 Go 語言實現。

## 功能特點

### 星曜系統
- **十四主星**：紫微、天府、天機、太陽、武曲、天同、廉貞、巨門、天相、天梁、七殺、破軍、太陰、貪狼
- **六吉星**：左輔、右弼、文昌、文曲、天魁、天鉞
- **六煞星**：擎羊、陀羅、火星、鈴星、地空、地劫
- **祿存、天馬**
- **丙級星**：紅鸞、天喜、孤辰、寡宿、龍池、鳳閣、天刑、天姚、解神、天巫、三台、八座
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
- **星曜亮度**：採用六級制（廟、旺、利、平、陷、不）
- **格局判斷**：收錄三奇加會、機月同梁、殺破狼等經典格局
- **來因宮偵測**：自動識別靈魂能量來源的「來因宮」
- **能量循環分析 (祿隨忌走)**：基於四化飛星的動態因果敘事
- **三方四正專業診斷**：提供本位(體)、對宮(用)、氣數位(氣)、財帛位(數)的深層解析

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

# 指定出生地經緯度
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -lat 25.033 -lon 121.565
```

### 作為 Library 使用

```go
package main

import (
    "fmt"
    "github.com/kaecer68/ziwei-zenith/pkg/engine"
)

func main() {
    e := engine.New()
    chart, _ := e.BuildChart(engine.BirthInfo{
        SolarYear:  1990,
        SolarMonth: 6,
        SolarDay:   15,
        Hour:       10,
        Gender:     "male",
    })
    fmt.Println(chart)
}
```

## 專案結構

```
ziwei-zenith/
├── cmd/ziwei-cli/          # CLI 應用程式
│   └── main.go
├── pkg/
│   ├── api/v1/             # JSON API 類型定義
│   ├── basis/              # 核心定義（星曜、宮位、五行等）
│   │   ├── definitions.go  # 干支、五行、柱定义
│   │   ├── stars.go        # 主星定義
│   │   ├── palaces.go      # 十二宮定義
│   │   ├── wuxing.go       # 五行、納音
│   │   ├── auspicious.go  # 吉星
│   │   ├── malefic.go      # 煞星
│   │   ├── secondary.go    # 丙級星
│   │   ├── transformation.go # 四化飛星
│   │   ├── dayun.go        # 大運類型
│   │   ├── pattern.go      # 格局定義
│   │   └── brightness.go   # 星曜亮度
│   └── engine/             # 計算引擎
│       ├── engine.go       # 主入口 + 輸出格式化
│       ├── lifepalace.go   # 命宮計算
│       ├── starplacement.go # 主星配置算法
│       ├── assistant.go    # 輔助星 + 四化星
│       ├── dayun.go        # 大運算法
│       ├── pattern.go      # 格局判斷
│       └── brightness.go   # 亮度計算
├── LICENSE                 # MIT License
└── README.md
```

## 技術規格

- **Language**: Go 1.25+
- **Dependencies**: lunar-zenith (高精確度農曆庫)
- **Output Formats**: Text (CLI), JSON (API)

## 參考文獻

- 《紫微斗數全書》
- 《紫微斗數命理學》
- 《現代紫微斗數論命寶典》

## 授權條款

本專案採用 MIT License - 詳見 [LICENSE](LICENSE) 檔案

## 聯繫

如有任何問題，歡迎提交 Issue 或 Pull Request