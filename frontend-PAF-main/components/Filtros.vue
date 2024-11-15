<template>
  <div class="filters">
    <div class="filter-item">
      <label for="nombre" class="label">Nombre</label>
      <input
        id="nombre"
        v-model="filtros.nombres"
        type="text"
        class="input"
        placeholder="Filtrar por nombre"
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
        <option value="Activo">Activo</option>
        <option value="Inactivo">Inactivo</option>
      </select>
    </div>
    <div class="filter-item">
      <label for="calidad" class="label">Calidad</label>
      <select id="calidad" v-model="filtros.calidad" class="select">
        <option value="">Todas</option>
        <option value="Contrato Fijo">Contrato Fijo</option>
        <option value="Contrato Temporal">Contrato Temporal</option>
        <option value="Contrato Indefinido">Contrato Indefinido</option>
        <option value="Contrato Parcial">Contrato Parcial</option>
      </select>
    </div>
    <div class="sort-item">
      <label for="sort" class="label">Ordenar por</label>
      <select v-model="sortBy" class="select">
        <option value="Nombres">Nombre</option>
        <option value="CodigoAsignatura">Código de Asignatura</option>
        <option value="Run">Run</option>
        <option value="FechaInicioContrato">Fecha de Inicio de Contrato</option>
        <option value="FechaUltimaModificacionProceso">Ultima Actualización del Proceso</option>
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
  nombres: '',
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: ''
})

const sortBy = ref('nombres')
const sortOrder = ref('asc')

watch(filtros, (newFilters) => {
  emit('filter', newFilters)
}, { deep: true })

const resetFilters = () => {
  filtros.value = { nombres: '', codigoPAF: '', codigoAsignatura: '', run: '', estadoProceso: '', calidad: '' }
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
.filters {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background-color: #f9fafb;
  border-radius: 8px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

.filter-item, .sort-item {
  display: flex;
  flex-direction: column;
}

.label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4b5563;
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
  border-color: #4f46e5;
  box-shadow: 0 0 0 1px #4f46e5;
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

.reset-btn {
  background-color: #4b5563;
  color: #ffffff;
}

.reset-btn:hover {
  background-color: #374151;
}

.sort-btn {
  background-color: #e5e7eb;
  color: #1f2937;
  margin-top: 0.5rem;
}

.sort-btn:hover {
  background-color: #d1d5db;
}
</style>
