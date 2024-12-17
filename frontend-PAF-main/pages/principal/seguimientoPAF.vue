<template>
  <div class="container">
    <div>
      <Filtros 
        @filter="filterData" 
        @sort="sortData" 
      />
    </div>
    <div>
      <div v-if="filteredPersonas.length > 0" class="contratos">
        <h2>Contratos Relacionados</h2>
        <table>
          <thead>
            <tr>
              <th>Código PAF</th>
              <th>Run</th>
              <th>Jerarquia</th>
              <th>Nombre de Asignatura</th>
              <th>Estado del Proceso</th>
              <th>Descripción del Proceso</th>
              <th>Fecha de la última Actualización de Estado</th>
              <th>Semestre</th>
              <th>Historial de Estados</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(contrato, index) in paginatedData" :key="index">
              <td>{{ contrato.PipelsoftData.IdPaf }}</td>
              <td>{{ contrato.PipelsoftData.RunEmpleado }}</td>
              <td>{{ contrato.PipelsoftData.Jerarquia }}</td>
              <td>{{ contrato.PipelsoftData.NombreAsignatura }}</td>
              <td>{{ contrato.PipelsoftData.CodEstado }}</td>
              <td>{{ contrato.PipelsoftData.DesEstado }}</td>
              <td>
                {{ new Date(contrato.PipelsoftData.UltimaModificacion).toLocaleDateString() }}
                {{ new Date(contrato.PipelsoftData.UltimaModificacion).toLocaleTimeString() }} 
              </td>
              <td>{{ contrato.PipelsoftData.Semestre }}</td>
              <td>
                <div v-for="(estado, idx) in contrato.PipelsoftData.historialEstados" :key="idx">
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
          <button @click="prevPage" :disabled="currentPage === 1">Anterior</button>
          <span>{{ currentPage }} de {{ totalPages }}</span>
          <button @click="nextPage" :disabled="currentPage === totalPages">Siguiente</button>
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
import Filtros from '../../components/Filtros.vue';
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

const contratos = ref<any[]>([]);
const errorMessage = ref('');
const filtros = ref({
  nombreUnidadMayor: '',
  nombreUnidadMenor: '',
  run: '',
  semestre: '',
  ruta: '/seguimientoPAF',
});
const sortBy = ref('');
const sortOrder = ref('asc');

// Paginación
const currentPage = ref(1);
const itemsPerPage = 20; // Número de elementos por página

const filterData = (newFilters: any) => {
  filtros.value = newFilters;
  currentPage.value = 1; // Resetear la página al filtrar
};

const sortData = (newSortBy: string, newSortOrder: string) => {
  sortBy.value = newSortBy;
  sortOrder.value = newSortOrder;
  currentPage.value = 1; // Resetear la página al ordenar
};

// Computed para filtrar y ordenar
const filteredPersonas = computed(() => {
  let filtered = contratos.value.filter(contrato => {
    return (
      (contrato.PipelsoftData.NombreUnidadMenor || '').toLowerCase().includes((filtros.value.nombreUnidadMenor || '').toLowerCase()) &&
      (contrato.PipelsoftData.NombreUnidadMayor || '').toLowerCase().includes((filtros.value.nombreUnidadMayor || '').toLowerCase()) &&
      (contrato.PipelsoftData.RunEmpleado || '').toLowerCase().includes((filtros.value.run || '').toLowerCase()) &&
      (contrato.PipelsoftData.Semestre || '').toLowerCase().includes((filtros.value.semestre || '').toLowerCase())
    );
  });

  if (sortBy.value) {
    filtered = filtered.sort((a, b) => {
      const compareA = a.PipelsoftData[sortBy.value];
      const compareB = b.PipelsoftData[sortBy.value];
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

const prevPage = () => {
  if (currentPage.value > 1) currentPage.value--;
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) currentPage.value++;
};

// Fetch inicial de datos
const fetchContratos = async () => {
  try {
    const response = await $axios.get(`/api/paf-en-linea/pipelsoft/contratos`);
    if (response.data && Array.isArray(response.data)) {
      contratos.value = response.data;
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
  grid-template-columns: 1fr 3fr;
  gap: 1rem;
  max-width: 100%;
}

/* Formulario */
form {
  margin-bottom: 20px;
}

label {
  font-family: "Helvetica Neue LT", sans-serif;
  color: #394049;
}

input {
  padding: 8px;
  margin-right: 10px;
  border: 1px solid #394049;
  border-radius: 4px;
}

button {
  padding: 8px 16px;
  background-color: #EA7600; /* Color institucional */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-family: "Bebas Neue Pro", sans-serif;
}

button:hover {
  background-color: #C8102E; /* Variante complementaria */
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

/* Paginación */
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination button {
  padding: 8px 16px;
  background-color: #EA7600;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  background-color: #f0f0f0;
  cursor: not-allowed;
}

.pagination span {
  font-size: 1rem;
}
</style>
