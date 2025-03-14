<template>
  <div class="container">
    <!-- Advertencias -->
    <div v-if="warnings.length" class="warnings">
      <h3>Advertencias:</h3>
      <ul>
        <li v-for="warning in warnings" :key="warning">{{ warning }}</li>
      </ul>
    </div>
    
    <!-- Filtros -->
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
    <br />
    <div>
      <button @click="Excel" class="btn confirm-btn">Importal Excel</button>
    </div>
  </div>
    
    <!-- Tabla -->
    <Tabla :data="filteredPersonas" :showButtons="false" @rowStatusChanged="handleRowStatusChanged" />
  </div>
</template>

<script setup lang="ts">
import Tabla from '../../components/Tabla.vue';
import { ref, computed, onMounted, defineEmits } from 'vue';

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

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

interface Persona {
  ID: number;
  IdPaf: number;
  CodigoAsignatura: string;
  Run: string;
  EstadoProceso: string;
  Calidad: string;
  Jerarquia: string;
  FechaInicioContrato: string;
  FechaFinContrato: string;
  NombreAsignatura: string;
  cupo: number;
  codigo_modificacion: number;
  bandera_modificacion: number;
  descripcion_modificacion: string;
  seccion: string;
  rowClass?: string;
  SemestrePaf: string;
  DesEstado: string;
  UltimaModificacion: string;
  Comentario: string;
  codigos_extras: string[];
  comentarios_extras: string[];
}

const personas = ref<Persona[]>([]);
const filtros = ref({
  codigoPAF: '',
  run: '',
  codigoAsignatura: '',
  estadoProceso: '',
  calidad: '',
  jerarquia: '',
  semestre: '',
  nombreAsignatura: '',
});
const sortBy = ref('nombres');
const sortOrder = ref('asc');
const warnings = ref<string[]>([]);

const filteredPersonas = computed(() => {
  warnings.value = [];
  let filtered = personas.value.filter(persona => {
    return (
      persona.comentarios_extras?.join(' ').toLowerCase().includes(filtros1.value.nombreAsignatura.toLowerCase() || '') &&
      persona.codigos_extras?.join(' ').toLowerCase().includes(filtros1.value.codigoAsignatura.toLowerCase() || '') &&
      (filtros1.value.estadoProceso ? persona.EstadoProceso.toString() === filtros1.value.estadoProceso : true) &&
      (filtros1.value.calidad ? persona.Calidad === filtros1.value.calidad : true) &&
      persona.IdPaf.toString().toLowerCase().includes(filtros1.value.codigoPAF.toLowerCase()) &&
      persona.Run.toLowerCase().includes(filtros1.value.run.toLowerCase()) &&
      (filtros1.value.jerarquia ? persona.Jerarquia === filtros1.value.jerarquia : true) &&
      persona.SemestrePaf?.toLowerCase().includes(filtros1.value.semestre.toLowerCase() || '')
    );
  });

  filtered = filtered.map(persona => {
    if (persona.bandera_modificacion === 1) {
      warnings.value.push(`Advertencia: La PAF ${persona.IdPaf} ha sido modificada.`);
      persona.rowClass = 'modified-row'; // Marca la fila como modificada
    } else if (persona.bandera_modificacion === 2) {
      warnings.value.push(`Advertencia: La PAF ${persona.IdPaf} ha sido eliminada.`);
      persona.rowClass = 'deleted-row'; // Marca la fila como eliminada
    }

    if (persona.codigo_modificacion === 1) {
      warnings.value.push(`${persona.descripcion_modificacion}.`);
    }
    return persona;
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

// Método para manejar la notificación a Tabla.vue
const handleRowStatusChanged = (persona: Persona) => {
  // Emitir el cambio de estado (modificación o eliminación)
};

onMounted(async () => {
  try {
    const response = await $axios.get('/api/paf-en-linea/historial');
    console.log('Personas:', response.data);
    personas.value = response.data;
    semestreMasActual.value = response.data
  .map((item: Persona) => item.SemestrePaf) // Extrae el semestre de cada objeto
  .sort((a: string, b: string) => {
    // Convierte los semestres en fechas para poder compararlos
    const [mesA, anioA] = a.split('-');
    const [mesB, anioB] = b.split('-');

    // Compara año y mes
    if (parseInt(anioA) === parseInt(anioB)) {
      return parseInt(mesB) - parseInt(mesA); // Ordena por mes descendente
    }
    return parseInt(anioB) - parseInt(anioA); // Ordena por año descendente
  })[0];
  if (sessionStorage.getItem('codigoPAF_h')) {
    codigoPAF_filtro.value = sessionStorage.getItem('codigoPAF_h') || '';
  }
  if (sessionStorage.getItem('run_h')) {
    run_filtro.value = sessionStorage.getItem('run_h') || '';
  }
  if (sessionStorage.getItem('codigoAsignatura_h')) {
    codigoAsignatura_filtro.value = sessionStorage.getItem('codigoAsignatura_h') || '';
  }
  if (sessionStorage.getItem('semestre_h') || sessionStorage.getItem('semestre_h') == '') {
    semestres_filtro.value = sessionStorage.getItem('semestre_h') || '';
  }
  else{
    semestres_filtro.value = semestreMasActual.value;
  }
  if (sessionStorage.getItem('estadoProceso_h') || sessionStorage.getItem('estadoProceso_h') == '') {
    estado_filtro.value = sessionStorage.getItem('estadoProceso_h') || '';
  }
  else{
    estado_filtro.value = 'C1D';
  }
  if (sessionStorage.getItem('calidad_h')) {
    calidad_filtro.value = sessionStorage.getItem('calidad_h') || '';
  }
  if (sessionStorage.getItem('nombreAsignatura_h')) {
    nombreAsignatura_filtro.value = sessionStorage.getItem('nombreAsignatura_h') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMenor_h')) {
    nombreUnidadMenor_filtro.value = sessionStorage.getItem('nombreUnidadMenor_h') || '';
  }
  if (sessionStorage.getItem('nombreUnidadMayor_h')) {
    nombreUnidadMayor_filtro.value = sessionStorage.getItem('nombreUnidadMayor_h') || '';
  }
  applyFilters();
  } catch (error) {
    console.error('Error al obtener personas:', error);
  }
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

const semestres = computed(() => {
  // Ordenar los semestres y eliminar duplicados
  const semestresOrdenados = [...new Set(personas.value.map(p => p.SemestrePaf))]
    .sort((a, b) => (String(a) || '').localeCompare(String(b) || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return semestresOrdenados; // O retorna [ultimoSemestre] si solo necesitas el último
});

const nombreAsig = computed(() => {
  const nombreAsignatura = [
    ...new Set(personas.value.flatMap(p => p.comentarios_extras))
  ].sort((a, b) => (String(a) || '').localeCompare(String(b) || ''));

  return nombreAsignatura;
});


const Excel = async () => {
  try {
    const response = await $axios.get("/api/paf-en-linea/historial/importarExcel", {
      responseType: "blob", // Importante para recibir el archivo correctamente
    });

    // Crear un enlace de descarga
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", "historial.xlsx"); // Nombre del archivo
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  } catch (error) {
    console.error("Error al descargar el archivo Excel:", error);
  }
};


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
  sessionStorage.setItem('codigoPAF_h', filtros1.value.codigoPAF);
  sessionStorage.setItem('run_h', filtros1.value.run);
  sessionStorage.setItem('codigoAsignatura_h', filtros1.value.codigoAsignatura);
  sessionStorage.setItem('semestre_h', filtros1.value.semestre);
  sessionStorage.setItem('estadoProceso_h', filtros1.value.estadoProceso);
  sessionStorage.setItem('calidad_h', filtros1.value.calidad);
  sessionStorage.setItem('nombreAsignatura_h', filtros1.value.nombreAsignatura);
  sessionStorage.setItem('nombreUnidadMenor_h', filtros1.value.nombreUnidadMenor);
  sessionStorage.setItem('nombreUnidadMayor_h', filtros1.value.nombreUnidadMayor);
}

// Resetear filtros
const emit = defineEmits(['filter', 'sort']);

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
  grid-template-columns: 25% 75%;
  max-width: 100%;
}

/* Advertencias */
.warnings {
  grid-column: 1 / -1;
  background-color: #FFF3CD; /* Fondo amarillo pálido */
  border: 1px solid #FFEBAA; /* Borde amarillo */
  border-radius: 5px;
  padding: 1rem;
  margin-top: 1rem;
  font-family: "Helvetica Neue LT", sans-serif;
  color: #856404; /* Texto amarillo oscuro */
}

.warnings h3 {
  margin: 0 0 0.5rem;
  font-family: "Bebas Neue Pro", sans-serif;
  color: #EA7600; /* Color institucional */
}

.warnings ul {
  margin: 0;
  padding: 0;
  list-style-type: none;
}

.warnings li {
  margin: 0.5rem 0;
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
