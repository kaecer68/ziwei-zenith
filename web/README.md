# Ziwei Zenith Web

Ziwei Zenith 的前端介面，使用 React + TypeScript + Vite，提供紫微排盤輸入、結果展示、紀錄管理與解盤視圖。

## 技術棧

- React 19
- TypeScript 5
- Vite 8
- Framer Motion
- TailwindCSS
- Axios

## 啟動與建置

```bash
# 安裝依賴
npm install

# 同步契約 port 設定（產生 .env.ports）
# 此步驟必須先執行，否則 vite.config.ts 會報錯
make sync-contracts

# 開發模式（一般啟動）
npm run dev

# 開發模式（定時重啟看門狗）- 防止 Vite 僵死
npm run dev:safe
# 或指定重啟間隔（秒）
node ../scripts/dev-watchdog.js 3600  # 1 小時重啟

# 生產建置
npm run build

# 本地預覽 build 結果
npm run preview
```

## 從根目錄啟動（推薦）

```bash
# 僅啟動後端（REST :8083 + gRPC :50053）
make run

# 僅啟動前端（開發模式）
make web-dev

# 僅啟動前端（定時重啟看門狗，預設 2 小時）
make web-dev-safe

# 同時啟動後端 + 前端看門狗
make dev-all
```

## 後端依賴

前端透過 `vite.config.ts` 的 proxy 設定串接本機後端。Proxy 目標 port 從 `../.env.ports` 讀取 `REST_PORT` 欄位，因此**必須先執行 `make sync-contracts`** 確保該檔案存在。

> **Port 來源**：`.env.ports` 由 `scripts/sync-contracts.sh` 根據 `contracts/runtime/ports.env` 生成，預設為 8083。若契約變更 port，前端 proxy 會自動同步。

主要 API：
- `POST /api/v1/calculate`
- `POST /api/v1/calculate/temporal`
- `GET /api/v1/records`
- `POST /api/v1/records`
- `DELETE /api/v1/records/:id`
- `GET /api/v1/tags`

### 流運欄位語義（前後端整合）

`TemporalPalaceData` 的時間層與宮位層必須分離處理：

- `stem` + `time_branch`：時間干支（例如 `辛卯`）
- `branch` + `palace`：流運落宮（例如 `未宮 / 遷移宮`）

前端不得使用 `stem + branch` 拼接時間干支，以避免流月/流日標籤錯誤。

## Vite 僵死問題解決方案

長時間運行的 Vite 開發服務器可能因 FSEvents 累積而進入僵死狀態（無法響應 HTTP 請求）。

### 看門狗功能（`scripts/dev-watchdog.js`）

- **定時重啟**：預設每 2 小時自動重啟 Vite
- **記憶體監控**：超過 512MB 自動重啟
- **健康檢查**：每 30 秒檢查 `/api/v1/health`，失敗則重啟
- **異常恢復**：Vite 異常退出時自動重啟

### 使用方法

```bash
# 在專案根目錄執行
make web-dev-safe

# 或自定義重啟間隔（秒）
node scripts/dev-watchdog.js 3600  # 1 小時
```

## 目錄結構

```text
web/
├── src/
│   ├── App.tsx
│   ├── main.tsx
│   └── components/
│       ├── ZiweiChart.tsx
│       ├── DirectoryView.tsx
│       └── InterpretationPanel.tsx
├── public/
├── package.json
└── vite.config.ts
```

## 開發規範

- 使用 TypeScript 嚴格型別，不使用 `any` 或 `@ts-ignore`
- 使用繁體中文作為使用者可見文案
- 維持元件職責單一，避免把後端邏輯混入 UI 元件
