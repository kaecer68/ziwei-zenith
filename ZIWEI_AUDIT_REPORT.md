# 紫微斗數算法稽核報告

## 稽核日期
2026-03-23

## 稽核範圍
- pkg/engine/lifepalace.go - 命宮身宮計算
- pkg/engine/starplacement.go - 主星配置
- pkg/engine/dayun.go - 大運流年計算
- pkg/basis/transformation.go - 四化對照表
- pkg/basis/wuxing.go - 五行局

## 研究方法
1. 網路搜尋紫微斗數標準算法規則
2. 對照古籍安星訣口訣
3. 實際案例驗證
4. 執行單元測試

---

## 各項算法稽核結果

### 1. 命宮計算 ✅ 正確

**規則** (安命身宮訣):
> 寅起正月，順數至生月，逆數生時為命宮

**代碼實現** (lifepalace.go:12-26):
```go
monthPos := (2 + lunarMonth - 1) % 12
mingIdx := (monthPos - int(hourBranch) + 12) % 12
```

**驗證結果**:
- 農曆5月巳時 → 命宮在丑(1) ✓
- 測試已通過

---

### 2. 身宮計算 ✅ 正確

**規則**:
> 寅起正月，順數至生月，順數生時為身宮

**代碼實現**:
```go
shenIdx := (monthPos + int(hourBranch)) % 12
```

**驗證結果**:
- 農曆5月巳時 → 身宮在亥(11) ✓
- 計算邏輯正確

---

### 3. 大運計算 ✅ 正確

**規則** (起大限訣):
> 大限由命宮起，陽男陰女順行；陰男陽女逆行，每十年過一宮限

**五行局起運歲數**:
- 水二局：2歲起運
- 木三局：3歲起運
- 金四局：4歲起運
- 土五局：5歲起運
- 火六局：6歲起運

**代碼實現** (dayun.go:7-36):
```go
isYang := (int(yearStem) % 2) == 0
direction := 1
if (sex == basis.SexMale && !isYang) || (sex == basis.SexFemale && isYang) {
    direction = -1
}
```

**驗證結果**:
- 陽男(甲年) → 順行 ✓
- 陽女(甲年) → 逆行 ✓
- 陰男(乙年) → 逆行 ✓
- 陰女(乙年) → 順行 ✓

---

### 4. 四化對照表 ✅ 正確

**口訣**:
> 甲廉破武陽，乙機梁紫陰，丙同機昌廉，丁陰同機巨，戊貪陰右機，己武貪梁曲，庚陽武府同，辛巨陽曲昌，壬梁紫左武，癸破巨陰貪

**代碼實現** (transformation.go:34-44):
```go
var TransformationTable = map[Stem][4]string{
    StemJia:  {"廉貞", "破軍", "武曲", "太陽"},  // 廉破武陽
    StemYi:   {"天機", "天梁", "紫微", "太陰"},  // 機梁紫陰
    StemBing: {"天同", "天機", "文昌", "廉貞"},  // 同機昌廉
    StemDing: {"太陰", "天同", "天機", "巨門"},  // 陰同機巨
    StemWu:   {"貪狼", "太陰", "右弼", "天機"},  // 貪陰右機
    StemJi:   {"武曲", "貪狼", "天梁", "文曲"},  // 武貪梁曲
    StemGeng: {"太陽", "武曲", "天府", "天同"},  // 陽武府同
    StemXin:  {"巨門", "太陽", "文曲", "文昌"},  // 巨陽曲昌
    StemRen:  {"天梁", "紫微", "左輔", "武曲"},  // 梁紫左武
    StemGui:  {"破軍", "巨門", "太陰", "貪狼"},  // 破巨陰貪
}
```

**驗證結果**: 十天干四化全部正確 ✓

---

### 5. 流年計算 ✅ 正確

**規則**:
> 流年命宮 = 流年地支所在宮位

**代碼實現** (dayun.go:40-50):
```go
func CalcLiuNian(yearBranch basis.Branch, currentYear int) basis.LiuNian {
    return basis.LiuNian{
        Year:   currentYear,
        Branch: yearBranch,
    }
}
```

**驗證結果**:
- 2024甲辰年 → 辰宮 ✓
- 2025乙巳年 → 巳宮 ✓

---

### 6. 流月計算 ✅ 正確

**規則** (安子斗訣 / 斗君):
> 流年歲建起正月，逆逢生月順回程，回程順至生時止，便是流年正月春

**代碼實現** (dayun.go:53-62):
```go
func CalcLiuYue(lnBranch basis.Branch, birthMonth int, birthHour basis.Branch, targetLunarMonth int) basis.Branch {
    month1Idx := (int(lnBranch) - (birthMonth - 1) + int(birthHour) + 12) % 12
    lYueIdx := (month1Idx + (targetLunarMonth - 1)) % 12
    return basis.Branch(lYueIdx)
}
```

**驗證結果**:
- 甲辰年(辰=4), 農曆4月生, 丑時
- 正月在寅宮(2), 順行排列十二個月 ✓

---

### 7. 流日計算 ✅ 正確

**規則**:
> 以流月所在的宮起初一，順行十二宮

**代碼實現** (dayun.go:65-69):
```go
func CalcLiuRi(lyBranch basis.Branch, targetLunarDay int) basis.Branch {
    lRiIdx := (int(lyBranch) + (targetLunarDay - 1)) % 12
    return basis.Branch(lRiIdx)
}
```

**驗證結果**:
- 流月在辰宮 → 初一辰宮, 初二巳宮... ✓

---

### 8. 紫微星定位 ✅ 正確

**規則** (起紫微星訣):
> 六五四三二，酉午亥辰丑，局數除日數，商數宮前走...

**代碼實現** (starplacement.go:8-27):
實現了完整的公式計算，包括餘數處理和奇偶判斷

**驗證結果**:
- 各五行局、各日數測試已通過 ✓

---

### 9. 主星配置 ✅ 正確

**規則**:
- 紫微星系：紫微逆去天機星，隔一太陽武曲辰，連接天同空二宮，廉貞居處方是真
- 天府星系：天府順行有太陰，貪狼而後巨門臨，隨來天相天梁繼，七殺空三是破軍

**代碼實現** (starplacement.go:29-76):
- 紫微星系逆時針排列 ✓
- 天府與紫微對宮關係 ✓
- 天府星系順時針排列 ✓

---

## 測試執行結果

```
=== RUN   TestCalcLifePalace
--- PASS: TestCalcLifePalace (0.00s)
=== RUN   TestDaYunDirection
--- PASS: TestDaYunDirection (0.00s)
=== RUN   TestDaYunAges
--- PASS: TestDaYunAges (0.00s)
=== RUN   TestLiuNian
--- PASS: TestLiuNian (0.00s)
=== RUN   TestLiuYueDouJun
--- PASS: TestLiuYueDouJun (0.00s)
=== RUN   TestLiuRi
--- PASS: TestLiuRi (0.00s)
=== RUN   TestCalcZiweiStarPos
--- PASS: TestCalcZiweiStarPos (0.00s)
=== RUN   TestPlaceMainStars
--- PASS: TestPlaceMainStars (0.00s)
PASS
ok      github.com/kaecer68/ziwei-zenith/pkg/engine
```

**所有測試通過 ✓**

---

## 結論

經過深入的網路研究、口訣對照和實際驗證，**目前代碼的紫微斗數算法實現是正確的**，包括：

1. ✅ 命宮身宮計算
2. ✅ 大運順逆與起運歲數
3. ✅ 四化對照表
4. ✅ 流年命宮
5. ✅ 流月斗君
6. ✅ 流日計算
7. ✅ 紫微星定位
8. ✅ 主星配置

**無需修復，算法正確！**

---

## 參考資料

1. 紫微斗數安星訣 (iztro.com)
2. 紫微斗數全書
3. 各派排盤軟體對照驗證
