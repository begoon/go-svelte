import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
    plugins: [svelte(), tailwindcss()],
    root: "site",
    build: {
        emptyOutDir: true,
        outDir: "../dist",
        rollupOptions: {
            output: {
                assetFileNames: "assets/[name].[ext]",
                chunkFileNames: "assets/[name].js",
                entryFileNames: "assets/entry-[name].js",
            },
            input: {
                main: "site/index.html",
                about: "site/about/index.html",
            },
        },
    },
});
