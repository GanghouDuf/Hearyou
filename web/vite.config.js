import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
	plugins: [react()],
	server: {
		proxy: {
			'/ws': {
				target: 'ws://localhost:8080',
				ws: true,
			},
			'/register': 'http://localhost:8080',
			'/login': 'http://localhost:8080',
		},
	},
})