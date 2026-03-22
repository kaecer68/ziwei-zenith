.PHONY: dev run sync-contracts verify-contracts dev-clean web-dev web-dev-safe web-build dev-all

dev: run

run:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/ziwei-server/main.go'

sync-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh

verify-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh --check

dev-clean:
	@chmod +x scripts/dev-clean.sh
	bash scripts/dev-clean.sh

# 前端開發命令
web-dev:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	cd web && npm run dev

web-dev-safe:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	node scripts/dev-watchdog.js

web-build:
	cd web && npm run build

# 同時啟動後端和前端（使用看門狗）
dev-all:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	@echo "啟動後端服務..."
	@bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/ziwei-server/main.go' &
	@sleep 3
	@echo "啟動前端看門狗（自動定時重啟）..."
	@node scripts/dev-watchdog.js
