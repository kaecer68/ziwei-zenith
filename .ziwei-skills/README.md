# 紫微斗數技能圖譜 (Ziwei Skills)

> **⚠️ 重要聲明**：修改紫微斗數相關代碼前，**必須**完成以下學習路徑。否則任何修改都將視為無效補丁。

---

## 學習路徑（必讀）

```
完成以下學習模組，方可進行代碼修改：
```

```
Level 1: 入門基礎
├── 01-FOUNDATION/01-STEM_BRANCH.md       ◄ 天干地支（陰陽五行）
├── 01-FOUNDATION/02-PALACE_SYSTEM.md     ◄ 十二宮系統
└── 01-FOUNDATION/03-LUNAR_CALENDAR.md    ◄ 農曆與節氣

Level 2: 星曜認知
├── 02-STARS/01-MAIN_STARS.md             ◄ 主星系統
├── 02-STARS/02-ASSISTANT_STARS.md        ◄ 輔助星
├── 02-STARS/03-TRANSFORM_STARS.md        ◄ 四化星（祿權科忌）
└── 02-STARS/04-BRIGHTNESS.md             ◄ 星曜亮度

Level 3: 排盤核心（重要！修改代碼前必讀）
├── 03-CHARTING/01-BIRTH_CHART.md         ◄ 本命盤排法
├── 03-CHARTING/02-PALACE_CYCLE.md        ◄ 宮位流轉
├── 03-CHARTING/03-SANFANG_SIZHENG.md     ◄ 三方四正
└── 03-CHARTING/04-FLY_HUA.md            ◄ 飛星四化計算

Level 4: 動態運限
├── 04-TEMPORAL/01-DA_YUN.md              ◄ 大限計算
├── 04-TEMPORAL/02-LIU_NIAN.md            ◄ 流年計算
├── 04-TEMPORAL/03-LIU_YUE.md             ◄ 流月計算
└── 04-TEMPORAL/04-LIU_RI.md             ◄ 流日計算

Level 5: 驗證測試
└── 05-VALIDATION/01-TEST_CASES.md        ◄ 測試用例與驗證
```

---

## 模組地圖

```
                    ┌─────────────────────────────────────┐
                    │         紫微斗數技能圖譜             │
                    └─────────────────────────────────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
        ▼                           ▼                           ▼
   ┌─────────┐               ┌─────────┐               ┌─────────┐
   │ Foundation│               │  Stars  │               │ Charting │
   │  基礎理論 │               │ 星曜系統 │               │ 排盤核心 │
   └─────────┘               └─────────┘               └─────────┘
        │                           │                           │
   天干地支                    主星                      本命盤
   十二宮                     四化                      三方四正
   農曆                       亮度                      飛星四化
        │                           │                           │
        └───────────────────────────┼───────────────────────────┘
                                    ▼
                            ┌─────────────┐
                            │  Temporal   │
                            │   動態運限   │
                            └─────────────┘
                                    │
                            大限 → 流年 → 流月 → 流日
                                    │
                                    ▼
                            ┌─────────────┐
                            │ Validation  │
                            │   驗證測試   │
                            └─────────────┘
```

---

## 前置知識檢查清單

修改代碼前，請確認已理解以下概念：

- [ ] 天干地支的陰陽五行屬性
- [ ] 十二宮的順序與意義（命宮→父母宮）
- [ ] 五行局（火6/水2/木3/金4/土5）
- [ ] 主星、輔助星的區別
- [ ] 四化（祿權科忌）的意義
- [ ] 生年四化 vs 飛星四化的差異
- [ ] 三方四正的判定（對宮+兩個三合宮）
- [ ] 大限/流年/流月/流日的計算順序

---

## 代碼修改規範

任何紫微斗數相關代碼修改，**必須**：

1. **先閱讀相關 skill 文件**
2. **理解現有實現邏輯**
3. **確認修改影響範圍**
4. **執行驗證測試**

違反以上規範的修改，將被拒絕合併。

---

## 技能索引總表

### Foundation（基礎）
| 文件 | 技能點 | 依賴 | 代碼對照 |
|------|--------|------|----------|
| `01-FOUNDATION/01-STEM_BRANCH.md` | 天干地支、陰陽五行 | - | `pkg/basis/definitions.go` |
| `01-FOUNDATION/02-PALACE_SYSTEM.md` | 十二宮系統 | 01-STEM_BRANCH | `pkg/basis/palaces.go` |
| `01-FOUNDATION/03-LUNAR_CALENDAR.md` | 農曆節氣 | 01-STEM_BRANCH | 外部庫 `lunar-zenith` |

### Stars（星曜）
| 文件 | 技能點 | 依賴 | 代碼對照 |
|------|--------|------|----------|
| `02-STARS/01-MAIN_STARS.md` | 十四主星 | 01-FOUNDATION | `pkg/basis/stars.go` |
| `02-STARS/02-ASSISTANT_STARS.md` | 六吉星、六煞星 | 02-MAIN_STARS | `pkg/basis/auspicious.go`, `malefic.go` |
| `02-STARS/03-TRANSFORM_STARS.md` | 四化（祿權科忌） | 02-MAIN_STARS | `pkg/basis/transformation.go` |
| `02-STARS/04-BRIGHTNESS.md` | 星曜亮度（廟旺利陷） | 02-MAIN_STARS | `pkg/basis/brightness.go` |

### Charting（排盤）
| 文件 | 技能點 | 依賴 | 代碼對照 |
|------|--------|------|----------|
| `03-CHARTING/01-BIRTH_CHART.md` | 本命盤排法 | 01-FOUNDATION, 02-STARS | `pkg/engine/` |
| `03-CHARTING/02-PALACE_CYCLE.md` | 宮位流轉、映射計算 | 03-BIRTH_CHART | `web/src/components/ZiweiChart.tsx` |
| `03-CHARTING/03-SANFANG_SIZHENG.md` | 三方四正 | 03-BIRTH_CHART | `web/src/components/ZiweiChart.tsx` |
| `03-CHARTING/04-FLY_HUA.md` | 飛星四化計算 | 02-TRANSFORM_STARS | `web/src/components/ZiweiChart.tsx` |

### Temporal（動態運限）
| 文件 | 技能點 | 依賴 | 代碼對照 |
|------|--------|------|----------|
| `04-TEMPORAL/01-DA_YUN.md` | 大限計算 | 03-CHARTING | `pkg/engine/dayun.go` |
| `04-TEMPORAL/02-LIU_NIAN.md` | 流年計算 | 04-DA_YUN | `pkg/engine/dayun.go` |
| `04-TEMPORAL/03-LIU_YUE.md` | 流月計算 | 04-LIU_NIAN | `pkg/engine/dayun.go` |
| `04-TEMPORAL/04-LIU_RI.md` | 流日計算 | 04-LIU_YUE | `pkg/engine/dayun.go` |

### Validation（驗證）
| 文件 | 技能點 | 依賴 | 代碼對照 |
|------|--------|------|----------|
| `05-VALIDATION/01-TEST_CASES.md` | 測試用例 | 全部 | `pkg/engine/*_test.go` |

---

## 關鍵技能快速查找

### 按主題查找
| 主題 | 核心文件 | 應用場景 |
|------|---------|---------|
| 天干地支 | `01-STEM_BRANCH.md` | 排盤基礎 |
| 十二宮 | `02-PALACE_SYSTEM.md` | 宮位定義 |
| 四化飛星 | `03-TRANSFORM_STARS.md` + `04-FLY_HUA.md` | 飛化計算 |
| 三方四正 | `04-SANFANG_SIZHENG.md` | 論命核心 |
| 宮位映射 | `03-PALACE_CYCLE.md` | 動態視角 |
| 大限流年 | `04-DA_YUN.md` | 動態運限 |
| 星曜亮度 | `04-BRIGHTNESS.md` | 亮度計算 |
| API 規範 | `../AGENTS.md` | 前後端一致性 |

### 按代碼位置查找
| 代碼文件 | 相關技能文件 |
|----------|-------------|
| `pkg/basis/brightness.go` | `02-STARS/04-BRIGHTNESS.md` |
| `pkg/engine/dayun.go` | `04-TEMPORAL/*.md` |
| `web/src/components/ZiweiChart.tsx` | `03-CHARTING/02-04.md` |
| `cmd/ziwei-server/main.go` | `../AGENTS.md` API 章節 |

---

## 🚨 修改前必讀

### 技能驗證清單
修改代碼前，**必須**查看：[SKILL_CHECKLIST.md](./SKILL_CHECKLIST.md)

### 關鍵技能速查
| 主題 | 文件 | 重要性 |
|------|------|--------|
| 宮位映射公式 | `03-CHARTING/02-PALACE_CYCLE.md` | ⭐⭐⭐ 極重要 |
| 星曜亮度系統 | `02-STARS/04-BRIGHTNESS.md` | ⭐⭐⭐ 極重要 |
| 四化飛星計算 | `03-CHARTING/04-FLY_HUA.md` | ⭐⭐⭐ 極重要 |
| 三方四正判定 | `03-CHARTING/03-SANFANG_SIZHENG.md` | ⭐⭐⭐ 極重要 |
| 輔助星識別 | `02-STARS/02-ASSISTANT_STARS.md` | ⭐⭐ 重要 |
| 大限流年計算 | `04-TEMPORAL/01-DA_YUN.md` | ⭐⭐ 重要 |
| API 更新規範 | `../AGENTS.md` | ⭐⭐ 重要 |
| UI/UX 設計 | `../.agent/skills/ui-ux-pro-max/skill.md` | ⭐⭐ 重要 |
