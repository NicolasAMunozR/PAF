import { defineNuxtPlugin } from '#app'
import { createRouter, createWebHistory } from 'vue-router'

export default defineNuxtPlugin((nuxtApp) => {
  const routes = [
    { path: '/', component: () => import('~/views/LoginPage.vue') },
    { path: '/seguimientoPAF', component: () => import('~/views/SeguimientoPAF.vue') },
    { path: '/unidadMayorPAF', component: () => import('~/views/UnidadMayorPAF.vue') },
    { path: '/historyPAF', component: () => import('~/views/HistoryPAF.vue') },
    { path: '/personas', component: () => import('~/views/Personas.vue') },
    { path: '/profesorPAF', component: () => import('~/views/ProfesorPAF.vue') },
    { path: '/paf', component: () => import('~/views/PAF.vue') },
    { path: '/horario', component: () => import('~/views/Horario.vue') },
    { path: '/estadisticaUnidadPAF', component: () => import('~/views/EstadisticaUnidadPAF.vue') },
    { path: '/estadisticaPAF', component: () => import('~/views/EstadisticaPAF.vue') }
  ]

  const router = createRouter({
    history: createWebHistory(),
    routes
  })

  nuxtApp.vueApp.use(router)
})