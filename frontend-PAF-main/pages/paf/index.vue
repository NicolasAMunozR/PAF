<template>
  <div class="flex">
    <!-- Botón para volver en la esquina superior izquierda -->
    <div class="absolute top-4 left-4">
      <button @click="volver" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Volver
      </button>
    </div>

    <!-- Información de la PAF -->
    <div class="w-2/3 mt-12">
      <h1 v-if="paf.length > 0">Información de las PAFs:</h1>

      <!-- Mostrar la lista de personas -->
      <div v-if="paf.length > 0">
        <div v-for="persona in paf" :key="persona.CodigoPaf">
          <p><strong>Código PAF:</strong> {{ persona.CodigoPaf }}</p>
          <p><strong>RUN:</strong> {{ persona.Run }}</p>
          <p><strong>Nombre:</strong> {{ persona.Nombres }} {{ persona.PrimerApellido }} {{ persona.SegundoApellido }}</p>
          <p><strong>Correo:</strong> {{ persona.Correo }}</p>
          <div class="action-buttons">
            <button @click="dejarListaPaf(persona.CodigoPaf)" class="bg-blue-500 text-white py-2 px-4 rounded mr-4">Dejar lista la PAF</button>
          </div>
          <hr />
        </div>
      </div>

      <div v-else>
        <p>Cargando datos o no se encontraron registros para las PAFs.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { useNuxtApp } from '#app'

const route = useRoute()
const router = useRouter()

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const codigoPaf = ref(route.query.codigoPaf || '')

console.log('Código PAF desde los parámetros de la ruta:', route.params.codigoPaf)

interface Paf {
  CodigoPaf: string;
  Run: string;
  Nombres: string;
  PrimerApellido: string;
  SegundoApellido: string;
  Correo: string;
}

const paf = ref<Paf[]>([])  // Cambiado de null a un array vacío

const obtenerDatosPaf = async () => {
  try {
    console.log('Obteniendo datos de la PAF con código:', codigoPaf.value)
    if (!codigoPaf.value) {
      console.error('Código PAF no está disponible')
      return
    }
    console.log('Obteniendo datos de la PAF con código:', codigoPaf.value)
    const response = await $axios.get(`/pipelsoft/persona/paf/${codigoPaf.value}`)
    if (response.data) {
      paf.value = Array.isArray(response.data) ? response.data : [response.data] // Si no es un arreglo, lo convierte a uno
      console.log('Datos de las PAFs:', paf.value)
    } else {
      console.log('No se encontraron datos para el código PAF:', codigoPaf.value)
    }
  } catch (error) {
    console.error('Error al obtener los datos de la PAF:', error)
  }
}


// Función para volver a la página anterior
const volver = () => {
  router.go(-1) // Esto hace que el usuario regrese a la página anterior
}

// Función para dejar lista la PAF (POST + DELETE)
const dejarListaPaf = async (codigoPaf: string) => {
  try {
    // Realizar el POST para registrar el historial
    const postResponse = await $axios.post('/historial', {
      Codigo_paf: codigoPaf, // Usar la clave esperada por el backend
    })
    console.log('Historial creado:', postResponse.data)

    // Realizar el DELETE después del POST
    const deleteResponse = await $axios.delete(`/historial/${codigoPaf}`)
    console.log('Historial eliminado:', deleteResponse.data)

  } catch (error) {
    console.error('Error al procesar la PAF:', error)
  }
}

// Cargar los datos de la PAF cuando el componente se monta
onMounted(() => {
  obtenerDatosPaf()
})
</script>

<style scoped>
.table-container {
  width: 100%;
  padding: 20px;
  overflow-x: auto;
  background-color: #f9fafb;
}

button {
  position: absolute;
  top: 4rem; /* Ajusta la posición según lo necesario */
  left: 1rem; /* Ajusta la distancia desde la izquierda */
  z-index: 10; /* Asegura que el botón esté por encima de otros elementos */
}
</style>

