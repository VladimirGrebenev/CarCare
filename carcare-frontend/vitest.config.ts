import { defineConfig } from 'vitest/config';

export default defineConfig({
	test: {
		environment: 'jsdom',
		globals: true,
		include: ['tests/**/*.test.ts'],
		exclude: ['tests/**/*.spec.ts', 'tests/e2e/**'],
		setupFiles: [],
		coverage: {
			reporter: ['text', 'html']
		}
	}
});
