import fs from 'node:fs'
import path from 'node:path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

function resolveRestPort(): string {
  const portsFile = path.resolve(__dirname, '../.env.ports')
  if (!fs.existsSync(portsFile)) {
    throw new Error('.env.ports not found. Please run: make sync-contracts')
  }

  const content = fs.readFileSync(portsFile, 'utf-8')
  const match = content.match(/^REST_PORT=(\d+)$/m)
  if (!match) {
    throw new Error('REST_PORT not found in .env.ports')
  }
  return match[1]
}

const restPort = resolveRestPort()

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    proxy: {
      '/api': {
        target: `http://localhost:${restPort}`,
        changeOrigin: true,
      },
    },
  },
})
