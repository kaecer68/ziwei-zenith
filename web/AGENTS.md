# web/ - React 前端

## 概述

Web frontend 使用 React 19 + TypeScript 5 + Vite 8，提供紫微排盤的互動界面。

## 技術棧

- **Framework**: React 19
- **Language**: TypeScript 5.9
- **Build**: Vite 8
- **UI**: TailwindCSS, Framer Motion
- **HTTP**: Axios
- **Icons**: Lucide React

## 結構

```
web/
├── src/
│   ├── App.tsx              # 主應用組件
│   ├── main.tsx             # 入口
│   ├── components/          # React 組件
│   ├── styles/              # CSS 樣式
│   └── assets/              # 靜態資源
├── public/                  # 公共文件
├── vite.config.ts           # Vite 配置
└── package.json
```

## 開發命令

```bash
# 安裝依賴
npm install

# 啟動開發服務器
npm run dev

# 生產建置
npm run build

# 從根目錄啟動（推薦）
make web-dev
make web-dev-safe    # 帶看門狗
```

## API Proxy 配置

Vite 開發服務器代理設定：
```ts
proxy: {
  '/api': {
    target: 'http://localhost:8083',
    changeOrigin: true,
  }
}
```

Port 從 `.env.ports` 讀取（由 `make sync-contracts` 生成）。

## Vite 看門狗

長時間運行可能導致 FSEvents 累積僵死，使用看門狗自動重啟：

```bash
make web-dev-safe
```

功能：
- 每 2 小時自動重啟
- 記憶體超過 512MB 自動重啟
- 每 30 秒健康檢查

## 依賴注意

**必須先執行** `make sync-contracts`，否則 `vite.config.ts` 會報錯：
```
.env.ports not found. Please run: make sync-contracts
```

## 流運欄位語義

前後端整合時注意分離時間層與宮位層：
- `stem` + `time_branch` = 時間干支（如：辛卯）
- `branch` + `palace` = 流運落宮（如：未宮 / 遷移宮）

---

## UI/UX Skills - 紫微斗數命盤顯示規範

### 星曜顯示規範

#### 重要輔星識別
**必須單獨列顯示的星曜**（與主星同級）：
- 左輔、右弼
- 文昌、文曲
- 天魁、天鉞

這些星曜在四化系統中有重要作用，必須單獨成行顯示，不可與其他雜星混排。

#### 星曜資訊排列順序
每個星曜的顯示格式：
```
[圖標] [星曜名稱] [亮度] [生年四化] [本宮飛化] [選中宮位飛化]
```

範例：
```
⭐紫微星 廟 化祿 飛化權 命宮化科
⭐天機星 旺      飛化忌
✨左輔星 廟 化科
```

**注意**：亮度必須緊跟在星曜名稱後面，不可置右對齊。

#### 星曜亮度系統

**所有星曜都有亮度，不只是主星**：

| 星曜類別 | 是否需顯示亮度 | 備註 |
|---------|--------------|------|
| 十四主星 | ✅ 必須 | 根據地支宮位查表 |
| 六吉星（左輔、右弼、文昌、文曲、天魁、天鉞）| ✅ 必須 | 根據地支宮位查表 |
| 祿存、天馬 | ✅ 必須 | 根據地支宮位查表 |
| 六煞星（擎羊、陀羅、火星、鈴星、地空、地劫）| ✅ 必須 | 根據地支宮位查表 |
| 丙級雜曜（紅鸞、天喜等）| ❌ 不需要 | 通常無亮度概念 |

**API 數據結構**：
後端通過 `assistant_star_details` 欄位傳遞輔星和煞星的亮度資訊：
```typescript
interface PalaceData {
  star_details?: Array<{ name: string; brightness?: string }>;      // 主星
  assistant_star_details?: Array<{ name: string; brightness?: string }>; // 輔星、煞星
}
```

**前端實現**：
```typescript
// 輔星和煞星從 assistant_star_details 獲取，包含亮度資訊
const assistantStarList = entry.assistant_star_details || [];

// 顯示時
{assistantStarList.map((star) => (
  <div>
    {star.name} {star.brightness}
  </div>
))}
```

### 宮位變化映射規則

#### 映射計算公式
當用戶點擊宮位 B 時，宮位 A 的新名稱計算：
```typescript
// 正確公式：從 B 的角度看 A 是什麼宮
const shiftedName = palaceCycle[(indexB - indexA + 12) % 12];
```

**常見錯誤**：使用 `(indexA - indexB)` 會導致方向相反。

#### 動態映射實現
- `mappingBasePalace` 必須跟隨 `selectedPalaceName` 動態變化
- 預設選中命宮，顯示「命宮→命宮」
- 點擊財帛宮後，命宮顯示「命宮→財帛的官祿宮」

### 四化飛星系統

#### 四化來源區分
必須區分三種四化來源，使用不同顏色：
1. **生年四化**（紅色）：出生年天干決定，終身不變
2. **本宮飛化**（金色）：該宮位天干對本宮星曜的影響
3. **選中宮位飛化**（紫色）：當前選中宮位對所有宮位的影響

#### 互動實現要點
- 點擊任意宮位時，計算該宮位天干對全系統的四化影響
- 在所有宮位的星曜旁顯示選中宮位的四化標記
- 中間面板顯示選中宮位的完整四化資訊

### 頁籤設計規範

- 使用圓角方形（`border-radius: 0.5rem`），不可使用圓形
- 必須防止文字斷行：`white-space: nowrap`
- 圖標與文字水平對齊：`display: inline-flex`

---

## 禁止事項

- 禁止使用 `any` 或 `@ts-ignore`
- 禁止在 UI 組件中混入後端計算邏輯
- 禁止直接修改 `.env.ports`（從契約同步）
- **禁止將重要輔星（左輔、右弼、文昌、文曲、天魁、天鉞）與雜星混排**
- **禁止將亮度置右對齊，應緊跟星曜名稱**
- **禁止將四化顯示為獨立按鈕，應緊跟對應星曜**
