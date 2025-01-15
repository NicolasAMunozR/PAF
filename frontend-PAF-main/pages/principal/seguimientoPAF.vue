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
  
    <!-- Botón de reset -->
    <div class="filter-item">
      <button @click="resetFilters" class="btn reset-btn">Resetear</button>
    </div>

    <!-- Botón de confirmar filtros -->
    <div class="filter-item">
      <button @click="applyFilters" class="btn confirm-btn">Confirmar Filtros</button>
    </div>
  </div>
    <div>
      <div v-if="filteredPersonas.length > 0" class="contratos">
        <h2>Contratos Relacionados</h2>
        <table>
      <thead>
        <tr>
          <th @click="sortData('id_paf')">
            Código PAF
            <span>
              {{ sortBy === 'id_paf' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('run_empleado')">
            Run
            <span>
              {{ sortBy === 'run_empleado' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('jerarquia')">
            Jerarquia
            <span>
              {{ sortBy === 'jerarquia' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('nombre_asignatura_list')">
            Nombre de Asignatura
            <span>
              {{ sortBy === 'nombre_asignatura_list' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('cod_estado')">
            Estado del Proceso
            <span>
              {{ sortBy === 'cod_estado' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('des_estado')">
            Descripción del Proceso
            <span>
              {{ sortBy === 'des_estado' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('ultima_modificacion')">
            Fecha de la última Actualización de Estado
            <span>
              {{ sortBy === 'ultima_modificacion' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th @click="sortData('semestre')">
            Semestre
            <span>
              {{ sortBy === 'semestre' ? (sortOrder === 'asc' ? '▲' : '▼') : '-' }}
            </span>
          </th>
          <th>Historial de Estados</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(contrato, index) in paginatedData" :key="index">
          <td>{{ contrato.id_paf }}</td>
          <td>{{ contrato.run_empleado }}</td>
          <td>{{ contrato.jerarquia }}</td>
          <td>{{ Array.isArray(contrato.nombre_asignatura_list) ? contrato.nombre_asignatura_list.join(', ') : contrato.nombre_asignatura_list }}</td>
          <td>{{ contrato.cod_estado }}</td>
          <td>{{ contrato.des_estado }}</td>
          <td>{{ new Date(contrato.ultima_modificacion).toLocaleDateString() }}</td>
          <td>{{ contrato.semestre }}</td>
          <td>
            <div v-for="(estado, idx) in contrato.historial_estados" :key="idx">
              <p>
                <strong>{{ estadoProceso(estado.estado) }}</strong>: 
                {{ calcularTiempoEnEstado(estado.fechaInicio) }} días
              </p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

         <!-- Paginación -->
         <div class="pagination">
          <button 
            @click="goToPage(1)" 
            :disabled="currentPage === 1">
            «
          </button>
          <button 
            v-for="page in pageNumbers" 
            :key="page" 
            :class="{ active: currentPage === page }"
            @click="goToPage(Number(page))">
            <span v-if="page === '...'">...</span>
            <span v-else>{{ page }}</span>
          </button>
          <button 
            @click="goToPage(totalPages)" 
            :disabled="currentPage === totalPages">
            »
          </button>
        </div>
      </div>

      <div v-else-if="errorMessage" class="error">
        <p>{{ errorMessage }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

const semestreMasActual = ref('');
const semestres_filtro = ref('');
const run_filtro = ref('');
const nombreUnidadMenor_filtro = ref('');
const nombreUnidadMayor_filtro = ref('');
const estado_filtro = ref('');
const contratos = ref<any[]>([]);
const errorMessage = ref('');
const filtros = ref({
  nombreUnidadMayor: '',
  nombreUnidadMenor: '',
  run: '',
  semestre: '',
  estadoProceso: '',
  ruta: '/seguimientoPAF',
});
const sortBy = ref('');
const sortOrder = ref('asc');

// Paginación
const currentPage = ref(1);
const itemsPerPage = 10; // Número de elementos por página
const isFirstLoad = ref(true);

const filterData = (newFilters: any) => {
  
    filtros.value = { ...newFilters }; // Si no es la primera carga, no modificamos el semestre
  
  currentPage.value = 1; // Resetear la página al filtrar
};

const sortData = (key: string) => {
  if (sortBy.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortBy.value = key;
    sortOrder.value = 'asc';
  }
  currentPage.value = 1; // Resetear a la primera página
};

// Computed para filtrar y ordenar
const filteredPersonas = computed(() => {
  let filtered = contratos.value.filter(contrato => {
    return (
      (contrato.nombre_unidad_menor || '').toLowerCase().includes((filtros1.value.nombreUnidadMenor || '').toLowerCase()) &&
      (contrato.nombre_unidad_mayor || '').toLowerCase().includes((filtros1.value.nombreUnidadMayor || '').toLowerCase()) &&
      (contrato.run_empleado || '').toLowerCase().includes((filtros1.value.run || '').toLowerCase()) &&
      (contrato.semestre || '').toLowerCase().includes((filtros1.value.semestre || '').toLowerCase()) &&
      (filtros1.value.estadoProceso ? contrato.cod_estado?.toString() === filtros1.value.estadoProceso : true)
    );
  });

  if (sortBy.value) {
    filtered = filtered.sort((a, b) => {
      const compareA = a[sortBy.value];
      const compareB = b[sortBy.value];
      if (compareA < compareB) return sortOrder.value === 'asc' ? -1 : 1;
      if (compareA > compareB) return sortOrder.value === 'asc' ? 1 : -1;
      return 0;
    });
  }

  return filtered;
});


// Computed para la paginación
const totalPages = computed(() => {
  return Math.ceil(filteredPersonas.value.length / itemsPerPage);
});

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return filteredPersonas.value.slice(start, end);
});

const pageNumbers = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  const range = 2; // Cantidad de páginas a mostrar alrededor de la actual
  const pages = [];

  if (total <= range * 2 + 5) {
    for (let i = 1; i <= total; i++) {
      pages.push(i);
    }
  } else {
    let start = Math.max(1, current - range);
    let end = Math.min(total, current + range);

    if (current <= range + 1) {
      end = Math.min(total, range * 2 + 3);
    } else if (current > total - range - 1) {
      start = Math.max(1, total - range * 2 - 2);
    }

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (start > 1) {
      pages.unshift(1, '...');
    }

    if (end < total) {
      pages.push('...', total);
    }
  }

  return pages;
});

const goToPage = (page: number) => {
  if (typeof page === 'number') {
    currentPage.value = page;
  }
};

// Fetch inicial de datos
const fetchContratos = async () => {
  try {
    const response = await $axios.get(`/api/paf-en-linea/pipelsoft/FiltroSemestre/`);
    if (response.data && Array.isArray(response.data)) {
      contratos.value = response.data;
      // Supongamos que response.data es el arreglo que has mencionado
semestreMasActual.value = response.data
  .map(item => item.semestre) // Extrae el semestre de cada objeto
  .sort((a, b) => {
    // Convierte los semestres en fechas para poder compararlos
    const [mesA, anioA] = a.split('-');
    const [mesB, anioB] = b.split('-');

    // Compara año y mes
    if (parseInt(anioA) === parseInt(anioB)) {
      return parseInt(mesB) - parseInt(mesA); // Ordena por mes descendente
    }
    return parseInt(anioB) - parseInt(anioA); // Ordena por año descendente
  })[0]; // Devuelve el semestre más reciente
  if (sessionStorage.getItem('run')) {
    run_filtro.value = sessionStorage.getItem('run') || '';
  }
  if (sessionStorage.getItem('semestre') || sessionStorage.getItem('semestre') == '') {
    semestres_filtro.value = sessionStorage.getItem('semestre') || '';
  }
  else{
    semestres_filtro.value = semestreMasActual.value;
  }
  if (sessionStorage.getItem('estadoProceso')) {
    estado_filtro.value = sessionStorage.getItem('estadoProceso') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMenor')) {
    nombreUnidadMenor_filtro.value = sessionStorage.getItem('nombreUnidadMenor') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMayor')) {
    nombreUnidadMayor_filtro.value = sessionStorage.getItem('nombreUnidadMayor') || '';
  }
  applyFilters();
      
    } else {
      errorMessage.value = 'No se encontraron contratos.';
    }
  } catch (error) {
    errorMessage.value = 'Hubo un error al obtener los datos.';
    console.error(error);
  }
};

onMounted(() => {
  fetchContratos();
});

const filtros1 = ref({
  run: '',
  semestre: '',
  estadoProceso: '',
  nombreUnidadMenor: '',
  nombreUnidadMayor: '',
});

const semestres = computed(() => {
  // Ordenar los semestres y eliminar duplicados
  const semestresOrdenados = [...new Set(contratos.value.map(p => p.semestre))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return semestresOrdenados; // O retorna [ultimoSemestre] si solo necesitas el último
});

const nombreUnidadMay = computed(() => {
  const nombreUnidadMayor = [...new Set(contratos.value.map(p => p.nombre_unidad_mayor))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return nombreUnidadMayor; // O retorna [ultimoSemestre] si
});

const nombreUnidadMen = computed(() => {
      return [...new Set(
        contratos.value
          .filter(p => !nombreUnidadMayor_filtro.value || p.nombre_unidad_mayor === nombreUnidadMayor_filtro.value)
          .map(p => p.nombre_unidad_menor)
      )].sort((a, b) => String(a).localeCompare(String(b)));
    });

    const applyFilters = () => {
  filtros1.value = {
  run: run_filtro.value,
  semestre: semestres_filtro.value,
  estadoProceso: estado_filtro.value,
  nombreUnidadMenor: nombreUnidadMenor_filtro.value,
  nombreUnidadMayor: nombreUnidadMayor_filtro.value,

};
  sessionStorage.setItem('run', filtros1.value.run);
  sessionStorage.setItem('semestre', filtros1.value.semestre);
  sessionStorage.setItem('estadoProceso', filtros1.value.estadoProceso);
  sessionStorage.setItem('nombreUnidadMenor', filtros1.value.nombreUnidadMenor);
  sessionStorage.setItem('nombreUnidadMayor', filtros1.value.nombreUnidadMayor);
}


// Desacoplar emisión automática de orden de los filtros
// Emit function
const emit = defineEmits(['filter', 'sort']);

// Resetear filtros
const resetFilters = () => {
  filtros1.value = {
    run: '',
    semestre: '',
    estadoProceso: '',
    nombreUnidadMenor: '',
    nombreUnidadMayor: '',
  };
  run_filtro.value = '';
  semestres_filtro.value = '';
  estado_filtro.value = '';
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

// Función para mapear estados
const estadoProceso = (estado: string): string => {
  switch (estado) {
    case "1": return "Estado 1";
    case "2": return "Estado 2";
    case "3": return "Estado 3";
    case "4": return "Estado 4";
    case "5": return "Estado 5";
    case "6": return "Estado 6";
    default: return "Desconocido";
  }
};

// Calcular tiempo en estado
const calcularTiempoEnEstado = (fechaInicio: string): number => {
  const fechaActual = new Date();
  const fechaInicioEstado = new Date(fechaInicio);
  const diferenciaTiempo = fechaActual.getTime() - fechaInicioEstado.getTime();
  return Math.floor(diferenciaTiempo / (1000 * 3600 * 24));
};
</script>

<style scoped>
 /* Contenedor principal */
 .container {
  display: grid;
  grid-template-columns: 25% 75%;

  max-width: 100%;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 5px;
}

.pagination button {
  padding: 8px 12px;
  background-color: #f0f0f0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button.active {
  background-color: #EA7600;
  color: white;
}

.pagination button:disabled {
  background-color: #EA7600;
  cursor: not-allowed;
}


button {
  padding: 8px 16px;
  background-color: #EA7600; /* Color institucional */
  color: rgb(3, 0, 0);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-family: "Bebas Neue Pro", sans-serif;
}

/* Tabla */
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th, td {
  padding: 10px;
  text-align: left;
  border: 1px solid #ddd;
}

th {
  background-color: #394049;
  color: white;
}

.error {
  color: #C8102E; /* Color institucional para errores */
  font-weight: bold;
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
