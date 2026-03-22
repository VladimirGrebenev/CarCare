import { sveltekit } from '@sveltejs/kit/vite';

export default {
  plugins: [sveltekit()],
  build: {
    outDir: 'build',
  },
  server: {
    port: 5173,
  },
  base: '/',
  publicDir: 'public',
  // Ensure service worker and manifest are served
  assetsInclude: [
    '**/manifest.webmanifest',
    '**/service-worker.js',
    '**/offline.html',
    '**/icons/**'
  ]
};
