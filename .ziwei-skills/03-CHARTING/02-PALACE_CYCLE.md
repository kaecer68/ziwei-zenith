# 宮位流轉

> **前置要求**：[本命盤排法](./01-BIRTH_CHART.md)

---

## 1. 基本概念

### 命宮固定原則

> **重要**：紫微斗數以**地支固定方位**，宮名隨命宮位置而流轉。

```
傳統占星：宮位固定，星曜流轉
紫微斗數：星曜固定，宮名流轉（命宮為準）
```

### 十二宮順序（固定）

```javascript
const palaceCycle = [
  '命宮',     // 0
  '兄弟宮',   // 1
  '夫妻宮',   // 2
  '子女宮',   // 3
  '財帛宮',   // 4
  '疾厄宮',   // 5
  '遷移宮',   // 6（對宮）
  '僕役宮',   // 7
  '官祿宮',   // 8
  '田宅宮',   // 9
  '福德宮',   // 10
  '父母宮',   // 11
];
```

---

## 2. 宮位流轉公式

當命宮改變時，宮位重新排列：

### 2.1 宮位偏移計算

```
新宮名 = palaceCycle[(原宮名Index - 命宮偏移 + 12) % 12]
```

### 2.2 程式碼實現

```javascript
const resolveShiftedPalace = (sourcePalace, newBasePalace) => {
  const sourceIndex = palaceCycle.indexOf(sourcePalace);
  const baseIndex = palaceCycle.indexOf(newBasePalace);
  if (sourceIndex === -1 || baseIndex === -1) return null;
  // 正確公式：從 newBasePalace 的角度看 sourcePalace 是什麼宮
  return palaceCycle[(baseIndex - sourceIndex + 12) % 12];
};
```

> ⚠️ **重要**：公式必須是 `(baseIndex - sourceIndex)` 而不是 `(sourceIndex - baseIndex)`，方向相反會導致宮位映射錯誤！

---

## 3. 實際範例

### 範例：命宮從「寅」變為「辰」

假設原本命盤：
- 命宮在寅（地支）
- 兄弟宮在卯
- 夫妻宮在辰
- ...

當命宮變為「辰」時（偏移 +2）：
- 新的「命宮」= 辰
- 新的「兄弟宮」= 巳（辰 + 1）
- 新的「夫妻宮」= 午（辰 + 2）
- ...

---

## 4. 對應關係

### 4.1 對宮（6宮位）

| 原宮 | 對宮 |
|:----:|:-----:|
| 命宮 | 遷移宮 |
| 兄弟宮 | 僕役宮 |
| 夫妻宮 | 官祿宮 |
| 子女宮 | 田宅宮 |
| 財帛宮 | 疾厄宮 |
| 福德宮 | 父母宮 |

### 4.2 三合宮

| 本宮 | 三合宮1 | 三合宮2 |
|:----:|:-------:|:-------:|
| 命宮 | 財帛宮 | 官祿宮 |
| 兄弟宮 | 疾厄宮 | 田宅宮 |
| 夫妻宮 | 遷移宮 | 僕役宮 |

---

## 5. 代碼對照

| 實體 | 程式碼位置 |
|------|-----------|
| 宮位流轉函數 | `web/src/components/ZiweiChart.tsx` | `resolveShiftedPalace()` |
| 宮位順序 | `web/src/components/ZiweiChart.tsx` | `palaceCycle` |
| 三方四正計算 | `web/src/components/ZiweiChart.tsx` | `sanFangSet` |

---

## 6. 自我檢查

1. 如果命宮在午（5），那麼遷移宮在哪個地支？
2. 如果命宮在子（11），那麼三合宮（財帛宮和官祿宮）分別在哪個地支？
3. 命宮和哪個宮永遠是對宮關係？

---

**下一章**：[三方四正](../03-CHARTING/03-SANFANG_SIZHENG.md)
