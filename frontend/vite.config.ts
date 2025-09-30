// 1. 导入 path 模块
import { defineConfig } from 'vite'
import { resolve } from 'path' // 新增导入
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  // 2. 将 resolve 选项移至 plugins 外部，成为 defineConfig 的直接属性
  resolve: {
    alias: [
      {
        find: "@",
        replacement: resolve(__dirname, "./src")
      }
    ]
  }
})