export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],
  runtimeConfig: {
    public: {
      apiBase: 'http://localhost:3000' // 🔹 Ajusta si usas Docker
    }
  },
  app: {
    baseURL: '/' // 🔹 No uses '/paf-en-linea' si no es necesario
  },
  routeRules: {
    '/**': { prerender: true } // 🔹 Fuerza la generación estática de todas las rutas
  },
  nitro: {
    output: {
      publicDir: 'dist' // 🔹 Reemplaza generate.dir con esto
    }
  }
})
