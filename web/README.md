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

# 開發模式
npm run dev

# 生產建置
npm run build

# 本地預覽 build 結果
npm run preview
```

## 後端依賴

前端預設串接本機後端：`http://localhost:8083`

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
