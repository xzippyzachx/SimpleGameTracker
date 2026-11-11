import { fileURLToPath, URL } from 'node:url'
import { existsSync } from 'node:fs'
import { resolve } from 'node:path'
import { VitePWA } from 'vite-plugin-pwa'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  // .env.local is loaded in every mode and outranks a plain .env, so a
  // production build would silently bake in the localhost API if
  // .env.production were missing. Require the mode-specific file instead.
  const modeEnvFile = mode === 'production' ? '.env.production' : '.env.local'
  if (!existsSync(resolve(process.cwd(), modeEnvFile))) {
    throw new Error(
      `${modeEnvFile} is missing. Copy .env.example to ${modeEnvFile} ` +
      `and set VITE_PB_API before running this command.`
    )
  }

  // Fail fast instead of shipping a bundle whose API base is "undefined".
  if (!env.VITE_PB_API) {
    throw new Error(
      `VITE_PB_API is not set for mode "${mode}". ` +
      `Set it in ${modeEnvFile} (see .env.example).`
    )
  }

  return {
    plugins: [
      vue(),
      VitePWA({
        registerType: 'autoUpdate',
        includeAssets: ['icons/favicon.ico', 'icons/favicon.svg', 'icons/apple-touch-icon.png'],
        manifest: {
          name: 'Simple Game Tracker',
          short_name: 'SGT',
          theme_color: '#000000',
          background_color: "#181818",
          icons: [
            {
              src: 'icons/pwa-192x192.png',
              sizes: '192x192',
              type: 'image/png'
            },
            {
              src: 'icons/pwa-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'any'
            },
            {
              src: 'icons/pwa-maskable-192x192.png',
              sizes: '192x192',
              type: 'image/png',
              purpose: 'maskable'
            },
            {
              src: 'icons/pwa-maskable-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'maskable'
            }
          ]
        },
        devOptions: {
          enabled: true
        }
      })
    ],
    define: {
      __APP_ENV__: JSON.stringify(env.APP_ENV),
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
    },
    base: '/',
  }
})
