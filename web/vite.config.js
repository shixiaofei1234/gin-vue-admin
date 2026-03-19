import legacyPlugin from '@vitejs/plugin-legacy'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { viteLogo } from './src/core/config'
import Banner from 'vite-plugin-banner'
import * as path from 'path'
import * as dotenv from 'dotenv'
import * as fs from 'fs'
import vuePlugin from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import fullImportPlugin from './vitePlugin/fullImport/fullImport.js'
import VueFilePathPlugin from './vitePlugin/componentName/index.js'
import { svgBuilder } from 'vite-auto-import-svg'
import { AddSecret } from './vitePlugin/secret'

// @see https://cn.vitejs.dev/config/
export default ({ command, mode }) => {
  AddSecret("")
  const NODE_ENV = mode || 'development'
  
  // 修复1：完善环境变量加载（优先加载默认.env，再加载环境专属.env）
  const envFiles = [
    `.env`, // 默认环境文件
    `.env.${NODE_ENV}` // 环境专属文件（覆盖默认）
  ]
  
  // 加载并注入环境变量
  for (const file of envFiles) {
    try {
      // 修复2：增加文件存在性校验，避免文件不存在时报错
      if (fs.existsSync(file)) {
        const envConfig = dotenv.parse(fs.readFileSync(file))
        for (const k in envConfig) {
          process.env[k] = envConfig[k]
        }
      }
    } catch (e) {
      console.warn(`加载环境文件 ${file} 失败:`, e.message)
    }
  }

  viteLogo(process.env)

  const timestamp = Date.parse(new Date())

  const optimizeDeps = {}

  const alias = {
    '@': path.resolve(__dirname, './src'),
    'vue$': 'vue/dist/vue.runtime.esm-bundler.js',
  }

  const esbuild = {}

  const rollupOptions = {
    output: {
      entryFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      chunkFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      assetFileNames: 'assets/087AC4D233B64EB0[name].[hash].[ext]',
    },
  }

  // 修复3：代理配置前置处理，避免undefined
  const proxyConfig = {}
  const baseApi = process.env.VITE_BASE_API
  if (baseApi) {
    const basePath = process.env.VITE_BASE_PATH || ''
    const serverPort = process.env.VITE_SERVER_PORT || ''
    let target = basePath
    // 仅当端口存在时拼接，避免:undefined
    if (serverPort) {
      target += `:${serverPort}`
    }
    target += '/'
    
    // 修复4：校验target合法性，避免无效URL
    if (target && !target.includes('undefined')) {
      proxyConfig[baseApi] = {
        target: target,
        changeOrigin: true,
        rewrite: path => path.replace(new RegExp('^' + baseApi), ''),
      }
    }
  }

  // 修复5：插件数组过滤false值，避免Vite警告
  const plugins = [
    legacyPlugin({
      targets: ['Android > 39', 'Chrome >= 60', 'Safari >= 10.1', 'iOS >= 10.3', 'Firefox >= 54', 'Edge >= 15'],
    }),
    vuePlugin(),
    svgBuilder('./src/assets/icons/'),
    Banner(`\n Build based on gin-vue-admin \n Time : ${timestamp}`),
    VueFilePathPlugin("./src/pathInfo.json")
  ]

  // 条件添加vueDevTools插件，避免插入false
  if (process.env.VITE_POSITION === 'open') {
    plugins.push(vueDevTools({ launchEditor: process.env.VITE_EDITOR }))
  }

  const config = {
    base: '/', // 编译后js导入的资源路径
    root: './', // index.html文件所在位置
    publicDir: 'public', // 静态资源文件夹
    resolve: {
      alias,
    },
    define: {
      // 修复6：保留process.env，避免覆盖导致环境变量失效
      'process.env': process.env
    },
    server: {
      open: true, // 如果使用docker-compose开发模式，设置为false
      port: process.env.VITE_CLI_PORT || 3000, // 修复：添加默认端口，避免undefined
      proxy: proxyConfig // 使用处理后的代理配置
    },
    build: {
      minify: 'terser', // 是否进行压缩,boolean | 'terser' | 'esbuild',默认使用terser
      manifest: false, // 是否产出manifest.json
      sourcemap: false, // 是否产出sourcemap.json
      outDir: 'dist', // 产出目录
      terserOptions: {
        compress: {
          //生产环境时移除console
          drop_console: true,
          drop_debugger: true,
        },
      },
      rollupOptions,
    },
    esbuild,
    optimizeDeps,
    plugins: plugins, // 使用整理后的插件数组
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@use "@/style/element/index.scss" as *;`,
        }
      }
    },
  }

  // 区分开发/生产环境插件
  if (NODE_ENV === 'development') {
    config.plugins.push(fullImportPlugin())
  } else {
    config.plugins.push(
      AutoImport({ resolvers: [ElementPlusResolver()] }),
      Components({
        resolvers: [ElementPlusResolver({ importStyle: 'sass' })]
      })
    )
  }

  return config
}