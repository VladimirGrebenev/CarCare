import { sveltekit } from '@sveltejs/kit/vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [
    svelte(),
    sveltekit()
  ],
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
    outDir: 'build',
  },
  server: {
    port: 5173,
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
