// vite.config.js
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver, VantResolver } from "unplugin-vue-components/resolvers";
import { resolve } from "path";

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), "");
    const API_Endpoint = env.VITE_API_Endpoint;

    return {
        base: "/ledger/",
        build: {
            outDir: "dist",
            minify: "esbuild",
            rollupOptions: {
                input: {
                    main: resolve(__dirname, "index.html"),
                },
            },
            terserOptions: {
                compress: {
                    drop_console: true,
                    drop_debugger: true,
                    collapse_vars: true,
                },
                keep_fnames: false,
            },
        },
        plugins: [
            vue(),
            AutoImport({
                resolvers: [ElementPlusResolver(), VantResolver()],
            }),
            Components({
                resolvers: [ElementPlusResolver(), VantResolver()],
            }),
        ],
        server: {
            host: "0.0.0.0",
            proxy: {
                "/api": {
                    target: API_Endpoint,
                    changeOrigin: true,
                    secure: false,
                },
                "/authui": {
                    target: API_Endpoint,
                    changeOrigin: true,
                    secure: false,
                },
            },
        },
        resolve: {
            alias: {
                "@": resolve(__dirname, "src"),
            },
        },
    };
});
