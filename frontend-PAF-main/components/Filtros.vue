<template>
  <div class="filters">
    <!-- Filtro Nombre Asignatura -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
      <label for="nombreAsignatura" class="label">Nombre Asignatura</label>
      <input
        v-model="filtros.nombreAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de asignatura"
      />
    </div>
    
    <!-- Filtro Run -->
    <div class="filter-item">
      <label for="run" class="label">Run</label>
      <input
        v-model="filtros.run"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>
    
    <!-- Filtro Código de PAF -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
      <label for="codigoPAF" class="label">Código de PAF</label>
      <input
        v-model="filtros.codigoPAF"
        type="text"
        class="input" 
        placeholder="Filtrar por código de PAF"
      />
    </div>
    
    <!-- Filtro Código de Asignatura -->
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
      <label for="codigoAsignatura" class="label">Código de Asignatura</label>
      <input
        v-model="filtros.codigoAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por código de asignatura"
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
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
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
    <div class="filter-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
      <label for="calidad" class="label">Categoria</label>
      <select id="calidad" v-model="filtros.calidad" class="select">
        <option value="">Todas</option>
        <option value="PROFESOR HORAS CLASES">Profesor por hora</option>
      </select>
    </div>

    <!-- Filtro Nombre Unidad Mayor (solo visible si es seguimientoPAF) -->
    <div class="filter-item" v-if="isSeguimientoPAF">
      <label for="nombreUnidadMayor" class="label">Nombre Unidad Mayor</label>
      <input
        v-model="filtros.nombreUnidadMayor"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de unidad mayor"
      />
    </div>
    
    <!-- Filtro Nombre Unidad Menor (solo visible si es seguimientoPAF o unidadMayorPAF) -->
    <div class="filter-item" v-if="isSeguimientoPAF || isUnidadMayorPAF">
      <label for="nombreUnidadMenor" class="label">Nombre Unidad Menor</label>
      <input
        v-model="filtros.nombreUnidadMenor"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de unidad menor"
      />
    </div>
    
    <!-- Ordenar por (solo visible si es seguimientoPAF) -->
    <div class="sort-item" v-if="!isSeguimientoPAF && !isUnidadMayorPAF">
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
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineEmits, onMounted } from 'vue'
import UnidadMayorPAF from '../pages/unidadMayorPAF/index.vue'
import SeguimientoPAF from '../pages/seguimientoPAF/index.vue'
import { useRoute } from 'vue-router'

const emit = defineEmits<{
  (event: 'filter', filters: any): void
  (event: 'sort', sortBy: string, order: string): void
}>()

const filtros = ref({
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
})

const sortBy = ref('nombres')
const sortOrder = ref('asc')

const routes = [
  { path: '/seguimientoPAF', name: 'seguimientoPAF', component: SeguimientoPAF },
  { path: '/unidadMayorPAF', name: 'unidadMayorPAF', component: UnidadMayorPAF },
  // otras rutas...
];
const route = useRoute();
const isSeguimientoPAF = ref(false);
const isUnidadMayorPAF = ref(false);

watchEffect(() => {
  const currentRoute = useRoute().name;
  isSeguimientoPAF.value = currentRoute === 'seguimientoPAF';
  isUnidadMayorPAF.value = currentRoute === 'unidadMayorPAF';
});


watch(filtros, (newFilters) => {
  let filtersToEmit: Partial<typeof filtros.value> = { ...newFilters };

  // Filtra campos según las rutas
  if (isSeguimientoPAF.value) {
    // Si es seguimientoPAF, emite todos los filtros
    emit('filter', filtersToEmit);
  } else if (isUnidadMayorPAF.value) {
    // Si es UnidadMayorPAF, excluye 'nombreUnidadMenor'
    filtersToEmit = Object.fromEntries(
      Object.entries(newFilters).filter(([key]) => key == 'nombreUnidadMenor' || key == 'run' || key == 'semestre')
    );
    emit('filter', filtersToEmit);
  } else {
    // En cualquier otro caso, emite los filtros sin cambios
    emit('filter', filtersToEmit);
  }
}, { deep: true });

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
  }
  emit('filter', filtros.value)
}

const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  emit('sort', sortBy.value, sortOrder.value)
}

watch([sortBy, sortOrder], ([newSortBy, newSortOrder]) => {
  emit('sort', newSortBy, newSortOrder)
})
</script>


<style scoped>

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
