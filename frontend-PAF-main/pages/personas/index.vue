<template>
  <div class="container">
    <Filtros @filter="filterData" @sort="sortData" />
    <Tabla :data="filteredPersonas" />
  </div>
</template>

<script setup lang="ts">
import Filtros from '../../components/Filtros.vue'
import Tabla from '../../components/Tabla.vue'
import { ref, computed, onMounted } from 'vue'

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

interface Persona {
  ID: number;
  CodigoPAF: number;
  CodigoAsignatura: string;
  Run: string;
  Nombres: string;
  PrimerApellido: string;
  SegundoApellido: string;
  Correo: string;
  EstadoProceso: number;
  Calidad: string;
  Jerarquia: string;
  CantidadHoras: number;
  FechaInicioContrato: string;
  FechaFinContrato: string;
  FechaUltimaModificacionProceso: string;
  NombreAsignatura: string;
  NombreUnidadContratante: string;
  NombreUnidadMayor: string;
  CodigoUnidadContratante?: string;
  Bloque?: string;
  Seccion?: string;
  Semestre?: string;
  Cupo?: number;
}

const personas = ref<Persona[]>([]);
const filtros = ref({
  nombres: '',
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: '',
  jerarquia: '',
  nombreAsignatura: '',
  fechaUltimaModificacionProceso: ''
});
const sortBy = ref('nombres');
const sortOrder = ref('asc');

const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    console.log('persona:', persona);
    console.log('filtros:', filtros.value);
    return (
      persona.NombreAsignatura?.toLowerCase().includes(filtros.value.nombreAsignatura.toLowerCase() || '') &&
      persona.CodigoAsignatura?.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase() || '') &&
      (filtros.value.estadoProceso ? persona.EstadoProceso?.toString() === filtros.value.estadoProceso : true) &&
      (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
      persona.CodigoPAF?.toString().toLowerCase().includes(filtros.value.codigoPAF.toString().toLowerCase() || '') &&
      persona.Run?.toLowerCase().includes(filtros.value.run.toLowerCase() || '') &&
      (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true) &&
      persona.FechaUltimaModificacionProceso?.toLowerCase().includes(filtros.value.fechaUltimaModificacionProceso.toLowerCase() || '')
    );
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

onMounted(async () => {
  try {
    const response = await $axios.get('/pipelsoft/contratos');
    console.log('response:', response);
    personas.value = response.data.map((item: any) => ({
      CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
      Nombres: item.PipelsoftData.Nombres,
      PrimerApellido: item.PipelsoftData.PrimerApp,
      SegundoApellido: item.PipelsoftData.SegundoApp,
      CodigoPAF: item.PipelsoftData.IdPaf,
      Calidad: item.PipelsoftData.Calidad,
      Jerarquia: item.PipelsoftData.Jerarquia,
      EstadoProceso: item.PipelsoftData.CodEstado,
      Run: item.PipelsoftData.RunEmpleado,
      Cupo: item.HistorialPafData.cupo,
      NombreAsignatura: item.PipelsoftData.NombreAsignatura,
      FechaUltimaModificacionProceso: item.PipelsoftData.FechaUltimaModificacionProceso,
    }));
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
  max-width: 100%;
}
</style>