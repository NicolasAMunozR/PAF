import axios from 'axios'

export default defineNuxtPlugin(nuxtApp => {
  const axiosInstance = axios.create({
    baseURL: 'http://localhost:3000'  // Cambia esta URL según la dirección de tu backend Go
  })

  // Hacer que Axios esté disponible globalmente en la aplicación
  nuxtApp.provide('axios', axiosInstance)
})
