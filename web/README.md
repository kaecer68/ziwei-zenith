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

# 開發模式
npm run dev

# 生產建置
npm run build

# 本地預覽 build 結果
npm run preview
```

## 後端依賴

前端透過 `vite.config.ts` 的 proxy 設定串接本機後端。Proxy 目標 port 從 `../.env.ports` 讀取 `REST_PORT` 欄位，因此**必須先執行 `make sync-contracts`** 確保該檔案存在。

> **Port 來源**：`.env.ports` 由 `scripts/sync-contracts.sh` 根據 `contracts/runtime/ports.env` 生成，預設為 8083。若契約變更 port，前端 proxy 會自動同步。

主要 API：
- `POST /api/v1/calculate`
- `GET /api/v1/records`
- `POST /api/v1/records`
- `DELETE /api/v1/records/:id`
- `GET /api/v1/tags`

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
