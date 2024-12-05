<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-primary-color text-white">
        <tr>
          <th class="px-4 py-3 text-left font-semibold">Código de la PAF</th>
          <th class="px-4 py-3 text-left font-semibold">Código de la Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Run</th>
          <th class="px-4 py-3 text-left font-semibold">Nombre de Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Estado de Proceso</th>
          <th class="px-4 py-3 text-left font-semibold">Sección</th>
          <th class="px-4 py-3 text-left font-semibold">Cupos</th>
          <th class="px-4 py-3 text-left font-semibold">Semestre de PAF</th>
          <th v-if="showButtons" class="px-4 py-3 text-left font-semibold">Opciones</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="persona in paginatedData" :key="persona.Id" :class="persona.rowClass" class="hover:bg-gray-50 transition-colors">
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoPAF }} {{ persona.IdPaf }}</td>
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.NombreAsignatura }} {{ persona.nombre_asignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.EstadoProceso }} {{ persona.estado_proceso }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.seccion }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Cupo }} {{ persona.cupo }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Semestre }} {{ persona.SemestrePaf }}</td>
          <td v-if="showButtons" class="px-4 py-3">
            <a :href="`/paf?codigoPaf=${persona.CodigoPAF}`" class="button">Ver PAF</a>
            <br>
            <br>
            <a :href="`/horario?run=${persona.Run}`" class="button">Ver Horarios</a>
          </td>
        </tr>
      </tbody>
    </table>
    <!-- Paginación -->
    <div class="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
      <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            Showing
            {{ ' ' }}
            <span class="font-medium">{{ startIndex + 1 }}</span>
            {{ ' ' }}
            to
            {{ ' ' }}
            <span class="font-medium">{{ endIndex }}</span>
            {{ ' ' }}
            of
            {{ ' ' }}
            <span class="font-medium">{{ data.length }}</span>
            {{ ' ' }}
            results
          </p>
        </div>
        <div>
          <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
            <button 
              :disabled="currentPage === 1"
              @click="changePage(currentPage - 1)"
              class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
            >
              <span class="sr-only">Previous</span>
              <ChevronLeftIcon class="size-5" aria-hidden="true" />
            </button>
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="changePage(page)"
              :class="['relative inline-flex items-center px-4 py-2 text-sm font-semibold', page === currentPage ? 'bg-indigo-600 text-white' : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50']"
            >
              {{ page }}
            </button>
            <button 
              :disabled="currentPage === totalPages"
              @click="changePage(currentPage + 1)"
              class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
            >
              <span class="sr-only">Next</span>
              <ChevronRightIcon class="size-5" aria-hidden="true" />
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup>
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/20/solid'

const props = defineProps({
  data: {
    type: Array,
    required: true
  },
  showButtons: {
    type: Boolean,
    default: true
  }
})

const ITEMS_PER_PAGE = 10
const totalPages = computed(() => Math.ceil(props.data.length / ITEMS_PER_PAGE))
const currentPage = ref(1)

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * ITEMS_PER_PAGE
  const end = start + ITEMS_PER_PAGE
  return props.data.slice(start, end)
})

const startIndex = computed(() => (currentPage.value - 1) * ITEMS_PER_PAGE)
const endIndex = computed(() => Math.min(startIndex.value + ITEMS_PER_PAGE, props.data.length))

const visiblePages = computed(() => {
  const range = []
  const maxVisible = 5
  const halfVisible = Math.floor(maxVisible / 2)
  
  if (totalPages.value <= maxVisible) {
    for (let i = 1; i <= totalPages.value; i++) range.push(i)
  } else if (currentPage.value <= halfVisible) {
    for (let i = 1; i <= maxVisible; i++) range.push(i)
    range.push('...')
  } else if (currentPage.value > totalPages.value - halfVisible) {
    range.push('...')
    for (let i = totalPages.value - maxVisible + 1; i <= totalPages.value; i++) range.push(i)
  } else {
    range.push('...')
    for (let i = currentPage.value - halfVisible; i <= currentPage.value + halfVisible; i++) range.push(i)
    range.push('...')
  }
  return range
})

const changePage = (page) => {
  if (page !== '...' && page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}
</script>



<style scoped>
/* Colores institucionales */
:root {
  --primary-color: #EA7600; /* Color principal USACH */
  --secondary-color: #394049; /* Color secundario USACH */
  --accent-color: #C8102E; /* Complementario */
  --background-color: #1d558d; /* Fondo neutro */
  --button-background-color: #4CAF50; /* Color de fondo de los botones */
  --button-hover-color: #388E3C; /* Color de hover en los botones */
}

/* Fila modificada y eliminada */
.modified-row {
  background-color: yellow; /* Amarillo suave */
}

.deleted-row {
  background-color: red; /* Rojo suave */
}

/* Contenedor de la tabla */
.table-container {
  width: 100%;
  padding: 20px;
  overflow-x: auto;
  background-color: var(--background-color);
}

/* Estilos para la tabla */
table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #030101;
  background-color: var(--primary-color);
}

tbody td {
  padding: 12px;
  font-size: 0.875rem;
  color: var(--secondary-color);
}

/* Estilos para los botones */
.button {
  display: inline-block;
  padding: 8px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  color: #090000;
  background-color: #4CAF50;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.button:hover {
  background-color: var(--button-hover-color);
}

.hover\:bg-gray-50:hover {
  background-color: #f1f5f9;
}
</style>
