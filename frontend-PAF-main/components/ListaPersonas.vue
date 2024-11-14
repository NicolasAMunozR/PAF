<template>
  <div class="container">
    <Filtros @filter="filterData" @sort="sortData" />
    <Tabla :data="filteredPersonas" />
  </div>
</template>

<script setup lang="ts">
import Filtros from './Filtros.vue'
import Tabla from './Tabla.vue'
import { ref, computed, onMounted } from 'vue'
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

// Definición completa de la interfaz Persona
interface Persona {
  ID: number
  CodigoPAF: string
  CodigoAsignatura: string
  Run: string
  Nombres: string
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
}

const personas = ref<Persona[]>([])
const filtros = ref({
  nombres: '',
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: '',
  jerarquia: ''
})
const sortBy = ref('nombres')
const sortOrder = ref('asc')

// Computed para datos filtrados y ordenados
const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    return (
      persona.Nombres.toLowerCase().includes(filtros.value.nombres.toLowerCase()) &&
      persona.CodigoAsignatura.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase()) &&
      (filtros.value.estadoProceso ? persona.EstadoProceso === filtros.value.estadoProceso : true) &&
      (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
      persona.CodigoPAF.toLowerCase().includes(filtros.value.codigoPAF.toLowerCase()) &&
      persona.Run.toLowerCase().includes(filtros.value.run.toLowerCase()) &&
      (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true)
    )
  })

  // Ordenamiento
  filtered = filtered.sort((a, b) => {
    const compareA = a[sortBy.value as keyof Persona]
    const compareB = b[sortBy.value as keyof Persona]

    if (compareA < compareB) return sortOrder.value === 'asc' ? -1 : 1
    if (compareA > compareB) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })

  return filtered
})

onMounted(async () => {
  try {
    const response = await $axios.get('/pipelsoft/persona')
    personas.value = response.data
    console.log('Personas:', personas.value)
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
</style>
