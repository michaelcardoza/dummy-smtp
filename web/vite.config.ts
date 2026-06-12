import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(),svelte()],
  build: {
    outDir: '../internal/infrastructure/web/dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8025',
      '/events': {
        target: 'http://localhost:8025',
        changeOrigin: true,
      },
    },
  },
})
