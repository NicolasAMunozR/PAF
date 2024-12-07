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
  EstadoProceso: string;
  Calidad: string;
  Jerarquia: string;
  FechaInicioContrato: string;
  FechaFinContrato: string;
  FechaUltimaModificacionProceso: string;
  NombreAsignatura: string;
  NombreUnidadContratante: string;
  NombreUnidadMayor: string;
  CodigoUnidadContratante?: string;
  Bloque?: string;
  seccion?: string;
  Semestre?: string;
  Cupo?: number;
  Id: number;
  SemestrePaf: string;
  DesEstado: string;
  nombreUnidadMenor: string;
  nombreUnidadMayor: string;
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
  fechaUltimaModificacionProceso: '',
  id: '',
  semestre: '',
  nombreUnidadMenor: '',
  nombreUnidadMayor: '',
});
const sortBy = ref('nombres');
const sortOrder = ref('asc');

const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    return (
      persona.NombreAsignatura?.toLowerCase().includes(filtros.value.nombreAsignatura.toLowerCase() || '') &&
      persona.CodigoAsignatura?.toLowerCase().includes(filtros.value.codigoAsignatura.toLowerCase() || '') &&
      (filtros.value.estadoProceso ? persona.EstadoProceso?.toString() === filtros.value.estadoProceso : true) &&
      (filtros.value.calidad ? persona.Calidad === filtros.value.calidad : true) &&
      persona.CodigoPAF?.toString().toLowerCase().includes(filtros.value.codigoPAF.toString().toLowerCase() || '') &&
      persona.Run?.toLowerCase().includes(filtros.value.run.toLowerCase() || '') &&
      (filtros.value.jerarquia ? persona.Jerarquia === filtros.value.jerarquia : true) &&
      persona.FechaUltimaModificacionProceso?.toLowerCase().includes(filtros.value.fechaUltimaModificacionProceso.toLowerCase() || '') &&
      persona.SemestrePaf?.toLowerCase().includes(filtros.value.semestre.toLowerCase() || '') &&
      persona.nombreUnidadMenor?.toLowerCase().includes(filtros.value.nombreUnidadMenor.toLowerCase() || '') &&
      persona.nombreUnidadMayor?.toLowerCase().includes(filtros.value.nombreUnidadMayor.toLowerCase() || '')
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
    personas.value = response.data.map((item: any) => {
      const bloquesArray = item.HistorialPafData.Bloque || []; // Asegurar que Bloque sea un arreglo (vacío si es null o undefined)

      // Verificar si el arreglo no está vacío antes de hacer el map
      const bloque = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.bloques).join("/") : "";
      const CodigoA = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.codigoAsignatura).join("/") : "";
      const Cupo = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.cupos).join("/") : "";
      const seccion = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.seccion).join("/") : "";

      return {
      CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
      Nombres: item.PipelsoftData.Nombres,
      PrimerApellido: item.PipelsoftData.PrimerApp,
      SegundoApellido: item.PipelsoftData.SegundoApp,
      CodigoPAF: item.PipelsoftData.IdPaf,
      Calidad: item.PipelsoftData.Categoria,
      Jerarquia: item.PipelsoftData.Jerarquia,
      EstadoProceso: item.PipelsoftData.CodEstado,
      DesEstado: item.PipelsoftData.DesEstado,
      SemestrePaf: item.PipelsoftData.Semestre,
      Run: item.PipelsoftData.RunEmpleado,
      Cupo,
      NombreAsignatura: item.PipelsoftData.NombreAsignatura,
      FechaUltimaModificacionProceso: item.PipelsoftData.UltimaModificacion,
      Id: item.PipelsoftData.Id,
      seccion,
      nombreUnidadMayor: item.PipelsoftData.NombreUnidadMayor,
      nombreUnidadMenor: item.PipelsoftData.NombreUnidadMenor,
    };
  });
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