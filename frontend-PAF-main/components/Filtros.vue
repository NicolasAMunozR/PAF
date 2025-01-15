<template>
  <div class="filters">
        <!-- Filtro Run -->
        <div class="filter-item" v-if="!isPaf">
      <label for="run" class="label">Run</label>
      <input
        v-model="filtros.run"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>

    <!-- Filtro Código de PAF -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF && !isPaf">
      <label for="codigoPAF" class="label">Código de PAF</label>
      <input
        v-model="filtros.codigoPAF"
        type="text"
        class="input" 
        placeholder="Filtrar por código de PAF"
      />
    </div>

    <!-- Filtro Nombre Asignatura -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF && !isPaf">
      <label for="nombreAsignatura" class="label">Nombre Asignatura</label>
      <input
        v-model="filtros.nombreAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de asignatura"
      />
    </div>
    
    <!-- Filtro Código de Asignatura -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF && !isPaf">
      <label for="codigoAsignatura" class="label">Código de Asignatura</label>
      <input
        v-model="filtros.codigoAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por código de asignatura"
      />
    </div>

    <!-- Filtro Nombre Unidad Mayor (solo visible si es seguimientoPAF) -->
    <div class="filter-item" v-if="!isUnidadMayorPAF && !isPaf">
      <label for="nombreUnidadMayor" class="label">Nombre Unidad Mayor</label>
      <input
        v-model="filtros.nombreUnidadMayor"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de unidad mayor"
      />
    </div>
    
    <!-- Filtro Nombre Unidad Menor (solo visible si es seguimientoPAF o unidadMayorPAF) -->
    <div class="filter-item" v-if="showButton && !isPaf">
      <label for="nombreUnidadMenor" class="label">Nombre Unidad Menor</label>
      <input
        v-model="filtros.nombreUnidadMenor"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de unidad menor"
      />
    </div>

    <!-- Filtro semestre -->
    <div class="filter-item">
      <label for="semestre" class="label">Semestre</label>
      <input
        v-model="filtros.semestre"
        type="text"
        class="input"
        placeholder="Filtrar por código Semestre"
      />
    </div>
    
    <!-- Filtro Estado de Proceso -->
    <div class="filter-item" v-if="!isPaf">
      <label for="estadoProceso" class="label">Estado de Proceso</label>
      <select id="estadoProceso" v-model="filtros.estadoProceso" class="select">
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
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF && !isPaf">
      <label for="calidad" class="label">Categoria</label>
      <select id="calidad" v-model="filtros.calidad" class="select">
        <option value="">Todas</option>
        <option value="PROFESOR HORAS CLASES">Profesor por hora</option>
      </select>
    </div>
    
    <!-- Ordenar por (solo visible si es seguimientoPAF) -->
    <div class="sort-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF && !isPaf">
      <label for="sort" class="label">Ordenar por</label>
      <select v-model="sortBy" class="select">
        <option value="NombreAsignatura">Nombre de Asignatura</option>
        <option value="CodigoAsignatura">Código de Asignatura</option>
        <option value="Run">Run</option>
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
</template>

<script setup lang="ts">
import { ref, defineEmits, onMounted, watchEffect } from 'vue';
import { useRoute } from 'vue-router';
// Define las props que el componente espera recibir
const props = defineProps<{
  showButton: boolean; // Define el tipo de la propiedad
}>();
let semestres = ref('');
let codigoPAF_filtro = ref('');
let run_filtro = ref('');
let codigoAsignatura_filtro = ref('');
let calidad_filtro = ref('');
let nombreAsignatura_filtro = ref('');
let nombreUnidadMenor_filtro = ref('');
let nombreUnidadMayor_filtro = ref('');

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
const filtros = ref({
  codigoPAF: codigoPAF_filtro && route.path === "/principal/personas" ? codigoPAF_filtro : '',
  run: run_filtro && route.path === "/principal/personas" ? run_filtro : '',
  codigoAsignatura: codigoAsignatura_filtro && route.path === "/principal/personas" ? codigoAsignatura_filtro : '',
  semestre: semestres ? semestres : '',
  estadoProceso: estadoProceso && route.path === "/principal/personas" ? estadoProceso : '',
  calidad: calidad_filtro && route.path === "/principal/personas" ? calidad_filtro : '',
  nombreAsignatura: nombreAsignatura_filtro && route.path === "/principal/personas" ? nombreAsignatura_filtro : '',
  fechaUltimaModificacionProceso: '',
  nombreUnidadMenor: nombreUnidadMenor_filtro && route.path === "/principal/personas" ? nombreUnidadMenor_filtro : '',
  nombreUnidadMayor: nombreUnidadMayor_filtro && route.path === "/principal/personas" ? nombreUnidadMayor_filtro : '',
  ruta: '',
});


const aplicar = () => {
  if (route.path === "/principal/seguimientoPAF") {
    isSeguimientoPAF.value = true;
  } else if (route.path === "/principal/unidadMayorPAF") {
    isUnidadMayorPAF.value = true;
  } else if (route.path === "/principal/personas/paf") {
    isPaf.value = true;
  }
}

// Resetear filtros al montar el componente
onMounted(() => {
  if (props.showButton) {
    // Realiza alguna acción si el botón está visible
  }

  const codigoPAF_Filtro = localStorage.getItem('codigoPAF_filtro');
    const run_Filtro = localStorage.getItem('run_filtro');
    const codigoAsignatura_Filtro = localStorage.getItem('codigoAsignatura_filtro');
    const estadoProceso_Filtro = localStorage.getItem('estadoProceso_filtro');
    const calidad_Filtro = localStorage.getItem('calidad_filtro');
    const nombreAsignatura_Filtro = localStorage.getItem('nombreAsignatura_filtro');
    const nombreUnidadMenor_Filtro = localStorage.getItem('nombreUnidadMenor_filtro');
    const nombreUnidadMayor_Filtro = localStorage.getItem('nombreUnidadMayor_filtro');

  const semestre = localStorage.getItem('semestre');
  
  if (semestre) {
    if(semestre === "todos"){
      semestres.value = "";
    }
    else{
    semestres.value = semestre;
    }
  }
  if (codigoPAF_Filtro) {
    codigoPAF_filtro.value = codigoPAF_Filtro;
  }
  if (run_Filtro) {
    run_filtro.value = run_Filtro;
  }
  if (codigoAsignatura_Filtro) {
    codigoAsignatura_filtro.value = codigoAsignatura_Filtro;
  }
  if (estadoProceso_Filtro) {
    estadoProceso.value = estadoProceso_Filtro;
  }
  if (calidad_Filtro) {
    calidad_filtro.value = calidad_Filtro;
  }
  if (nombreAsignatura_Filtro) {
    nombreAsignatura_filtro.value = nombreAsignatura_Filtro;
  }
  if (nombreUnidadMenor_Filtro) {
    nombreUnidadMenor_filtro.value = nombreUnidadMenor_Filtro;
  }
  if (nombreUnidadMayor_Filtro) {
    nombreUnidadMayor_filtro.value = nombreUnidadMayor_Filtro;
  }
  setTimeout(() => {
    if(!semestre) {
    window.location.reload(); // Recarga la página completa
  } 
    applyFilters();
  }, 100);
  aplicar();
  applyFilters();
});

// Desacoplar emisión automática de orden de los filtros
const applyFilters = () => {
  let filtersToEmit: Partial<typeof filtros.value> = { ...filtros.value };

  // Filtra campos según las rutas
  if (isSeguimientoPAF.value) {
    filtersToEmit = Object.fromEntries(
      Object.entries(filtros.value).filter(([key]) => key === 'nombreUnidadMenor' || key === 'nombreUnidadMayor' || key === 'run' || key === 'semestre' || key === 'estadoProceso')
    );
    emit('filter', filtersToEmit);
  } else if (isUnidadMayorPAF.value) {
    filtersToEmit = Object.fromEntries(
      Object.entries(filtros.value).filter(([key]) => key === 'nombreUnidadMenor' || key === 'run' || key === 'semestre' || key === 'estadoProceso')
    );
    emit('filter', filtersToEmit);
  } else if (isPaf.value) {
    filtersToEmit = Object.fromEntries(
      Object.entries(filtros.value).filter(([key]) => key === 'semestre')
    );
    emit('filter', filtersToEmit);
  } else {
    emit('filter', filtersToEmit);
  }

  // Emitir ordenamiento solo cuando se confirma
  emit('sort', sortBy.value, sortOrder.value);
};

// Resetear filtros
const resetFilters = () => {
  filtros.value = {
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
  };
  sortBy.value = 'nombres'; // Resetear ordenamiento
  sortOrder.value = 'asc'; // Resetear dirección de ordenamiento
  emit('filter', filtros.value);
  emit('sort', sortBy.value, sortOrder.value);
  applyFilters();
};

// Cambiar dirección de ordenamiento
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
};
</script>

<style scoped>
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
}

/* Estilos de los elementos de filtro */
.filter-item, .sort-item {
  display: flex;
  flex-direction: column;
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
