import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {ElementPlusResolver} from 'unplugin-vue-components/resolvers'
import ElementPlus from 'unplugin-element-plus/vite'

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        ElementPlus({
            useSource: true,
        }),
        AutoImport({
            resolvers: [ElementPlusResolver({importStyle: 'sass'})],
        }),
        Components({
            resolvers: [ElementPlusResolver({importStyle: 'sass'})],
        }),
        vueDevTools(),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        },
    },
    css: {
        preprocessorOptions: {
            scss: {
                additionalData: `
        @use "@/assets/index.scss" as *;
        `,
            }
        }
    },
    // 配置服务器代理
    server: {
        proxy: {
            '/gapi': {
                target: 'http://127.0.0.1:8792',                                            // 目标地址
                changeOrigin: true,                                                         // 是否改变源地址
                rewrite: (path) => path.replace(/^\/gapi/, '') // 重写路径
            }
        }
    },
    build: {
        rollupOptions: {
            input: {
                main: './public/index.html'
            }
        }
    },
    dev:{
        rollupOptions: {
            input: {
                main: './public/index.html'
            }
        }
    }
})
