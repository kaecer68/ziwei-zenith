# AI Skills: 流運一致性維護（Ziwei Zenith）

此文件整理「大限/流年/流月/流日」在前後端協作時不可忽略的規則，供後續 AI 迭代直接套用。

## 1) 欄位語義（不可混用）

`TemporalPalaceData` 必須遵守以下語義：

- `stem`: 流運時間天干（時間層）
- `time_branch`: 流運時間地支（時間層）
- `branch`: 流運落宮地支（宮位層）
- `palace`: 流運落宮宮名（宮位層）

禁止用 `stem + branch` 直接拼「時間干支」，正確應為 `stem + time_branch`。

## 2) 單一來源建構規則（前端）

在 `web/src/App.tsx` 以 `temporalPositions` 作為命盤頁流運資料唯一來源：

- 命盤中央流運標籤
- 命盤宮位高亮（大限/流年/流月/流日）
- 「流運定位」面板四格

上述三處不得各自重算或各自組字，避免東改西壞。

## 3) REST / gRPC / OpenAPI 對齊清單

每次調整流運欄位，必須同步以下四處：

1. `contracts/openapi/ziwei-zenith.yaml`
2. `proto/ziwei.proto`
3. `pkg/api/v1/types.go`
4. `cmd/ziwei-server/main.go` + `pkg/service/grpc_server.go`

並確認 REST `/api/v1/calculate`、`/api/v1/calculate/temporal` 及 gRPC `ZiweiService/Calculate` 回傳語義一致。

## 4) 命名正規化規則（宮位別名）

- `僕役宮` 與 `交友宮` 視為同義。
- 三方四正計算、查表、飛入比對皆需先正規化。
- 若未正規化，會出現少一宮或 `—酉` 類型空值問題。

## 5) 驗證流程（每次必跑）

```bash
openapi-generator validate -i contracts/openapi/ziwei-zenith.yaml
go build ./...
npm run build --prefix web
```

若改動 `proto/ziwei.proto`，需先再生 gRPC：

```bash
PATH="$(go env GOPATH)/bin:$PATH" protoc \
  --go_out=pkg/api/grpc/v1 --go_opt=paths=source_relative \
  --go-grpc_out=pkg/api/grpc/v1 --go-grpc_opt=paths=source_relative \
  -I proto proto/ziwei.proto
```
