<template>
  <div class="flex">
    <!-- Tabla de horarios -->
    <div class="w-2/3">
      <h1>Horario para el RUN: {{ run }}</h1>

      <div v-if="persona">
        <table class="w-full text-sm bg-white divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-4 py-2 font-medium text-gray-900">Módulo</th>
              <th v-for="dia in dias" :key="dia" class="px-4 py-2 font-medium text-gray-900">{{ dia }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(horario, index) in horarios" :key="index">
              <td class="px-4 py-2 text-gray-700">{{ horario.modulo }}</td>
              <td v-for="dia in dias" :key="dia" class="px-4 py-2 text-gray-700">
                {{ obtenerAsignaturaPorDia(dia, horario) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p>Cargando datos...</p>
      </div>
    </div>

    <!-- Cuadro de asignaturas al lado de la tabla -->
    <div class="w-1/3 pl-4">
      <h2 class="font-semibold text-lg">Asignaturas</h2>
      <div v-if="persona.length > 0">
      <div
        v-for="(p, index) in persona"
        :key="index"
        class="p-4 mb-4 bg-gray-100 rounded shadow"
      >
        <p><strong>Código de Asignatura:</strong> {{ p.CodigoAsignatura }}</p>
        <p><strong>Nombre de Asignatura:</strong> {{ p.NombreAsignatura }}</p>
        <p><strong>Horas Semanales:</strong> {{ p.CantidadHoras }}</p>
        <p><strong>Jefatura:</strong> {{ p.NombreUnidadContratante }}</p>
      </div>
    </div>

      <div v-else>
        <p>No se encontraron asignaturas.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import { useNuxtApp } from '#app'

const route = useRoute()
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const run = ref(route.query.run || '')
interface Persona {
  CodigoAsignatura: string;
  NombreAsignatura: string;
  CantidadHoras: number;
  NombreUnidadContratante: string;
}

const persona = ref<Persona[]>([]) // Cambiar de 'Persona | null' a 'Persona[]'

const dias = ['Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado']
const horarios = ref([
  { modulo: '08:15 - 09:35' },
  { modulo: '09:50 - 11:10' },
  { modulo: '11:25 - 12:45' },
  { modulo: '13:45 - 15:05' },
  { modulo: '15:20 - 16:40' },
  { modulo: '16:55 - 18:15' },
  { modulo: '18:45 - 20:05' },
  { modulo: '20:05 - 21:25' },
  { modulo: '21:25 - 22:45' }
])

// Función para obtener datos de la persona usando el RUN
const obtenerDatosPersona = async () => {
  try {
    const response = await $axios.get(`/pipelsoft/personas/rut/${run.value}`)
    persona.value = response.data // Ahora es un array de personas
    console.log('Datos de las personas:', persona.value)
  } catch (error) {
    console.error('Error al obtener los datos de las personas:', error)
  }
}
// Función para obtener la asignatura según el día y el horario
const obtenerAsignaturaPorDia = (dia: string, horario: any) => {
  return 'Asignatura ejemplo' // Esto debe cambiarse según los datos obtenidos
}

onMounted(() => {
  obtenerDatosPersona()
})
</script>

<style scoped>
table {
  width: 100%;
  table-layout: auto;
}
</style>
