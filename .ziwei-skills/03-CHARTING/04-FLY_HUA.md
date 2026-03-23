# 飛星四化計算

> **前置要求**：
> - [四化星](../02-STARS/03-TRANSFORM_STARS.md)
> - [本命盤排法](./01-BIRTH_CHART.md)
> - [三方四正](./03-SANFANG_SIZHENG.md)
>
> ⚠️ **重要**：飛星四化是紫微斗數動態論命的核心。

---

## 1. 飛星四化的定義

**飛星四化** = 宮干的四化飛入到該星曜所在的宮位

```
宮干（時間干）
    │
    ▼
查四化表（甲→廉破武陽...）
    │
    ▼
找到被激活的星曜
    │
    ▼
飛入該星曜所在的宮位
```

---

## 2. 飛星四化計算流程

### Step 1: 取得宮干

每個宮位有一個宮干（天干）。

### Step 2: 查四化表

```javascript
const transformationTable = {
  '甲': ['廉貞', '破軍', '武曲', '太陽'], // [化祿, 化權, 化科, 化忌]
  '乙': ['天機', '天梁', '紫微', '太陰'],
  // ...
};
```

### Step 3: 遍歷12宮找星曜

```javascript
const calculateFlyHua = (palaceGan, palaces) => {
  const results = [];
  const table = transformationTable[palaceGan];
  if (!table) return results;

  const transTypes = ['祿', '權', '科', '忌'];

  for (let i = 0; i < 4; i++) {
    const starName = table[i];
    const transType = transTypes[i];

    // 在12宮中找這個星曜
    for (const [palaceName, data] of Object.entries(palaces)) {
      const allStars = [...data.stars, ...data.assistant_stars];
      if (allStars.includes(starName)) {
        results.push({
          type: transType,
          star: starName,
          targetPalace: palaceName,
        });
        break; // 找到後跳出
      }
    }
  }
  return results;
};
```

---

## 3. 飛星四化的論命意義

### 3.1 本命飛星

本命盤中各宮的宮干產生的飛星四化。

### 3.2 大限飛星

大限期間，該大限宮干產生的飛星四化。

### 3.3 流年飛星

流年期間，流年宮干產生的飛星四化。

### 3.4 流月/流日飛星

以此類推...

---

## 4. 顯示格式

| 層級 | 顯示格式 | 範例 |
|------|---------|------|
| 本命飛星 | 宮名+四化名 | 命宮祿、命宮權 |
| 大限飛星 | 大限+四化名 | 大限祿、大限權 |
| 流年飛星 | 流年+四化名 | 流年祿、流年權 |
| 流月飛星 | 流月+四化名 | 流月祿、流月權 |
| 流日飛星 | 流日+四化名 | 流日祿、流日權 |

---

## 5. 代碼對照

| 實體 | 程式碼位置 | 說明 |
|------|-----------|------|
| 四化計算 | `cmd/ziwei-server/main.go` | `calculateTemporalTransforms()` |
| 四化查詢表 | `pkg/basis/transformation.go` | `TransformationTable` |
| 前端飛化 | `web/src/components/PalaceDetailView.tsx` | `calculateFlyHua()` |
| 命盤飛化 | `web/src/components/ZiweiChart.tsx` | `calculatePalaceFlyHua()` |

---

## 6. 實例演練

假設「丙寅」年出生，命宮在寅宮。

**問題**：寅宮的宮干是「丙」，那麼丙干的四化飛星是什麼？

**解答**：
1. 丙 → 四化表 → 天同化祿、天機化權、文昌化科、廉貞化忌
2. 在命盤中找到這些星曜的位置：
   - 天同星在 ○○宮 → ○○宮有「天同化祿」
   - 天機星在 ○○宮 → ○○宮有「天機化權」
   - 文昌星在 ○○宮 → ○○宮有「文昌化科」
   - 廉貞星在 ○○宮 → ○○宮有「廉貞化忌」

---

## 7. 自我檢查

1. 什麼是飛星四化？
2. 宮干和年干在四化計算上有什麼不同？
3. 假設宮干是「甲」，那個星曜被化祿？哪些被化忌？
4. 飛星四化和生年四化在論命時有什麼區別？

---

**下一章**：[大限計算](../04-TEMPORAL/01-DA_YUN.md)
