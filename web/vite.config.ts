import fs from 'node:fs'
import path from 'node:path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

function resolveRuntimeEnvFile(): string {
  const envPortsFile = process.env.ENV_PORTS_FILE
  if (envPortsFile) {
    return path.resolve(envPortsFile)
  }

  const portsFile = path.resolve(__dirname, '../.env.ports')
  if (!fs.existsSync(portsFile)) {
    throw new Error('.env.ports not found. Please run: make sync-contracts')
  }
  return portsFile
}

function mustRuntimeValue(key: string): string {
  const content = fs.readFileSync(resolveRuntimeEnvFile(), 'utf-8')
  const match = content.match(new RegExp(`^${key}=(.+)$`, 'm'))
  if (!match) {
    throw new Error(`${key} not found in .env.ports`)
  }
  return match[1].trim()
}

const apiTarget = mustRuntimeValue('VITE_API_TARGET')

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    proxy: {
      '/api': {
        target: apiTarget,
        changeOrigin: true,
      },
    },
  },
})
