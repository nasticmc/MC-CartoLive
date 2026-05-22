import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { readFileSync } from 'node:fs';

const packageJSON = JSON.parse(readFileSync(new URL('./package.json', import.meta.url), 'utf8')) as { version?: string };

function buildNumber(): string {
  if (process.env.VITE_BUILD_NUMBER) return process.env.VITE_BUILD_NUMBER;
  if (process.env.GITHUB_SHA) return process.env.GITHUB_SHA.slice(0, 7);
  return new Date().toISOString().replace(/[-:]/g, '').slice(0, 13);
}

export default defineConfig({
  plugins: [react()],
  define: {
    __APP_VERSION__: JSON.stringify(packageJSON.version ?? '1.0.0'),
    __BUILD_NUMBER__: JSON.stringify(buildNumber())
  },
  test: {
    environment: 'jsdom'
  }
});
