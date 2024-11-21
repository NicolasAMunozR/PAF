<template>
      <div class="container">
    <Filtros @filter="filterData" @sort="sortData" />
    <!-- Pasa filteredPersonas como prop al componente Tabla, y oculta los botones -->
    <Tabla :data="filteredPersonas" :showButtons="false" />
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
  
  // Computed para datos filtrados y ordenados
  const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    console.log(persona) // Verifica si los datos son correctos
    return (
      (persona.CodigoAsignatura?.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase()) ?? false) &&
      (filtros.value.estadoProceso ? persona.EstadoProceso.toString() === filtros.value.estadoProceso : true) &&
      (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
      (persona.CodigoPAF?.toLowerCase().includes(filtros.value.codigoPAF.toLowerCase()) ?? false) &&
      (persona.Run?.toLowerCase().includes(filtros.value.run.toLowerCase()) ?? false) &&
      (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true)
    )
  })

  console.log('Filtrados:', filtered) // Verifica los datos filtrados
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
  </style>
  