import path from "node:path"
import { TanStackRouterVite } from "@tanstack/router-vite-plugin"
import react from "@vitejs/plugin-react-swc"
import { defineConfig } from "vite"
import mkcert from "vite-plugin-mkcert" // awesome guys // https://stackoverflow.com/questions/69417788/vite-https-on-localhost

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    // host: '0.0.0.0'
    host: '192.168.1.12',
    cors: {
      origin: "*"
    },
    allowedHosts: true,
    proxy: {
      '/vault-unseal': {
        // target: 'https://192.168.1.12:8751',
        target: 'https://navyseals.pesto.io:8751',
        changeOrigin: true, // Optional: Change origin header to match frontend
        secure: false, // Optional: Disable checks for https on target
      },
      '/vault-status': {
        // target: 'https://192.168.1.12:8751',
        target: 'https://navyseals.pesto.io:8751',
        changeOrigin: true, // Optional: Change origin header to match frontend
        secure: false, // Optional: Disable checks for https on target
      },
      '/vault-init': {
        // target: 'https://192.168.1.12:8751',
        target: 'https://navyseals.pesto.io:8751',
        changeOrigin: true, // Optional: Change origin header to match frontend
        secure: false, // Optional: Disable checks for https on target
      },
      '/vault-seal': {
        // target: 'https://192.168.1.12:8751',
        target: 'https://navyseals.pesto.io:8751',
        changeOrigin: true, // Optional: Change origin header to match frontend
        secure: false, // Optional: Disable checks for https on target
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  plugins: [react(), TanStackRouterVite(), mkcert()],
})
