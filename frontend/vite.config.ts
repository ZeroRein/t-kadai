import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: true,
    proxy: {
      "/api": {
        target: "http://app:8080",
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // backend handles /api group
      },
    },
  },
});
