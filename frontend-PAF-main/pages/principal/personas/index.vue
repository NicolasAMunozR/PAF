<template>
  <div class="container">
    <div class="filters">
        <!-- Filtro Run -->
        <div class="filter-item">
      <label for="run" class="label">Run</label>
      <input
        v-model="run_filtro"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>

    <!-- Filtro Código de PAF -->
    <div class="filter-item" >
      <label for="codigoPAF" class="label">Código de PAF</label>
      <input
        v-model="codigoPAF_filtro"
        type="text"
        class="input" 
        placeholder="Filtrar por código de PAF"
      />
    </div>

    <!-- Filtro Código de Asignatura -->
    <div class="filter-item" >
      <label for="codigoAsignatura" class="label">Código de Asignatura</label>
      <input
        v-model="codigoAsignatura_filtro"
        type="text"
        class="input"
        placeholder="Filtrar por código de asignatura"
      />
    </div>

    <!-- Filtro Nombre Asignatura -->
    <div class="filter-item" >
      <label for="nombreAsignatura" class="label">Nombre Asignatura</label>
      <select id="nombreAsignatura" v-model="nombreAsignatura_filtro" class="select">
        <option value="">Todos</option>
        <option v-for="asi in nombreAsig" :key="asi" :value="asi">
          {{ asi }}
        </option>
    </select>
    </div>

    <!-- Filtro Nombre Unidad Mayor (solo visible si es seguimientoPAF) -->
    <div class="filter-item" >
      <label for="nombreUnidadMayor" class="label">Nombre Unidad Mayor</label>
      <select id="nombreUnidadMayor" v-model="nombreUnidadMayor_filtro" class="select">
      <option value="">Todos</option>
      <option
        v-for="may in nombreUnidadMay"
        :key="may"
        :value="may"
      >
        {{ may }}
      </option>
    </select>
    </div>
    
    <!-- Filtro Nombre Unidad Menor -->
    <div class="filter-item">
      <label for="nombreUnidadMenor" class="label">Nombre Unidad Menor</label>
      <select id="nombreUnidadMenor" v-model="nombreUnidadMenor_filtro" class="select">
        <option value="">Todos</option>
        <option v-for="men in nombreUnidadMen" :key="men" :value="men">
          {{ men }}
        </option>
      </select>
    </div>

    <!-- Filtro semestre -->
    <div class="filter-item">
    <label for="semestre" class="label">Semestre</label>
    <select id="semestre" v-model="semestres_filtro" class="select">
      <option value="">Todos</option>
      <option
        v-for="sem in semestres"
        :key="sem"
        :value="sem"
      >
        {{ sem }}
      </option>
    </select>
  </div>
    
    <!-- Filtro Estado de Proceso -->
    <div class="filter-item" >
      <label for="estadoProceso" class="label">Estado de Proceso</label>
      <select id="estadoProceso" v-model="estado_filtro" class="select">
        <option value="">Todos</option>
        <option value="A1">Estado: Sin Solicitar</option>
        <option value="A2">Estado: Enviada al Interesado</option>
        <option value="A3">Estado: Enviada al Validador</option>
        <option value="B1">Estado: Aprobada por Validador</option>
        <option value="B9">Estado: Rechazada por Validador</option>
        <option value="C1D">Estado: Aprobada por Dir. Pregrado</option>
        <option value="C9D">Estado: Rechazada por Dir. de Pregrado</option>
        <option value="F1">Estado: Aprobada por RRHH</option>
        <option value="F9">Estado: Rechazada por RRHH</option>
        <option value="A9">Estado: Anulada</option>
      </select> 
    </div>
    
    <!-- Filtro Categoría -->
    <div class="filter-item">
      <label for="calidad" class="label">Categoria</label>
      <select id="calidad" v-model="calidad_filtro" class="select">
        <option value="">Todas</option>
        <option value="PROFESOR HORAS CLASES">Profesor por hora</option>
      </select>
    </div>
    
    <!-- Ordenar por (solo visible si es seguimientoPAF) -->
    <div class="sort-item">
      <label for="sort" class="label">Ordenar por</label>
      <select v-model="sortBy" class="select">
        <option value="NombreAsignatura">Nombre de Asignatura</option>
        <option value="CodigoAsignatura">Código de Asignatura</option>
        <option value="Run">Run</option>
        <option value="CodigoPAF">Código de PAF</option>
        <option value="FechaUltimaModificacionProceso">Última Actualización del Proceso</option>
      </select>
      <button @click="toggleSortOrder" class="btn sort-btn">
        Ordenar {{ sortOrder === 'asc' ? 'Ascendente' : 'Descendente' }}
      </button>
    </div>
    
    <!-- Botón de reset -->
    <div class="filter-item">
      <button @click="resetFilters" class="btn reset-btn">Resetear</button>
    </div>

    <!-- Botón de confirmar filtros -->
    <div class="filter-item">
      <button @click="applyFilters" class="btn confirm-btn">Confirmar Filtros</button>
    </div>
  </div>
  <Tabla :data="filteredPersonas" :showButton="false"/>
</div>
</template>

<script setup lang="ts">
import Tabla from '../../../components/Tabla.vue'
import { ref, computed, onMounted } from 'vue'

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const semestreMasActual = ref('');
const semestres_filtro = ref('');
const codigoPAF_filtro = ref('');
const run_filtro = ref('');
const codigoAsignatura_filtro = ref('');
const calidad_filtro = ref('');
const nombreAsignatura_filtro = ref('');
const nombreUnidadMenor_filtro = ref('');
const nombreUnidadMayor_filtro = ref('');
const estado_filtro = ref('');
const jerarquia_filtro = ref('');

onMounted(async () => {
  try {
    const response = await $axios.get('/api/paf-en-linea/pipelsoft/contratos');
    semestreMasActual.value = response.data
  .map((item: { PipelsoftData: { Semestre: any; }; }) => item.PipelsoftData.Semestre) // Extrae el semestre de cada objeto
  .sort((a: { split: (arg0: string) => [any, any]; }, b: { split: (arg0: string) => [any, any]; }) => {
    // Convierte los semestres en fechas para poder compararlos
    const [mesA, anioA] = a.split('-');
    const [mesB, anioB] = b.split('-');

    // Compara año y mes
    if (parseInt(anioA) === parseInt(anioB)) {
      return parseInt(mesB) - parseInt(mesA); // Ordena por mes descendente
    }
    return parseInt(anioB) - parseInt(anioA); // Ordena por año descendente
  })[0]; // Devuelve el semestre más reciente
    personas.value = response.data.map((item: any) => {
      const bloquesArray = item.HistorialPafData.Bloque || []; // Asegurar que Bloque sea un arreglo (vacío si es null o undefined)

      // Verificar si el arreglo no está vacío antes de hacer el map
      const bloque = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.bloques).join(" / ") : "";
      const CodigoA = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.codigoAsignatura).join(" / ") : "";
      const Cupo = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.cupos).join(" / ") : "";
      const seccion = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.seccion).join(" / ") : "";

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
  if (sessionStorage.getItem('codigoPAF')) {
    codigoPAF_filtro.value = sessionStorage.getItem('codigoPAF') || '';
  }
  if (sessionStorage.getItem('run')) {
    run_filtro.value = sessionStorage.getItem('run') || '';
  }
  if (sessionStorage.getItem('codigoAsignatura')) {
    codigoAsignatura_filtro.value = sessionStorage.getItem('codigoAsignatura') || '';
  }
  if (sessionStorage.getItem('semestre') || sessionStorage.getItem('semestre') == '') {
    semestres_filtro.value = sessionStorage.getItem('semestre') || '';
  }
  else{
    semestres_filtro.value = semestreMasActual.value;
  }
  if (sessionStorage.getItem('estadoProceso') || sessionStorage.getItem('estadoProceso') == '') {
    estado_filtro.value = sessionStorage.getItem('estadoProceso') || '';
  }
  else{
    estado_filtro.value = 'B1';
  }
  if (sessionStorage.getItem('calidad')) {
    calidad_filtro.value = sessionStorage.getItem('calidad') || '';
  }
  if (sessionStorage.getItem('nombreAsignatura')) {
    nombreAsignatura_filtro.value = sessionStorage.getItem('nombreAsignatura') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMenor')) {
    nombreUnidadMenor_filtro.value = sessionStorage.getItem('nombreUnidadMenor') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMayor')) {
    nombreUnidadMayor_filtro.value = sessionStorage.getItem('nombreUnidadMayor') || '';
  }
  applyFilters();
  } catch (error) {
    console.error('Error al obtener personas:', error);
  }
});

const filtros = ref({
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  semestre: semestreMasActual,
  estadoProceso: '',
  calidad: '',
  nombreAsignatura: '',
  fechaUltimaModificacionProceso: '',
  nombreUnidadMenor: '',
  nombreUnidadMayor: '',
  ruta: '',
  jerarquia: '',
});
const filtros1 = ref({
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  semestre: '',
  estadoProceso: '',
  calidad: '',
  nombreAsignatura: '',
  fechaUltimaModificacionProceso: '',
  nombreUnidadMenor: '',
  nombreUnidadMayor: '',
  ruta: '',
  jerarquia: '',
});

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
// Removed duplicate declaration of filtros
const semestres = computed(() => {
  // Ordenar los semestres y eliminar duplicados
  const semestresOrdenados = [...new Set(personas.value.map(p => p.SemestrePaf))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return semestresOrdenados; // O retorna [ultimoSemestre] si solo necesitas el último
});

const nombreUnidadMay = computed(() => {
  const nombreUnidadMayor = [...new Set(personas.value.map(p => p.nombreUnidadMayor))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return nombreUnidadMayor; // O retorna [ultimoSemestre] si
});

const nombreUnidadMen = computed(() => {
      return [...new Set(
        personas.value
          .filter(p => !nombreUnidadMayor_filtro.value || p.nombreUnidadMayor === nombreUnidadMayor_filtro.value)
          .map(p => p.nombreUnidadMenor)
      )].sort((a, b) => String(a).localeCompare(String(b)));
    });

const nombreAsig = computed(() => {
  return [...new Set(
        personas.value
          .filter(p => (!nombreUnidadMayor_filtro.value || p.nombreUnidadMayor === nombreUnidadMayor_filtro.value) && (!nombreUnidadMenor_filtro.value || p.nombreUnidadMenor === nombreUnidadMenor_filtro.value))
          .map(p => p.NombreAsignatura)
      )].sort((a, b) => String(a).localeCompare(String(b)));
    });


const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    return (
      persona.NombreAsignatura?.toLowerCase().includes(filtros1.value.nombreAsignatura.toLowerCase() || '') &&
      persona.CodigoAsignatura?.toLowerCase().includes(filtros1.value.codigoAsignatura.toLowerCase() || '') &&
      (filtros1.value.estadoProceso ? persona.EstadoProceso?.toString() === filtros1.value.estadoProceso : true) &&
      (filtros1.value.calidad ? persona.Calidad === filtros1.value.calidad : true) &&
      persona.CodigoPAF?.toString().toLowerCase().includes(filtros1.value.codigoPAF.toString().toLowerCase() || '') &&
      persona.Run?.toLowerCase().includes(filtros1.value.run.toLowerCase() || '') &&
      (filtros1.value.jerarquia ? persona.Jerarquia === filtros1.value.jerarquia : true) &&
      persona.FechaUltimaModificacionProceso?.toLowerCase().includes(filtros1.value.fechaUltimaModificacionProceso.toLowerCase() || '') &&
      persona.SemestrePaf?.toLowerCase().includes(filtros1.value.semestre.toLowerCase() || '') &&
      persona.nombreUnidadMenor?.toLowerCase().includes(filtros1.value.nombreUnidadMenor.toLowerCase() || '') &&
      persona.nombreUnidadMayor?.toLowerCase().includes(filtros1.value.nombreUnidadMayor.toLowerCase() || '')
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


const shouldShowTable = ref(true);

const sortData = (newSortBy: string, newSortOrder: string) => {
  sortBy.value = newSortBy;
  sortOrder.value = newSortOrder;
};



//filtros 

const props = defineProps<{
  showButton: boolean; // Define el tipo de la propiedad
}>();

const sortBy = ref('nombres');
const sortOrder = ref('asc');

const route = useRoute();
const isSeguimientoPAF = ref(false);
const isUnidadMayorPAF = ref(false);
const isPaf = ref(false);
const emit = defineEmits<{
  (event: 'filter', filters: any): void;
  (event: 'sort', sortBy: string, order: string): void;
}>();

const estadoProceso = ref('B1');
const applyFilters = () => {
  filtros1.value = {
  codigoPAF: codigoPAF_filtro.value,
  run: run_filtro.value,
  codigoAsignatura: codigoAsignatura_filtro.value,
  semestre: semestres_filtro.value,
  estadoProceso: estado_filtro.value,
  calidad: calidad_filtro.value,
  nombreAsignatura: nombreAsignatura_filtro.value,
  fechaUltimaModificacionProceso: '',
  nombreUnidadMenor: nombreUnidadMenor_filtro.value,
  nombreUnidadMayor: nombreUnidadMayor_filtro.value,
  ruta: '',
  jerarquia: '',
};
  sessionStorage.setItem('codigoPAF', filtros1.value.codigoPAF);
  sessionStorage.setItem('run', filtros1.value.run);
  sessionStorage.setItem('codigoAsignatura', filtros1.value.codigoAsignatura);
  sessionStorage.setItem('semestre', filtros1.value.semestre);
  sessionStorage.setItem('estadoProceso', filtros1.value.estadoProceso);
  sessionStorage.setItem('calidad', filtros1.value.calidad);
  sessionStorage.setItem('nombreAsignatura', filtros1.value.nombreAsignatura);
  sessionStorage.setItem('nombreUnidadMenor', filtros1.value.nombreUnidadMenor);
  sessionStorage.setItem('nombreUnidadMayor', filtros1.value.nombreUnidadMayor);
}


// Desacoplar emisión automática de orden de los filtros
// Resetear filtros
const resetFilters = () => {
  filtros1.value = {
    codigoPAF: '',
    run: '',
    codigoAsignatura: '',
    semestre: '',
    estadoProceso: '',
    calidad: '',
    nombreAsignatura: '',
    fechaUltimaModificacionProceso: '',
    nombreUnidadMenor: '',
    nombreUnidadMayor: '',
    ruta: '',
    jerarquia: '',
  };
  codigoPAF_filtro.value = '';
  run_filtro.value = '';
  codigoAsignatura_filtro.value = '';
  semestres_filtro.value = '';
  estado_filtro.value = '';
  calidad_filtro.value = '';
  nombreAsignatura_filtro.value = '';
  nombreUnidadMenor_filtro.value = '';
  nombreUnidadMayor_filtro.value = '';
  sortBy.value = 'nombres'; // Resetear ordenamiento
  sortOrder.value = 'asc'; // Resetear dirección de ordenamiento
  emit('filter', filtros1.value);
  emit('sort', sortBy.value, sortOrder.value);
  applyFilters();
};

// Cambiar dirección de ordenamiento
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
};
</script>

<style scoped>
.container {
  display: grid;
  grid-template-columns: 25% 75%;

  max-width: 100%;
}

/*filtros */
/* Estilos adicionales para el botón de confirmación */
.confirm-btn {
  background-color: var(--primary-color);
  color: rgb(0, 0, 0);
}

.confirm-btn:hover {
  background-color: #e57000;
}

/* Colores institucionales */
:root {
  --primary-color: #EA7600; /* Color principal USACH */
  --secondary-color: #394049; /* Color secundario USACH */
  --accent-color: #C8102E; /* Complementario */
  --background-color: #f9fafb; /* Fondo neutro */
  --button-background-color: #4CAF50; /* Color de fondo de los botones */
  --button-hover-color: #388E3C; /* Color de hover en los botones */
}

/* Contenedor de filtros */
.filters {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background-color: var(--background-color);
  border-radius: 8px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
  max-width: 100%;
}

/* Estilos de los elementos de filtro */
.filter-item, .sort-item {
  display: flex;
  flex-direction: column;
  position: relative; /* Necesario para aplicar z-index */
  z-index: 10; /* Ajusta el valor según la necesidad */
}

.label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--secondary-color);
}

.input, .select, .btn {
  padding: 0.5rem;
  margin-top: 0.25rem;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  transition: all 0.2s ease;
}

.input:focus, .select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 1px var(--primary-color);
}

.select {
  padding: 0.5rem;
  background-color: #f9fafb;
}

.btn {
  cursor: pointer;
  display: inline-flex;
  justify-content: center;
  align-items: center;
}

/* Botón de reset */
.reset-btn {
  background-color: var(--secondary-color);
  color: #1a0909;
}

.reset-btn:hover {
  background-color: #374151;
}

/* Botón de ordenar */
.sort-btn {
  background-color: #e5e7eb;
  color: #1f2937;
  margin-top: 0.5rem;
}

.sort-btn:hover {
  background-color: #d1d5db;
}
</style>