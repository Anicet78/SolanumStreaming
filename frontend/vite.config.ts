import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  server: {
    host: "0.0.0.0",
    allowedHosts: ["localhost", "solanumstreaming.mooo.com"],
    port: 5173,
    proxy: {
      "/api/auth": {
        target: "http://auth:8081",
        rewrite: (p) => p.replace(/^\/api\/auth/, ""),
      },
      "/api/movies": {
        target: "http://movies:8082",
        rewrite: (p) => p.replace(/^\/api\/movies/, ""),
      },
      "/api/stream": {
        target: "http://stream:8083",
        rewrite: (p) => p.replace(/^\/api\/stream/, ""),
      },
    },
  },
  plugins: [
    solid(),
    tailwindcss()
  ],
})
