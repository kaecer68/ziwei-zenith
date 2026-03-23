# pkg/basis/ - 核心定義與類型

## 概述

Basis package 定義紫微斗數的所有基礎類型、常數和查找表，是整個系統的數據模型層。

## 結構

```
pkg/basis/
├── definitions.go         # 干支、五行、柱定義
├── stars.go              # 十四主星定義
├── palaces.go            # 十二宮定義
├── wuxing.go             # 五行、納音
├── auspicious.go         # 六吉星
├── malefic.go            # 六煞星
├── secondary.go          # 丙級星 + 小星神煞
├── transformation.go     # 四化飛星
├── lifecycle.go          # 十二長生 & 博士十二神
├── dayun.go              # 大運類型
├── pattern.go            # 格局定義
└── brightness.go         # 星曜亮度 (7級制)
```

## 核心類型

### 基礎枚舉
- `Branch` - 地支 (0=子, 11=亥)
- `Stem` - 天干 (0=甲, 9=癸)
- `Sex` - 性別
- `HourBranch` - 時辰

### 複合類型
- `Pillar` - 干支柱 (年/月/日/時)
- `BirthInfo` - 出生信息 (輸入參數)
- `Palace` - 宮位 (十二宮名)
- `Star` - 主星
- `Brightness` - 星曜亮度

## 使用約定

### 類型轉換
```go
// 枚舉轉字符串
branch := basis.BranchZi // 0
str := branch.String()   // "子"

// 名稱查找
b, err := basis.BranchByName("寅") // 2
```

### iota 定義模式
```go
const (
    StarZiwei Star = iota   // 紫微
    StarTianfu              // 天府
    // ...
)
```

## 依賴關係

- **被依賴**: `pkg/engine`, `pkg/service`, `cmd/*`
- **無外部依賴**: 純數據定義

## 查找表

大量使用 map 進行 O(1) 查找：
- `StarBrightnessTable` - 星曜亮度表
- `NaYinTable` - 納音五行表

## 注意事項

- 所有類型必須實現 `String()` 方法
- 枚舉值從 0 開始，使用 iota
- 中文註解說明每個常數的含義

## 禁止事項

- 禁止在此處添加業務邏輯（只放定義）
- 禁止修改已定義的枚舉順序（會破壞查找表）
