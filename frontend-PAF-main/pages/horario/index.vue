<template>
  <div class="flex">
    <!-- Botón para volver en la esquina superior izquierda -->
    <div class="absolute top-4 left-4">
      <button @click="volver" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Volver
      </button>
    </div>

    <!-- Tabla de horarios -->
    <div class="w-2/3 mt-12">
      <h1>Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}</h1>

      <div v-if="persona.length > 0">
        <table class="w-full text-sm bg-white divide-y divide-gray-300 border">
          <thead>
            <tr class="bg-gray-200">
              <th class="px-4 py-2 font-medium text-gray-900 border-r">Módulo</th>
              <th v-for="dia in dias" :key="dia" class="px-4 py-2 font-medium text-gray-900 border-r">{{ dia }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(horario, index) in horarios" :key="index" :class="index % 2 === 0 ? 'bg-gray-100' : 'bg-white'">
              <td class="px-4 py-2 text-gray-700 border-r">{{ horario.modulo }}</td>
              <td v-for="dia in dias" :key="dia" class="px-4 py-2 text-gray-700 border-r">
                {{ obtenerAsignaturaPorDia(dia, horario.modulo) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p>Cargando datos o no se encontraron registros para el RUN.</p>
      </div>
    </div>

    <!-- Cuadro de asignaturas al lado de la tabla -->
    <div class="w-1/3 pl-4 mt-12">
      <h2 class="font-semibold text-lg">Asignaturas</h2>
      <div v-if="persona.length > 0">
        <div
          v-for="(p, index) in persona"
          :key="index"
          :style="{ backgroundColor: colores[index % colores.length] }"
          class="p-4 mb-4 rounded shadow"
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
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { useNuxtApp } from '#app'

const route = useRoute()
const router = useRouter() // Obtén el enrutador

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const run = ref(route.query.run || '')

interface Persona {
  CodigoAsignatura: string;
  NombreAsignatura: string;
  CantidadHoras: number;
  NombreUnidadContratante: string;
  Nombres: string;
  PrimerApellido: string;
  SegundoApellido: string;
}

const persona = ref<Persona[]>([]) // Cambiar de 'Persona | null' a 'Persona[]'

// Colores para los cuadros de asignaturas
const colores = ['#FFCDD2', '#F8BBD0', '#E1BEE7', '#D1C4E9', '#C5CAE9', '#BBDEFB', '#B3E5FC', '#B2EBF2', '#B2DFDB', '#C8E6C9']

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
  return '' // Esto debe cambiarse según los datos obtenidos
}

// Función para volver a la página anterior
const volver = () => {
  router.go(-1) // Esto hace que el usuario regrese a la página anterior
}

onMounted(() => {
  obtenerDatosPersona()
})
</script>

<style scoped>
table {
  width: 100%;
  table-layout: auto;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid #d1d5db; /* Borde para diferenciar módulos */
}

button {
  position: absolute;
  top: 4rem; /* Ajusta la posición según lo necesario */
  left: 1rem; /* Ajusta la distancia desde la izquierda */
  z-index: 10; /* Asegura que el botón esté por encima de otros elementos */
}
</style>
