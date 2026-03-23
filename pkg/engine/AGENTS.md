# pkg/engine/ - 紫微斗數計算引擎

## 概述

Engine package 實現紫微斗數的核心排盤算法，包括主星配置、輔助星計算、大運流年等所有數學邏輯。

## 結構

```
pkg/engine/
├── engine.go              # 主入口，ZiweiChart 結構體
├── lifepalace.go          # 命宮、身宮計算
├── starplacement.go       # 十四主星配置算法
├── assistant.go           # 輔助星 + 小星神煞
├── lifecycle.go           # 十二長生 & 博士十二神
├── dayun.go               # 大運算法
├── pattern.go             # 格局判斷
├── brightness.go          # 星曜亮度計算
├── interpretation.go      # 解盤邏輯
└── knowledge_base.go      # 知識庫數據
```

## 關鍵入口

| 功能 | 函數 | 文件 |
|------|------|------|
| 排盤 | `BuildChart(birth)` | engine.go |
| 命宮計算 | `CalcLifePalace(year, month, day, hour)` | lifepalace.go |
| 主星配置 | `PlaceMainStars(lifePalaceBranch)` | starplacement.go |
| 大運計算 | `CalculateDaYun(...)` | dayun.go |

## 依賴關係

- **依賴**: `pkg/basis` (所有類型定義)
- **被依賴**: `pkg/service` (服務層調用)

## 算法約定

- 使用 `Branch` (地支) 作為宮位索引，值域 0-11
- 所有位置計算使用模 12 運算
- 性別使用 `basis.SexMale` / `basis.SexFemale`

## 測試

```bash
go test -v ./pkg/engine/...
```

主要測試文件：
- `engine_test.go` - 整體排盤
- `lifepalace_test.go` - 命宮計算
- `starplacement_test.go` - 主星配置
- `dayun_test.go` - 大運計算

## 注意事項

- 日期轉換由 `lunar-zenith` 處理，本層只接收已轉換的 `BirthInfo`
- 所有計算純函數化，無副作用
- 大文件：`assistant.go` (14KB), `dayun_verification_test.go` (11KB)

## 禁止事項

- 禁止直接操作網路或外部資源（純邏輯層）
- 禁止修改 `pkg/basis` 中的類型定義
