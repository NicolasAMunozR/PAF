<template>
    <div class="container">
      <!-- Advertencias -->
      <div v-if="warnings.length" class="warnings">
        <h3>Advertencias:</h3>
        <ul>
          <li v-for="warning in warnings" :key="warning">{{ warning }}</li>
        </ul>
      </div>
      
      <!-- Filtros -->
      <Filtros @filter="filterData" @sort="sortData" />
      
      <!-- Tabla -->
      <Tabla :data="filteredPersonas" :showButtons="false" @rowStatusChanged="handleRowStatusChanged" />
    </div>
  </template>
  
  <script setup lang="ts">
  import Filtros from '../components/Filtros.vue';
  import Tabla from '../components/Tabla.vue';
  import { ref, computed, onMounted } from 'vue';
  
  const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };
  
  interface Persona {
    ID: number;
    IdPaf: number;
    CodigoAsignatura: string;
    Run: string;
    EstadoProceso: string;
    Calidad: string;
    Jerarquia: string;
    FechaInicioContrato: string;
    FechaFinContrato: string;
    NombreAsignatura: string;
    cupo: number;
    codigo_modificacion: number;
    bandera_modificacion: number;
    descripcion_modificacion: string;
    seccion: string;
    rowClass?: string;
    SemestrePaf: string;
    DesEstado: string;
    UltimaModificacion: string;
  }
  
  const personas = ref<Persona[]>([]);
  const filtros = ref({
    codigoPAF: '',
    run: '',
    codigoAsignatura: '',
    estadoProceso: '',
    calidad: '',
    jerarquia: '',
    semestre: '',
    nombreAsignatura: '',
    
  });
  const sortBy = ref('nombres');
  const sortOrder = ref('asc');
  const warnings = ref<string[]>([]);
  
  const filteredPersonas = computed(() => {
    warnings.value = [];
    let filtered = personas.value.filter(persona => {
      return (
        persona.NombreAsignatura?.toLowerCase().includes(filtros.value.nombreAsignatura.toLowerCase() || '') &&
        persona.CodigoAsignatura?.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase() || '') &&
        (filtros.value.estadoProceso ? persona.EstadoProceso.toString() === filtros.value.estadoProceso : true) &&
        (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
        persona.IdPaf.toString().toLowerCase().includes(filtros.value.codigoPAF.toLowerCase()) &&
        persona.Run.toLowerCase().includes(filtros.value.run.toLowerCase()) &&
        (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true) &&
        persona.SemestrePaf?.toLowerCase().includes(filtros.value.semestre.toLowerCase() || '')
      );
    });
  
    filtered = filtered.map(persona => {
      if (persona.bandera_modificacion === 1) {
        warnings.value.push(`Advertencia: La PAF ${persona.IdPaf} ha sido modificada.`);
        persona.rowClass = 'modified-row'; // Marca la fila como modificada
      } else if (persona.bandera_modificacion === 2) {
        warnings.value.push(`Advertencia: La PAF ${persona.IdPaf} ha sido eliminada.`);
        persona.rowClass = 'deleted-row'; // Marca la fila como eliminada
      }
  
      if (persona.codigo_modificacion === 1) {
        warnings.value.push(`${persona.descripcion_modificacion}.`);
      }
      return persona;
    });
  
    filtered = filtered.sort((a, b) => {
      const compareA = a[sortBy.value as keyof Persona];
      const compareB = b[sortBy.value as keyof Persona];
  
      if (compareA !== undefined && compareB !== undefined) {
        if (compareA < compareB) return sortOrder.value === 'asc' ? -1 : 1;
        if (compareA > compareB) return sortOrder.value === 'asc' ? 1 : -1;
      }
      return 0;
    });
  
    return filtered;
  });
  
  // Método para manejar la notificación a Tabla.vue
  const handleRowStatusChanged = (persona: Persona) => {
    // Emitir el cambio de estado (modificación o eliminación)
  };
  
  onMounted(async () => {
    try {
      const response = await $axios.get('/historial');
      personas.value = response.data;
    } catch (error) {
      console.error('Error al obtener personas:', error);
    }
  });
  
  const filterData = (newFilters: any) => {
    filtros.value = newFilters;
  };
  
  const sortData = (newSortBy: string, newSortOrder: string) => {
    sortBy.value = newSortBy;
    sortOrder.value = newSortOrder;
  };
  </script>
  
  <style scoped>
  .container {
    display: grid;
    grid-template-columns: 1fr 3fr;
    gap: 1rem;
    font-family: "Helvetica Neue LT", sans-serif;
    background-color: #f9f9f9;
    padding: 20px;
    border: 1px solid #394049;
    border-radius: 8px;
    max-width: 100%;
  }
  
  /* Advertencias */
  .warnings {
    grid-column: 1 / -1;
    background-color: #FFF3CD; /* Fondo amarillo pálido */
    border: 1px solid #FFEBAA; /* Borde amarillo */
    border-radius: 5px;
    padding: 1rem;
    margin-top: 1rem;
    font-family: "Helvetica Neue LT", sans-serif;
    color: #856404; /* Texto amarillo oscuro */
  }
  
  .warnings h3 {
    margin: 0 0 0.5rem;
    font-family: "Bebas Neue Pro", sans-serif;
    color: #EA7600; /* Color institucional */
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
  