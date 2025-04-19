import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/ws': {
				target: process.env.SERVER_URL || 'ws://backend:3000',
				ws: true,
				changeOrigin: true,
			},
			'/api': {
				target: process.env.SERVER_URL || 'ws://backend:3000',
				changeOrigin: true
			}
		}
	}
});
