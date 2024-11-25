<template>
  <div class="container">
    <Filtros @filter="filterData" @sort="sortData" />
    <!-- Pasa filteredPersonas como prop al componente Tabla -->
    <Tabla :data="filteredPersonas" :showButtons="false" />
    <!-- Advertencia -->
    <div v-if="warnings.length" class="warnings">
      <h3>Advertencias:</h3>
      <ul>
        <li v-for="warning in warnings" :key="warning">{{ warning }}</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import Filtros from '../../components/Filtros.vue'
import Tabla from '../../components/Tabla.vue'
import { ref, computed, onMounted } from 'vue'

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

interface Persona {
  ID: number
  CodigoPAF: string
  CodigoAsignatura: string
  Run: string
  PrimerApellido: string
  SegundoApellido: string
  Correo: string
  EstadoProceso: string
  Calidad: string
  Jerarquia: string
  CantidadHoras: number
  FechaInicioContrato: string
  FechaFinContrato: string
  FechaUltimaModificacionProceso: string
  NombreAsignatura: string
  NombreUnidadContratante: string
  NombreUnidadMayor: string
  cupo: number
  codigo_modificacion: number
  bandera_modificacion: number
  descripcion_modificacion: string
}

const personas = ref<Persona[]>([])
const filtros = ref({
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: '',
  jerarquia: ''
})
const sortBy = ref('nombres')
const sortOrder = ref('asc')
const warnings = ref<string[]>([])

// Computed para datos filtrados y ordenados
const filteredPersonas = computed(() => {
  warnings.value = [] // Resetear advertencias

  let filtered = personas.value.filter(persona => {
    // Filtrado según criterios de búsqueda
    return (
      (persona.CodigoAsignatura?.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase()) ?? false) &&
      (filtros.value.estadoProceso ? persona.EstadoProceso.toString() === filtros.value.estadoProceso : true) &&
      (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
      (persona.CodigoPAF?.toLowerCase().includes(filtros.value.codigoPAF.toLowerCase()) ?? false) &&
      (persona.Run?.toLowerCase().includes(filtros.value.run.toLowerCase()) ?? false) &&
      (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true)
    )
  })

  // Revisar banderas de modificación y generar advertencias
  filtered = filtered.map(persona => {
    if (persona.bandera_modificacion === 1) {
      warnings.value.push(
        `Advertencia: La PAF ${persona.CodigoPAF} ha sido modificado.`
      )
    } else if (persona.bandera_modificacion === 2) {
      warnings.value.push(
        `Advertencia: La PAF ${persona.CodigoPAF} ha sido eliminado.`
      )
    }

    if (persona.codigo_modificacion === 1) {
      warnings.value.push(
        `${persona.descripcion_modificacion}.`
      )
    }

    return persona
  })

  // Ordenamiento
  filtered = filtered.sort((a, b) => {
    const compareA = a[sortBy.value as keyof Persona]
    const compareB = b[sortBy.value as keyof Persona]

    if (compareA !== undefined && compareB !== undefined) {
      if (compareA < compareB) return sortOrder.value === 'asc' ? -1 : 1
      if (compareA > compareB) return sortOrder.value === 'asc' ? 1 : -1
    }
    return 0
  })

  return filtered
})

// onMounted no se modifica
onMounted(async () => {
  try {
    const response = await $axios.get('/historial')
    personas.value = response.data
    console.log('Personas transformadas:', personas.value)
  } catch (error) {
    console.error('Error al obtener personas:', error)
  }
})

const filterData = (newFilters: any) => {
  filtros.value = newFilters
}

const sortData = (newSortBy: string, newSortOrder: string) => {
  sortBy.value = newSortBy
  sortOrder.value = newSortOrder
}
</script>

<style scoped>
.container {
  display: grid;
  grid-template-columns: 1fr 3fr;
  gap: 1rem;
}

.warnings {
  grid-column: 1 / -1;
  background-color: #fff3cd;
  border: 1px solid #ffeeba;
  border-radius: 5px;
  padding: 1rem;
  margin-top: 1rem;
}

.warnings h3 {
  margin: 0 0 0.5rem;
  color: #856404;
}

.warnings ul {
  margin: 0;
  padding: 0;
  list-style-type: none;
}

.warnings li {
  margin: 0.5rem 0;
}
</style>
