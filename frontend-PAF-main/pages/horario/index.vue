<template>
  <div class="flex">
    <!-- Botón para volver -->
    <div class="absolute top-4 left-4">
      <button @click="volver" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Volver
      </button>
    </div>

    <!-- Tabla de horarios -->
    <div class="w-2/3 mt-12">
      <h1>Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}</h1>

      <div v-if="persona.length > 0">
        <div class="mb-4">
          <label for="semestre">Seleccionar Semestre:</label>
          <select id="semestre" v-model="semestreSeleccionado" class="ml-2">
            <option v-for="sem in semestres" :key="sem" :value="sem">{{ sem }}</option>
          </select>
        </div>

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
                <div v-for="bloque in bloquesPorDia(dia, index + 1)" :key="bloque.nombre" :style="{ backgroundColor: bloque.color }" class="rounded p-1 text-black">
                  {{ bloque.nombre }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p>Cargando datos o no se encontraron registros para el RUN.</p>
      </div>
    </div>

    <!-- Cuadro lateral -->
    <div class="w-1/3 pl-4 mt-12">
      <h2 class="font-semibold text-lg">PAF</h2>
      <div v-if="persona.length > 0">
        <div
          v-for="(p, index) in personaFiltrada"
          :key="index"
          :style="{ backgroundColor: colores[index % colores.length] }"
          class="p-4 mb-4 rounded shadow"
        >
          <p><strong>Código de Asignatura:</strong> {{ p.CodigoAsignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.NombreAsignatura }}</p>
          <p><strong>Horas Semanales:</strong> {{ p.CantidadHoras }}</p>
          <p><strong>Jefatura:</strong> {{ p.NombreUnidadContratante }}</p>
          <p><strong>Bloque:</strong> {{ p.Bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.Cupo }}</p>
          <p><strong>Sección:</strong> {{ p.Seccion }}</p>
          <p><strong>Semestre:</strong> {{ p.Semestre }}</p>
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
import { onMounted, ref, computed } from 'vue'
import { useNuxtApp } from '#app'

const route = useRoute()
const router = useRouter()

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
  Dia?: string;
  Bloque?: string;
  Cupo?: number;
  Seccion?: string;
  Semestre?: string;
}

const persona = ref<Persona[]>([])
const colores = ['#FFCDD2', '#F8BBD0', '#E1BEE7', '#D1C4E9', '#C5CAE9']
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

const semestreSeleccionado = ref('')
const semestres = computed(() => [...new Set(persona.value.map(p => p.Semestre))])
const personaFiltrada = computed(() => persona.value.filter(p => p.Semestre === semestreSeleccionado.value))

const bloquesPorDia = (dia: string, modulo: number) => {
  // Mapear iniciales de días con sus nombres completos
  const inicialDia: { [key: string]: string } = {
    Lunes: 'L',
    Martes: 'M',
    Miércoles: 'W', // 'W' para distinguir de miércoles
    Jueves: 'J',
    Viernes: 'V',
    Sábado: 'S',
  };

  // Filtrar asignaturas para el día específico y el módulo actual
  return personaFiltrada.value
    .filter((p) => {
      if (!p.Bloque) return false;

      // Separar bloques (por ejemplo: M2-M5-V1 se divide en ["M2", "M5", "V1"])
      const bloques = p.Bloque.split('-');

      // Buscar un bloque que coincida con el día y el módulo
      return bloques.some((b) => {
        const diaBloque = b.charAt(0); // Primera letra es el día (M, V, etc.)
        const moduloBloque = b.slice(1); // Resto es el número de módulo

        // Comparar con el día actual y el módulo actual
        return inicialDia[dia] === diaBloque && parseInt(moduloBloque) === modulo;
      });
    })
    .map((p) => ({
      nombre: p.NombreAsignatura,
      seccion: p.Seccion,
      color: colores[persona.value.indexOf(p) % colores.length],
    }));
};


const obtenerDatosPersona = async () => {
  try {
    const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`)
    persona.value = response.data.map((item: any) => ({
      CodigoAsignatura: item.pipelsoft_data.CodigoAsignatura,
      Nombres: item.pipelsoft_data.Nombres,
      NombreAsignatura: item.pipelsoft_data.NombreAsignatura,
      PrimerApellido: item.pipelsoft_data.PrimerApellido,
      SegundoApellido: item.pipelsoft_data.SegundoApellido,
      CantidadHoras: item.pipelsoft_data.CantidadHoras,
      NombreUnidadContratante: item.pipelsoft_data.NombreUnidadContratante,
      Dia: item.profesor_data.dia,
      Bloque: item.profesor_data.bloque,
      Cupo: item.profesor_data.cupo,
      Seccion: item.profesor_data.seccion,
      Semestre: item.profesor_data.semestre,
    }))
  } catch (error) {
    console.error('Error al obtener los datos:', error)
  }
}

const volver = () => {
  router.go(-1)
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
