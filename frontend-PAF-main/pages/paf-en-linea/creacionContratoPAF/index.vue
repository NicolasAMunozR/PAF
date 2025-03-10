<template>
    <div class="container">
  <div class="filters">
        <div class="filter-item">
      <label for="run" class="label">Run</label>
      <input
        v-model="run_filtro"
        type="text"
        class="input"
        placeholder="Filtrar por Run"
      />
    </div>

    <div class="filter-item">
      <label for="unidadMayor" class="label">Unidad Mayor</label>
      <input
        v-model="unidadMayor_filtro"
        type="text"
        class="input"
        placeholder="Filtrar por unidadMayor"
      />
    </div>
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
    <div class="filter-item">
      <button @click="applyFilters" class="btn confirm-btn">Confirmar Filtros</button>
    </div>
    <br />
    <div class="filter-item">
    <label for="semestre" class="label">Semestre para generar contratos</label>
    <select id="semestre" v-model="semestre_contrato" class="select">
      <option
        v-for="sem in semestres"
        :key="sem"
        :value="sem"
      >
        {{ sem }}
      </option>
    </select>
  </div>
    <div class="filter-item">
      <button @click="applyContrats" class="btn confirm-btn">Generar contratos</button>
    </div>
  </div>
  <Tabla :data="filteredPersonas" :showButton="false" :showButtons="false" :show="true" :contratos="false"/>
</div>
</template>

<script setup lang="ts">
import Tabla from '../../../components/Tabla.vue'
import { ref, computed, onMounted } from 'vue'

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const personas = ref<Persona[]>([]);
const unidadMayor_filtro = ref('');
const semestres_filtro = ref('');
const run_filtro = ref('');
const semestre_contrato = ref('');
const semestreMasActual = ref('');

const sortBy = ref('nombres');
const sortOrder = ref('asc');

const semestres = computed(() => {
  // Ordenar los semestres y eliminar duplicados
  const semestresOrdenados = [...new Set(personas.value.map(p => p.Semestre))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return semestresOrdenados; // O retorna [ultimoSemestre] si solo necesitas el último
});

const emit = defineEmits<{
  (event: 'filter', filters: any): void;
  (event: 'sort', sortBy: string, order: string): void;
}>();

interface Persona {
  Run: string;
  Semestre: string;
  NombreAsignatura: string;
  seccion:  string;
  CantidadHoras:  string;
  CodigoAsignatura: string;
  NombreUnidadMayor: string;
}

const filtros1 = ref({
  run: '',
  codigoAsignatura: '',
  semestre: '',
  nombreAsignatura: '',
  cantidadHoras: '',
  seccion: '',
  nombreUnidadMayor: '',
});

const applyContrats = async () => {
  const response = await $axios.get(`/api/paf-en-linea/archivo/archivo/generar/Automatico/${semestre_contrato}`);
  console.log(response.data);
}

const applyFilters = async () => {
  filtros1.value = {
    run: run_filtro.value,
    codigoAsignatura: '',
    semestre: semestres_filtro.value,
    nombreAsignatura: '',
    cantidadHoras: '',
    seccion: '',
    nombreUnidadMayor: unidadMayor_filtro.value,
  }
  /*
  if(filtros1.value.nombreUnidadMayor !== ''){
    const response = await $axios.get(`/api/paf-en-linea/archivo/profesores/sin-contrato/${filtros1.value.nombreUnidadMayor}`);
    console.log(response.data);
    // Asegurar que `profesores_no_comunes` es un array antes de mapearlo
    if (Array.isArray(response.data.profesores_no_comunes)) {
      personas.value = response.data.profesores_no_comunes.map((item: any) => ({
        Run: item.run,
        Semestre: item.semestre,
        NombreAsignatura: item.nombre_asignatura,
        seccion: item.seccion,
        CantidadHoras: String(item.CantidadHoras), // Convertir a string si es necesario
        CodigoAsignatura: item.codigo_asignatura,
      }));
    } else {
      console.error("El campo 'profesores_no_comunes' no es un array");
    }
  } else {
    const response = await $axios.get('/api/paf-en-linea/archivo/profesores/sin-contrato');

    // Asegurar que `profesores_no_comunes` es un array antes de mapearlo
    if (Array.isArray(response.data.profesores_no_comunes)) {
      personas.value = response.data.profesores_no_comunes.map((item: any) => ({
        Run: item.run,
        Semestre: item.semestre,
        NombreAsignatura: item.nombre_asignatura,
        seccion: item.seccion,
        CantidadHoras: String(item.CantidadHoras), // Convertir a string si es necesario
        CodigoAsignatura: item.codigo_asignatura,
        NombreUnidadMayor: item.NombreUnidadMayor ?? 'N/A',
      }));
    } else {
      console.error("El campo 'profesores_no_comunes' no es un array");
    }
  }
  */
}

const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
};

const resetFilters = () => {
  filtros1.value = {
    run: '',
    codigoAsignatura: '',
    semestre: '',
    nombreAsignatura: '',
    cantidadHoras: '',
    seccion: '',
    nombreUnidadMayor: '',
  }
  unidadMayor_filtro.value = '';
  semestres_filtro.value = '';
  run_filtro.value = '';
  sortBy.value = 'nombres'; // Resetear ordenamiento
  sortOrder.value = 'asc'; // Resetear dirección de ordenamiento
  emit('filter', filtros1.value);
  emit('sort', sortBy.value, sortOrder.value);
  applyFilters();
};

const filteredPersonas = computed(() => {
  let filtered = personas.value.filter(persona => {
    return (
      persona.NombreAsignatura?.toLowerCase().includes(filtros1.value.nombreAsignatura.toLowerCase() || '') &&
      persona.CodigoAsignatura?.toLowerCase().includes(filtros1.value.codigoAsignatura.toLowerCase() || '') &&
      persona.Run?.toLowerCase().includes(filtros1.value.run.toLowerCase() || '') &&
      persona.seccion?.toLowerCase().includes(filtros1.value.seccion.toLowerCase() || '') &&
      persona.Semestre?.toLowerCase().includes(filtros1.value.semestre.toLowerCase() || '') &&
      persona.CantidadHoras?.toLowerCase().includes(filtros1.value.cantidadHoras.toLowerCase() || '') &&
      persona.NombreUnidadMayor?.toLowerCase().includes(filtros1.value.nombreUnidadMayor.toLowerCase() || '')
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
    const response = await $axios.get('/api/paf-en-linea/archivo/profesores/sin-contrato');
    if (Array.isArray(response.data.profesores_no_comunes)) {
      personas.value = response.data.profesores_no_comunes.map((item: any) => ({
        Run: item.run,
        Semestre: item.semestre,
        NombreAsignatura: item.nombre_asignatura,
        seccion: item.seccion,
        CantidadHoras: String(item.CantidadHoras), // Convertir a string si es necesario
        CodigoAsignatura: item.codigo_asignatura,
        NombreUnidadMayor: item.NombreUnidadMayor ?? 'N/A',
      }));
      const semestresOrdenados = response.data.profesores_no_comunes
        .map((item: any) => item.semestre)
        .sort((a: string, b: string) => {
          const [mesA, anioA] = a.split('-');
          const [mesB, anioB] = b.split('-');
          if (parseInt(anioA) === parseInt(anioB)) {
            return parseInt(mesB) - parseInt(mesA); // Ordena por mes descendente
          }
          return parseInt(anioB) - parseInt(anioA); // Ordena por año descendente
        });

      // Asignar el semestre más reciente
      semestreMasActual.value = semestresOrdenados[0];
      semestres_filtro.value = semestreMasActual.value;
      semestre_contrato.value = semestreMasActual.value;
      applyFilters();
    } else {
      console.error("El campo 'profesores_no_comunes' no es un array");
    }
  } catch (error) {
    console.error("Error al obtener los datos:", error);
  }
});
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