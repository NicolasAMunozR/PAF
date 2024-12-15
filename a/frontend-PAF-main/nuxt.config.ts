// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],
  runtimeConfig: {
    public: {
      apiBase: 'http://localhost:3000' // Cambia esta URL según tu entorno
    }
  },
  devServer: {
    port: 3001
  },
  plugins: [
    '~/plugins/router.client.js',
  ],
})

