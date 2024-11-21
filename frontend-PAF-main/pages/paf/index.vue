<template>
  <div class="flex flex-col">
    <!-- Botón para volver debajo de la barra superior pero encima de la información de la PAF -->
    <div class="mt-4 ml-4">
      <button @click="volver" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Volver
      </button>
    </div>

    <!-- Información de la PAF -->
    <div class="w-2/3 mt-6">
      <h1 v-if="paf.length > 0" class="mb-4">Información de la PAF:</h1>

      <!-- Mostrar la lista de personas -->
      <div v-if="paf.length > 0">
        <div
          v-for="persona in paf"
          :key="persona.pipelsoft_data.CodigoPAF"
          class="paf-container mb-6 border rounded p-4 shadow-md"
        >
          <p><strong>Código PAF:</strong> {{ persona.pipelsoft_data.CodigoPAF }}</p>
          <p><strong>RUN:</strong> {{ persona.pipelsoft_data.Run }}</p>
          <p><strong>Nombre:</strong> {{ persona.pipelsoft_data.Nombres }} {{ persona.pipelsoft_data.PrimerApellido }} {{ persona.pipelsoft_data.SegundoApellido }}</p>
          <p><strong>Correo:</strong> {{ persona.pipelsoft_data.Correo }}</p>
          <p><strong>Asignatura:</strong> {{ persona.profesor_data.nombre_asignatura }}</p>
          <p><strong>Bloque:</strong> {{ persona.profesor_data.bloque }}</p>
          <p><strong>Unidad Contratante:</strong> {{ persona.pipelsoft_data.NombreUnidadContratante }}</p>

          <!-- Botón ubicado en la parte inferior -->
          <div class="flex justify-end mt-4">
            <button @click="dejarListaPaf(persona.pipelsoft_data.CodigoPAF)" class="bg-blue-500 text-white py-2 px-4 rounded">
              Dejar lista la PAF
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <p>Cargando datos o no se encontraron registros para la PAF.</p>
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
  pipelsoft_data: {
    CodigoPAF: string;
    Run: string;
    Nombres: string;
    PrimerApellido: string;
    SegundoApellido: string;
    Correo: string;
    NombreUnidadContratante: string;
  };
  profesor_data: {
    nombre_asignatura: string;
    bloque: string;
  };
}

const paf = ref<Paf[]>([])

const obtenerDatosPaf = async () => {
  try {
    console.log('Obteniendo datos de la PAF con código:', codigoPaf.value)
    if (!codigoPaf.value) {
      console.error('Código PAF no está disponible')
      return
    }
    const response = await $axios.get(`/contratos/codigo_paf/${codigoPaf.value}`)
    if (response.data) {
      paf.value = Array.isArray(response.data) ? response.data : [response.data] 
      console.log('Datos de las PAFs:', paf.value)
    } else {
      console.log('No se encontraron datos para el código PAF:', codigoPaf.value)
    }
  } catch (error) {
    console.error('Error al obtener los datos de la PAF:', error)
  }
}

const volver = () => {
  router.go(-1)
}

const dejarListaPaf = async (codigoPaf: string) => {
  try {
    const postResponse = await $axios.post('/historial', { Codigo_paf: codigoPaf })
    console.log('Historial creado:', postResponse.data)
    
    router.push('/personas')
  } catch (error) {
    console.error('Error al procesar la PAF:', error)
  }
}

onMounted(() => {
  obtenerDatosPaf()
})
</script>

<style scoped>
.paf-container {
  display: flex;
  flex-direction: column;
}

button {
  cursor: pointer;
}

.flex.justify-end {
  margin-top: auto;
}

.mt-4.ml-4 {
  margin-top: 16px;
  margin-left: 16px;
}
</style>
