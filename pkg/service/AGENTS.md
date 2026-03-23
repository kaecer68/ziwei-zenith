# pkg/service/ - 服務層

## 概述

Service package 提供統一的計算接口和 gRPC 服務實現，作為 engine 與外部協議之間的適配層。

## 結構

```
pkg/service/
├── calculate.go          # 統一計算邏輯封裝
└── grpc_server.go        # gRPC 服務實現
```

## 關鍵接口

### Calculate
```go
func Calculate(input CalculateInput) (*engine.ZiweiChart, error)
```

統一入口，處理：
- 時區轉換 (經度計算)
- 夏令時調整
- 陰曆/陽曆轉換
- 四柱計算

### gRPC Server
```go
type ZiweiServer struct {
    ziwei.UnimplementedZiweiServiceServer
}
```

實現服務：
- `Calculate` - 基本排盤
- `CalculateTemporal` - 流運計算
- `ListRecords` - 記錄列表
- `ListTags` - 標籤列表

## 數據流

```
REST/gRPC Request
       ↓
  grpc_server.go (協議轉換)
       ↓
  calculate.go (參數處理)
       ↓
  pkg/engine (核心計算)
       ↓
  Response
```

## 依賴關係

- **依賴**: `pkg/engine`, `pkg/basis`, `lunar-zenith`
- **被調用**: `cmd/ziwei-server`

## 契約約束

所有 API 欄位必須符合：
```
destiny-contracts/openapi/ziwei-zenith.yaml
```

## 注意事項

- 經度預設值：121.565 (台北)
- 性別字符串映射：`"male"` → `basis.SexMale`
- 流運字段語義：`stem` + `time_branch` = 時間干支

## 數據轉換規範

### 星曜亮度轉換

當需要將引擎計算結果轉換為 API 響應時，必須根據星曜類型選擇正確的亮度查詢函數：

```go
// 主星亮度
brightness := basis.BrightnessLevel(star, branch)

// 六吉星亮度
brightness := basis.AuspiciousBrightnessLevel(star, branch)

// 祿存、天馬亮度
brightness := basis.LuCunBrightnessLevel(star, branch)

// 六煞星亮度
brightness := basis.MaleficBrightnessLevel(star, branch)
```

### 類型斷言處理

`chart.AssistantStars` 存儲的是 `[]interface{}`，需要進行類型斷言：

```go
for _, s := range chart.AssistantStars[b] {
    switch star := s.(type) {
    case basis.AuspiciousStar:
        // 六吉星：左輔、右弼、文昌、文曲、天魁、天鉞
        brightness = basis.AuspiciousBrightnessLevel(star, b)
    case basis.LuCunStar:
        // 祿存、天馬
        brightness = basis.LuCunBrightnessLevel(star, b)
    case basis.MaleficStar:
        // 六煞星：擎羊、陀羅、火星、鈴星、地空、地劫
        brightness = basis.MaleficBrightnessLevel(star, b)
    }
}
```

## 禁止事項

- 禁止在此處實現業務算法（應調用 engine）
- 禁止修改 API 契約未定義的欄位
