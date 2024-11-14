<!-- Filtros.vue -->
<template>
  <div class="filters">
    <div class="filter-item">
      <label for="nombre" class="text-sm font-medium text-gray-700">Nombre</label>
      <input
        id="nombre"
        v-model="filtros.nombres"
        type="text"
        class="input"
        placeholder="Filtrar por nombre"
      />
    </div>
    <div class="filter-item">
      <label for="run" class="text-sm font-medium text-gray-700">Run</label>
      <input
        id="run"
        v-model="filtros.run"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>
    <div class="filter-item">
      <label for="codigoPAF" class="text-sm font-medium text-gray-700">Código de PAF</label>
      <input
        id="codigoPAF"
        v-model="filtros.codigoPAF"
        type="text"
        class="input"
        placeholder="Filtrar por código de PAF"
      />
    </div>
    <div class="filter-item">
      <label for="codigoAsignatura" class="text-sm font-medium text-gray-700">Código de Asignatura</label>
      <input
        id="codigoAsignatura"
        v-model="filtros.codigoAsignatura"
        type="text"
        class="input"
        placeholder="Filtrar por código de asignatura"
      />
    </div>
    <div class="filter-item">
      <label for="estadoProceso" class="text-sm font-medium text-gray-700">Estado de Proceso</label>
      <select id="estadoProceso" v-model="filtros.estadoProceso" class="select">
        <option value="">Todos</option>
        <option value="Activo">Activo</option>
        <option value="Inactivo">Inactivo</option>
      </select>
    </div>
    <div class="filter-item">
      <label for="calidad" class="text-sm font-medium text-gray-700">Calidad</label>
      <select id="calidad" v-model="filtros.calidad" class="select">
        <option value="">Todas</option>
        <option value="Contrato Fijo">Contrato Fijo</option>
        <option value="Contrato Temporal">Contrato Temporal</option>
        <option value="Contrato Indefinido">Contrato Indefinido</option>
        <option value="Contrato Parcial">Contrato Parcial</option>
      </select>
    </div>
    <div class="sort-item">
      <label for="sort" class="text-sm font-medium text-gray-700">Ordenar por</label>
      <select v-model="sortBy" class="select">
        <option value="Nombres">Nombre</option>
        <option value="CodigoAsignatura">Código de Asignatura</option>
        <option value="Run">Run</option>
        <option value="FechaInicioContrato">Fecha de Inicio de Contrato</option>
        <option value="FechaUltimaModificacionProceso">Ultima Actualización del Proceso</option>
      </select>
      <button @click="toggleSortOrder" class="btn bg-gray-300">
        Ordenar {{ sortOrder === 'asc' ? 'Ascendente' : 'Descendente' }}
      </button>
    </div>
    <div class="filter-item">
      <button @click="resetFilters" class="btn bg-gray-600 text-white">Resetear</button>
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
}
.filter-item, .sort-item {
  display: flex;
  flex-direction: column;
}
.input, .select, .btn {
  padding: 0.5rem;
  margin-top: 0.25rem;
  border-radius: 4px;
  border: 1px solid #ccc;
}
.btn {
  cursor: pointer;
  width: 100%;
}
.btn:hover {
  opacity: 0.8;
}
.select {
  padding: 0.5rem;
}
</style>
