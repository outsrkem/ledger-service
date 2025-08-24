import { defineConfig } from "vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";
import { resolve } from "path";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { visualizer } from "rollup-plugin-visualizer";
export default defineConfig({
    base: "/ledger/",
    build: {
        outDir: "dist",
        rollupOptions: {
            input: {
                main: resolve(__dirname, "index.html"),
            },
        },
    },
    plugins: [
        vue(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
        }),
        Components({
            resolvers: [ElementPlusResolver()],
        }),
        visualizer({
            open: false, // 构建后自动打开分析报告
            gzipSize: true,
        }),
    ],
    esbuild: {
        drop: ["console", "debugger"],
    },
    server: {
        proxy: {
            "/api": {
                target: "https://uias.localvm.outsrkem.top:30078",
                // rewrite: (path) => path.replace(/^\/api/, ""),
                changeOrigin: true,
            },
            "/authui": {
                target: "https://uias.localvm.outsrkem.top:30078",
                // rewrite: (path) => path.replace(/^\/api/, ""),
                changeOrigin: true,
            },
        },
    },
});
