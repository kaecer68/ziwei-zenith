#!/usr/bin/env node
/**
 * Vite Dev Server Watchdog
 * 自動監控並定時重啟 Vite 開發服務器，防止僵死
 * 
 * 使用方式：
 *   node scripts/dev-watchdog.js        # 預設 2 小時重啟
 *   node scripts/dev-watchdog.js 3600   # 自定義 1 小時重啟（秒）
 */

const { spawn } = require('child_process');
const path = require('path');

// 配置
const RESTART_INTERVAL_MS = (parseInt(process.argv[2], 10) || 7200) * 1000; // 預設 2 小時
const HEALTH_CHECK_INTERVAL_MS = 30000; // 30 秒檢查一次健康狀態
const MAX_MEMORY_MB = 512; // 超過此記憶體則重啟

let viteProcess = null;
let restartTimer = null;
let healthCheckTimer = null;
let startTime = Date.now();

function log(message) {
  const timestamp = new Date().toLocaleString('zh-TW');
  console.log(`[${timestamp}] ${message}`);
}

function getMemoryUsageMB() {
  if (!viteProcess) return 0;
  try {
    // 在 macOS/Linux 上使用 ps 獲取記憶體
    const { execSync } = require('child_process');
    const output = execSync(`ps -o rss= -p ${viteProcess.pid}`, { encoding: 'utf-8' });
    const rssKB = parseInt(output.trim(), 10);
    return rssKB / 1024; // 轉換為 MB
  } catch (e) {
    return 0;
  }
}

async function checkHealth() {
  try {
    const response = await fetch('http://localhost:5174/api/v1/health');
    if (!response.ok) {
      log('⚠️ 健康檢查失敗，準備重啟...');
      restartVite();
    }
  } catch (error) {
    log('⚠️ 無法連接到 Vite 服務器，可能已僵死');
    restartVite();
  }
}

function startVite() {
  log('🚀 啟動 Vite 開發服務器...');
  
  const webDir = path.resolve(__dirname, '../web');
  
  viteProcess = spawn('npm', ['run', 'dev'], {
    cwd: webDir,
    stdio: 'inherit',
    env: {
      ...process.env,
      FORCE_COLOR: '1',
    }
  });

  viteProcess.on('error', (err) => {
    log(`❌ 啟動失敗: ${err.message}`);
    process.exit(1);
  });

  viteProcess.on('exit', (code) => {
    if (code !== 0 && code !== null) {
      log(`⚠️ Vite 異常退出 (code: ${code})，5 秒後重啟...`);
      setTimeout(startVite, 5000);
    }
  });

  startTime = Date.now();
  log(`✅ Vite 已啟動 (PID: ${viteProcess.pid})`);
}

function restartVite() {
  log('🔄 正在重啟 Vite 服務器...');
  
  if (viteProcess) {
    viteProcess.kill('SIGTERM');
    
    // 5 秒後強制終止
    setTimeout(() => {
      if (viteProcess && !viteProcess.killed) {
        log('⚠️ 強制終止進程...');
        viteProcess.kill('SIGKILL');
      }
    }, 5000);
  }

  // 等待終止後重新啟動
  setTimeout(startVite, 6000);
}

function scheduleRestart() {
  const hours = Math.floor(RESTART_INTERVAL_MS / 3600000);
  const minutes = Math.floor((RESTART_INTERVAL_MS % 3600000) / 60000);
  log(`⏰ 已設定定時重啟: ${hours} 小時 ${minutes} 分鐘`);

  restartTimer = setInterval(() => {
    const uptime = (Date.now() - startTime) / 1000;
    log(`⏰ 定時重啟觸發 (運行時間: ${Math.floor(uptime / 60)} 分鐘)`);
    restartVite();
  }, RESTART_INTERVAL_MS);
}

function scheduleHealthCheck() {
  healthCheckTimer = setInterval(async () => {
    // 記憶體檢查
    const memoryMB = getMemoryUsageMB();
    if (memoryMB > MAX_MEMORY_MB) {
      log(`⚠️ 記憶體使用過高 (${Math.round(memoryMB)}MB > ${MAX_MEMORY_MB}MB)，執行重啟...`);
      restartVite();
      return;
    }

    // 健康檢查（只在運行一段時間後執行）
    const uptime = (Date.now() - startTime) / 1000;
    if (uptime > 60) { // 啟動 1 分鐘後開始檢查
      await checkHealth();
    }
  }, HEALTH_CHECK_INTERVAL_MS);
}

// 優雅終止處理
process.on('SIGINT', () => {
  log('👋 收到 SIGINT，正在關閉...');
  cleanup();
});

process.on('SIGTERM', () => {
  log('👋 收到 SIGTERM，正在關閉...');
  cleanup();
});

function cleanup() {
  if (restartTimer) clearInterval(restartTimer);
  if (healthCheckTimer) clearInterval(healthCheckTimer);
  
  if (viteProcess) {
    viteProcess.kill('SIGTERM');
    setTimeout(() => process.exit(0), 2000);
  } else {
    process.exit(0);
  }
}

// 主程序
log('🔧 Vite Dev Watchdog 啟動');
log(`   - 定時重啟間隔: ${Math.floor(RESTART_INTERVAL_MS / 3600000)} 小時`);
log(`   - 健康檢查間隔: ${HEALTH_CHECK_INTERVAL_MS / 1000} 秒`);
log(`   - 記憶體上限: ${MAX_MEMORY_MB} MB`);
log('');

startVite();
scheduleRestart();
scheduleHealthCheck();

// 保持進程運行
process.stdin.resume();
