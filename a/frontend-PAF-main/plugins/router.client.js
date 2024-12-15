// plugins/router.client.js
import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '@/views/LoginPage.vue';
import Personas from '@/views/Personas.vue';

export default defineNuxtPlugin(nuxtApp => {
  // Verificamos si la aplicación se ejecuta en el cliente
  if (process.client) {
    const routes = [
      { path: '/', component: LoginPage },
      { path: '/personas', component: Personas },
    ];

    const router = createRouter({
      history: createWebHistory(),
      routes,
    });

    // Montamos el router en la aplicación Vue solo si no está ya configurado
    nuxtApp.vueApp.use(router);
  }
});
