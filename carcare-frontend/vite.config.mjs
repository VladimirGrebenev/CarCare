import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  test: {
    environment: 'jsdom',
    globals: true,
    include: ['tests/**/*.test.ts', 'tests/**/*.spec.ts'],
    setupFiles: [],
    coverage: {
      reporter: ['text', 'html']
    }
  },
  build: {
    outDir: 'build'
  },
  server: {
    port: 5173
  },
  base: '/',
  publicDir: 'public',
  assetsInclude: [
    '**/manifest.webmanifest',
    '**/service-worker.js',
    '**/offline.html',
    '**/icons/**'
  ]
});
