import tailwindcss from "@tailwindcss/vite";
import Aura from '@primeuix/themes/aura';

export default defineNuxtConfig({
  compatibilityDate: '2025-05-15',
  devtools: { enabled: false },
  runtimeConfig: {
    backendUrl: 'http://127.0.0.1:5000',
  },
  css: ['~/assets/css/main.css'],
  modules: [ '@primevue/nuxt-module' ],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  primevue: {
    options: {
      theme: {
        preset: Aura
      }
    }
  },
  routeRules: {
    '/**': {
      headers: {
        'x-content-type-options': 'nosniff',
        'referrer-policy': 'strict-origin-when-cross-origin',
        'x-frame-options': 'DENY',
        'permissions-policy': 'camera=(), microphone=(), geolocation=()',
      },
    },
  },
})
