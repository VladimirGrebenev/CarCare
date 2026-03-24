import { defineConfig } from 'vitest/config';

export default defineConfig(async () => {
	const { sveltekit } = await import('@sveltejs/kit/vite');

	return {
		plugins: [sveltekit()],
		test: {
			environment: 'jsdom',
			globals: true,
			include: ['tests/**/*.test.ts'],
			exclude: ['tests/**/*.spec.ts', 'tests/e2e/**'],
			setupFiles: [],
			coverage: {
				provider: 'v8',
				reporter: ['text', 'html']
			}
		}
	};
});
