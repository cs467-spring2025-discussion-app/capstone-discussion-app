import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());

  const fast_api_target = env.VITE_FASTAPI_URL || "http://localhost:8080";
  const auth_api_target = env.VITE_AUTH_API_URL || "http://localhost:3001"

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve("./src"),
      },
    },
    server: {
      host: "0.0.0.0",
      port: 5173,
      proxy: {
        // Proxy API requests
        "/authapi": {
          target: auth_api_target,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/authapi/, ""),
        },
        "/fastapi": {
          target: fast_api_target,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/fastapi/, ""),
        },
      },
    },
    preview: {
      host: "0.0.0.0",
      port: 4173,
    },
  };
});
