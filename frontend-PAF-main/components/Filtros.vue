<template>
  <div class="filters">
    <div class="filter-item">
      <label for="nombreAsignatura" class="label">Nombre Asignatura</label>
      <input
        id="nombreAsignatura"
        v-model="filtros.nombreAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por nombre de asignatura"
      />
    </div>
    <div class="filter-item">
      <label for="run" class="label">Run</label>
      <input
        id="run"
        v-model="filtros.run"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>
    <div class="filter-item">
      <label for="codigoPAF" class="label">Código de PAF</label>
      <input
        id="codigoPAF"
        v-model="filtros.codigoPAF"
        type="text"
        class="input"
        placeholder="Filtrar por código de PAF"
      />
    </div>
    <div class="filter-item">
      <label for="codigoAsignatura" class="label">Código de Asignatura</label>
      <input
        id="codigoAsignatura"
        v-model="filtros.codigoAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por código de asignatura"
      />
    </div>
    <div class="filter-item">
      <label for="estadoProceso" class="label">Estado de Proceso</label>
      <select id="estadoProceso" v-model="filtros.estadoProceso" class="select">
        <option value="">Todos</option>
        <option value=1>Estado 1</option>
        <option value=2>Estado 2</option>
        <option value=3>Estado 3</option>
        <option value=4>Estado 4</option>
        <option value=5>Estado 5</option>
        <option value=6>Estado 6</option>
      </select>
    </div>
    <div class="filter-item">
      <label for="calidad" class="label">Calidad</label>
      <select id="calidad" v-model="filtros.calidad" class="select">
        <option value="">Todas</option>
        <option value="Baja">Calidad Baja</option>
        <option value="Media">Calidad Media</option>
        <option value="Alta">Calidad Alta</option>
      </select>
    </div>
    <div class="sort-item">
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
    <div class="filter-item">
      <button @click="resetFilters" class="btn reset-btn">Resetear</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineEmits } from 'vue'

const emit = defineEmits<{
  (event: 'filter', filters: any): void
  (event: 'sort', sortBy: string, order: string): void
}>()

const filtros = ref({
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: '',
  nombreAsignatura: '',
  fechaUltimaModificacionProceso: ''
})

const sortBy = ref('nombres')
const sortOrder = ref('asc')

watch(filtros, (newFilters) => {
  emit('filter', newFilters)
}, { deep: true })

const resetFilters = () => {
  filtros.value = { nombreAsignatura: '', codigoPAF: '', codigoAsignatura: '', run: '', estadoProceso: '', calidad: '', fechaUltimaModificacionProceso: '' }
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
