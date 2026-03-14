import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  // base must match your GitHub repo name
  base: "/aloe-raw/",
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "https://aloe-raw-api.onrender.com",
        changeOrigin: true,
      },
      "/auth": {
        target: "https://aloe-raw-api.onrender.com",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
  },
});
