import fs from 'node:fs'
import path from 'node:path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

function resolveRestPort(): string {
  const portsFile = path.resolve(__dirname, '../.env.ports')
  if (!fs.existsSync(portsFile)) {
    return '8083'
  }

  const content = fs.readFileSync(portsFile, 'utf-8')
  const match = content.match(/^REST_PORT=(\d+)$/m)
  return match?.[1] ?? '8083'
}

const restPort = resolveRestPort()

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: `http://localhost:${restPort}`,
        changeOrigin: true,
      },
    },
  },
})
