import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  base: "/h",
  server: {
    proxy: {
      "/api": {
        rewrite: (path) => path.replace(/^\/api/, ""),
        target: "http://localhost:44551",
      },
      "/per": { target: "ws://localhost:44551", ws: true },
    },
  },
});
