# 紫微斗數技能驗證清單

> 修改代碼前，請確認已檢查以下關鍵技能點。

---

## 🔴 極重要（必須正確，否則導致嚴重錯誤）

### 1. 宮位映射計算
- [ ] 公式是否為 `(baseIndex - sourceIndex + 12) % 12`？
- [ ] 是否使用 `palaceCycle` 數組？
- [ ] 結果是否正確顯示 `-->新宮名`？

**驗證範例**：
```typescript
// 命宮→官祿的財帛宮（從官祿順時針數8位到命宮）
const result = palaceCycle[(8 - 0 + 12) % 12]; // 應為「財帛宮」
```

**參考文件**：`.ziwei-skills/03-CHARTING/02-PALACE_CYCLE.md`

---

### 2. 星曜亮度系統
- [ ] 六吉星是否有亮度計算？
- [ ] 六煞星是否有亮度計算？
- [ ] 祿存、天馬是否有亮度計算？
- [ ] 亮度是否顯示在星曜名稱後面？
- [ ] 是否使用 `assistant_star_details` 欄位？

**驗證範例**：
```go
// 六吉星
AuspiciousBrightnessLevel(AuspiciousZuofu, BranchZi) // 應返回 "旺"

// 六煞星
MaleficBrightnessLevel(MaleficQingyang, BranchChou)  // 應返回 "廟"

// 祿存
LuCunBrightnessLevel(LuCun, BranchYin)              // 應返回 "旺"
```

**參考文件**：`.ziwei-skills/02-STARS/04-BRIGHTNESS.md`

---

### 3. 星曜層級顯示
- [ ] 六吉星是否單獨列顯示？
- [ ] 是否與主星使用相同格式？
- [ ] 是否顯示圖標、名稱、亮度、四化？

**正確格式**：
```
⭐紫微星 廟 化祿 飛化權
✨左輔星 旺 化科
✨右弼星 得
```

**錯誤格式**：
```
左輔 右弼 文昌 文曲  // ❌ 混排成一列
```

**參考文件**：`.ziwei-skills/02-STARS/02-ASSISTANT_STARS.md`

---

## 🟡 重要（影響功能正確性）

### 4. 四化飛星計算
- [ ] 是否正確查找四化表？
- [ ] 是否遍歷所有宮位找星曜？
- [ ] 生年四化、本宮飛化、選中宮位飛化是否區分顏色？

**參考文件**：`.ziwei-skills/03-CHARTING/04-FLY_HUA.md`

---

### 5. 三方四正計算
- [ ] 對宮是否為 `(index + 6) % 12`？
- [ ] 三合宮是否為 `(index + 4) % 12` 和 `(index + 8) % 12`？
- [ ] 高亮顯示是否正確？

**參考文件**：`.ziwei-skills/03-CHARTING/03-SANFANG_SIZHENG.md`

---

### 6. 大限流年計算
- [ ] 順逆行是否正確（陽男陰女順行、陰男陽女逆行）？
- [ ] 起運歲數是否正確？

**參考文件**：`.ziwei-skills/04-TEMPORAL/`

---

## 🟢 一般（代碼質量）

### 7. API 一致性
- [ ] REST API 是否返回與 gRPC 相同的欄位？
- [ ] 新增欄位是否同時更新兩個協議？
- [ ] JSON 欄位名稱是否使用 snake_case？

### 8. UI/UX 規範
- [ ] 是否遵循 `.agent/skills/ui-ux-pro-max/skill.md`？
- [ ] 頁籤是否使用圓角方形（非圓形）？
- [ ] 按鈕是否有 hover 效果？

---

## 技能文件索引

### 按主題查找
| 主題 | 核心文件 |
|------|---------|
| 宮位映射 | `.ziwei-skills/03-CHARTING/02-PALACE_CYCLE.md` |
| 星曜亮度 | `.ziwei-skills/02-STARS/04-BRIGHTNESS.md` |
| 輔助星 | `.ziwei-skills/02-STARS/02-ASSISTANT_STARS.md` |
| 四化飛星 | `.ziwei-skills/03-CHARTING/04-FLY_HUA.md` |
| 三方四正 | `.ziwei-skills/03-CHARTING/03-SANFANG_SIZHENG.md` |
| 大限流年 | `.ziwei-skills/04-TEMPORAL/01-DA_YUN.md` |
| API 規範 | `AGENTS.md` + `pkg/service/AGENTS.md` |
| UI/UX | `.agent/skills/ui-ux-pro-max/skill.md` |

### 按代碼位置查找
| 代碼文件 | 相關技能文件 |
|----------|-------------|
| `pkg/basis/brightness.go` | `.ziwei-skills/02-STARS/04-BRIGHTNESS.md` |
| `pkg/engine/dayun.go` | `.ziwei-skills/04-TEMPORAL/*.md` |
| `web/src/components/ZiweiChart.tsx` | `.ziwei-skills/03-CHARTING/02-04.md` |
| `cmd/ziwei-server/main.go` | `AGENTS.md` API 章節 |

---

## 常見錯誤對照表

| 錯誤描述 | 正確做法 | 參考文件 |
|---------|---------|---------|
| 宮位映射方向相反 | 使用 `(baseIndex - sourceIndex)` | `02-PALACE_CYCLE.md` |
| 輔星沒有亮度 | 添加 `assistant_star_details` | `04-BRIGHTNESS.md` |
| 六吉星混排 | 單獨列顯示，與主星同級 | `02-ASSISTANT_STARS.md` |
| 只更新 gRPC | 同時更新 REST | `AGENTS.md` |
| 頁籤使用圓形 | 使用圓角方形 | `ui-ux-pro-max/skill.md` |

---

## 驗證命令

```bash
# 1. 檢查宮位映射
curl -s http://localhost:8083/api/v1/calculate -d '{...}' | grep -o '"[^"]*→[^"]*"'

# 2. 檢查輔星亮度
curl -s http://localhost:8083/api/v1/calculate -d '{...}' | grep -A2 'assistant_star_details'

# 3. 檢查四化
curl -s http://localhost:8083/api/v1/calculate -d '{...}' | grep -A5 'fly_hua'

# 4. 運行測試
go test ./pkg/engine/...
```

---

**最後更新**：2026-03-23
**版本**：v1.0
