<template>
  <div>
    <Tabla :data="personas" />
  </div>
</template>

<script setup lang="ts">
import Tabla from './Tabla.vue'
import { ref, onMounted } from 'vue'
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

interface Persona {
  id: number
  nombre: string
  correo: string
}

const personas = ref<Persona[]>([])

onMounted(async () => {
  try {
    const response = await $axios.get('/personas')
    personas.value = response.data
  } catch (error) {
    console.error('Error al obtener personas:', error)
  }
})
</script>
