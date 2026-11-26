import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  server: {
    host: "0.0.0.0",
    port: 5173,
    proxy: {
      "/api/auth": {
        target: "http://localhost:8081",
        rewrite: (p) => p.replace(/^\/api\/auth/, ""),
      },
      "/api/movies": {
        target: "http://localhost:8082",
        rewrite: (p) => p.replace(/^\/api\/movies/, ""),
      },
      "/api/stream": {
        target: "http://localhost:8083",
        rewrite: (p) => p.replace(/^\/movies/, ""),
      },
    },
  },
  plugins: [
    solid(),
    tailwindcss()
  ],
})
