<template>
    <div>
      <h1 class="titulo-principal">Datos SAI y PAF</h1>
      <h1>Semestre seleccionado: <select v-model="semestreSeleccionado" @change="obtenerSemestres">
      <option v-for="semestre in semestresDisponibles" :key="semestre" :value="semestre">
        {{ semestre }}
      </option>
    </select>
    </h1>
    <div v-if="unidadSeleccionada" class="unidad-seleccionada">
      <h4 class="cantidad-text">Unidad Mayor Seleccionada: {{ unidadSeleccionada }}</h4>
      <!-- Botón para recargar los datos -->
      <button @click="recargarPagina" class="btn-recargar">Recargar Datos Iniciales</button>
    </div>
    <br />
    <!-- Cantidades -->
    <p class="cantidad-text">Cantidad de PAF: <strong>{{ totalPafPipelsoft }}</strong></p>
    <p class="cantidad-text">Cantidad de docentes: <strong>{{ cantidadPersonasSai }}</strong></p>
    <p class="cantidad-text">Cantidad de docentes con PAF: <strong>{{ cantidadPafUnicas }}</strong></p>
    <p class="cantidad-text">Cantidad de PAF activas: <strong>{{ cantidadPafActivas }}</strong></p>
    <p class="cantidad-text">Cantidad de PAF inactivas: <strong>{{ (totalPafPipelsoft - cantidadPafActivas) }}</strong></p>
    <p class="cantidad-text">se entiende por activa que no se encuente rechazada, anulada o sin solicitar</p>
    <br />
    <p class="cantidad-text"> Estados de la PAF:</p>
    <br />
    <!-- Estado de avance -->
    <div class="estado-linea">
      <template v-for="(cantidad, estado, index) in cantidadPafPorEstado" :key="'estado-' + estado">
        <div
          :class="['estado-rectangulo', `estado-${estado}`, { 'estado-seleccionado': estado === estadoSeleccionado }]"
          @click="mostrarDetalles(estado)"
        >
          Estado {{ estado }}
        </div>
        <!-- Agregar flecha, excepto después del último estado -->
        <span v-if="index < Object.keys(cantidadPafPorEstado).length - 1" class="estado-flecha">➔</span>
      </template>
    </div>
    
    <!-- Detalles desplegables por estado -->
    <div v-if="estadoSeleccionado !== null" class="detalles-container">
      <h4 class="subtitulo">Detalles del Estado {{ estadoSeleccionado }}</h4>
      <p class="detalle-text">Cantidad de PAF en este estado: <strong>{{ cantidadPafPorEstado[estadoSeleccionado] }}</strong></p>
    </div>

    <!-- Gráficos -->
    <div v-if="profesoresChartData && pafChartData" class="grafico-container">
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores con PAF y sin PAF</h4>
        <Pie :data="profesoresChartData" :options="profesoresChartData.options" />
      </div>
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores sin PAF activas y con PAF activas</h4>
        <Pie :data="pafChartData" :options="pafChartData.options" />
      </div>
    </div>
    <div class="grafico">
      <div v-if="pafPorEstadoChartData" class="pie-chart1">
        <h4 class="subtitulo">Cantidad de PAF por Estado</h4>
        <Pie :data="pafPorEstadoChartData" :options="pafPorEstadoChartData.options" />
      </div>
    </div>
    <br />
    <br />
    <div v-if="pafPorUnidadMayorChartData" class="bar-chart">
      <h4 class="subtitulo">Cantidad de PAF por Unidad Mayor</h4>
      <Bar 
  :data="pafPorUnidadMayorChartData" 
  :options="pafPorUnidadMayorChartData.options"
/>
    </div>
       <!-- Modal para el gráfico dinámico -->
    <div v-if="modalAbierto === 1" class="modal-overlay">
      <div class="modal-content">
        <h3 class="modal-title">Detalles de:</h3>
        <Bar v-if="graficoModalData" :data="graficoModalData" />
        <div v-if="mostrarBotonModal === true" class="boton-lista">
          <button @click="obtenerListaProfesoresSinPaf" class="btn-lista">Obtener lista de profesores</button>
          <br />
        </div>
        <button @click="cerrarModal" class="modal-close-button">Cerrar</button>
      </div>
    </div>
        <div  v-if="modalAbierto === 2" class="modal-overlay">
          <div class="modal-content">
      <div v-if="listaProfesores.length > 0 || listaProfesores !== null" class="lista-profesores">
        <h2>Contratos Relacionados</h2>
        <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead>
        <tr>
          <th class="col-medium">
            Run
          </th>
          <th class="col-small">
            Sección
          </th>
          <th class="col-large">
            Nombre de Asignatura
          </th>
          <th class="col-small">
            Codigo de Asignatura
          </th>
          <th class="col-medium">
            Semestre
          </th>
          <th class="col-small">
            Bloque
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(contrato, index) in paginatedData" :key="index">
          <td>{{ contrato.run }}</td>
          <td>{{ contrato.seccion }}</td>
          <td>{{ contrato.nombre_asignatura }}</td>
          <td>{{ contrato.codigo_asignatura }}</td>
          <td>{{ contrato.semestre }}</td>
          <td>{{ contrato.bloque }}</td>
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
        <button @click="cerrarModal" class="modal-close-button">Cerrar</button>
      </div>
      </div>
        </div>
      </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Pie } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, BarElement } from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import { Bar } from 'vue-chartjs'; // Importar gráfico de barras
import { isConstructorDeclaration } from 'typescript';

const { $axios } = useNuxtApp();

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, ChartDataLabels, BarElement);

const mostrarModal = ref(false); // Controla si el modal está visible
const graficoModalData = ref(null); // Datos para el gráfico en el modal
const cantidadPersonasSai = ref(0);
const totalPafPipelsoft = ref(0);
const cantidadPafActivas = ref(0);
const cantidadPafUnicas = ref(0);
const cantidadPafPorEstado = ref({});
const promedioTiempoPorEstado = ref({});
const profesoresChartData = ref(null);
const pafPorEstadoChartData = ref(null);
const pafChartData = ref(null);
const estadoSeleccionado = ref(null);
const totalPaf = ref(0);
const pafPorUnidadMayorChartData = ref(null);
const totalPorcPaf = ref([]);
const unidadSeleccionada = ref(null); // Unidad seleccionada
const detalleUnidadSeleccionada = ref(null); // Detalles de la unidad seleccionada
const semestreSeleccionado = ref(''); // Almacenar semestre seleccionado
const semestresDisponibles = ref([]); // Semestres disponibles de la API
const mostrarBotonModal = ref(false);
const valor = ref(false);
const listaProfesores = ref([]);

const modalAbierto = ref(null);

    const abrirModal = (modal) => {
      modalAbierto.value = modal;
    };

// Paginación
const currentPage = ref(1);
const itemsPerPage = 8; // Número de elementos por página

const sortData = (key) => {
  if (sortBy.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortBy.value = key;
    sortOrder.value = 'asc';
  }
  currentPage.value = 1; // Resetear a la primera página
};

// Computed para la paginación
const totalPages = computed(() => {
  return Math.ceil(listaProfesores.value.length / itemsPerPage);
});

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return listaProfesores.value.slice(start, end);
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

const goToPage = (page) => {
  if (typeof page === 'number') {
    currentPage.value = page;
  }
};

const cargarProfesores = async () => {
  try {
    // Llamada a la API con parámetros de paginación
    const response = await $axios.get('/api/paf-en-linea/profesores/NoContrato', {
      params: {
        page: paginaActual.value, // Página actual
        size: tamanoPagina.value, // Tamaño de la página
      },
    });

    profesores.value = response.data.profesores; // Datos de la página actual
    totalProfesores.value = response.data.total; // Total de registros (viene de la API)
  } catch (error) {
    console.error('Error al cargar la lista de profesores:', error);
  }
};

const cambiarPagina = (nuevaPagina) => {
  if (nuevaPagina > 0 && nuevaPagina <= paginasTotales.value) {
    paginaActual.value = nuevaPagina;
    cargarProfesores(); // Cargar los datos de la nueva página
  }
};
// Función para obtener los semestres de la respuesta de la API
const obtenerSemestres = async () => {
  try {
    unidadSeleccionada.value = null;
    const response = await $axios.get(`/api/paf-en-linea/pipelsoft/contratos`);
    // Extraer los semestres únicos de la respuesta
    const semestres = response.data.map(item => item.PipelsoftData.Semestre);
    // Filtrar y ordenar los semestres en base al año (YY) y mes (MM)
    const semestresUnicos = [...new Set(semestres)]
    .filter(semestre => typeof semestre === 'string' && semestre.includes('-')) // Filtrar valores válidos
    .sort((a, b) => {
      const [monthA, yearA] = a.split('-'); // Obtener mes y año de "MM-YY"
      const [monthB, yearB] = b.split('-');
      
      // Primero ordenar por año (YY) de forma ascendente
      if (yearA !== yearB) {
        return yearA.localeCompare(yearB);
      }
      // Si los años son iguales, ordenar por mes (MM) de forma ascendente
      return monthA.localeCompare(monthB);
    });

    semestresDisponibles.value = semestresUnicos;

    // Establecer semestre por defecto como el más reciente, solo si no hay una selección previa
    if (semestreSeleccionado.value === '') {
      semestreSeleccionado.value = semestresUnicos[semestresUnicos.length - 1];  // Seleccionar el semestre más reciente (último semestre)
    }
    await Promise.all([
      fetchCantidadPersonasSai(),
      fetchCantidadPafSai(),
      fetchCantidadPafPorEstado(),
      fetchPafPorUnidadMayor(),
      ]);
      configurarGraficos();
  } catch (error) {
    console.error('Error al obtener los semestres:', error);
  }
};
// Las demás funciones para obtener datos y configurar gráficos son iguales.
const fetchCantidadPersonasSai = async () => {
  try {
    if (!semestreSeleccionado.value) {
      obtenerSemestres();
    }
    const response = await $axios.get(`/api/paf-en-linea/estadisticas/${semestreSeleccionado.value}`);
    totalPafPipelsoft.value = response.data.total_pipelsoft;
    cantidadPersonasSai.value = response.data.total_profesores;
    cantidadPafUnicas.value = response.data.total_pipelsoft_unicos;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};

const mostrarDetalles = (estado) => {
    estadoSeleccionado.value = estado;
};

const recargarPagina = () => {
  window.location.reload(); // Recarga la página completa
};

const fetchCantidadPafSai = async () => {
  try {
    if (!semestreSeleccionado.value) {
      obtenerSemestres();
    }
    const response = await $axios.get(`/api/paf-en-linea/estadisticas/PafActivas/${semestreSeleccionado.value}`);
    cantidadPafActivas.value = response.data.conteo;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};

const cerrarModal = () => {
  modalAbierto.value = null;
  mostrarModal.value = false;
  graficoModalData.value = null; // Limpiar datos del gráfico
};

const fetchPafPorUnidadMayor = async () => {
  try {
    if (!semestreSeleccionado.value) {
      obtenerSemestres();
    }
    const response = await $axios.get(`/api/paf-en-linea/estadisticas/frecuencia-unidades-mayores/${semestreSeleccionado.value}`);
    const unidadesData = response.data;
    pafPorUnidadMayorChartData.value = {
      labels: Object.keys(unidadesData),
      datasets: [
        {
          label: 'Cantidad de PAF por Unidad Mayor',
          data: Object.values(unidadesData),
          backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65', '#FF7043'],
        },
      ],
      options: {
        responsive: true,
        plugins: {
          tooltip: {
            enabled: true,
          },
        },
        onClick: async (event, elements) => {
          if (elements.length === 0) return;

          const index = elements[0].index;
          const label = pafPorUnidadMayorChartData.value.labels[index];
          const value = pafPorUnidadMayorChartData.value.datasets[0].data[index];

          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }

          unidadSeleccionada.value = label;
          detalleUnidadSeleccionada.value = value;

          const response = await $axios.get(`/api/paf-en-linea/estadisticas/unidad-mayor/unidades-menores-frecuencia/${label}/${semestreSeleccionado.value}`);
          const unidadesData = response.data;
          graficoModalData.value = {
            labels: Object.keys(unidadesData),
            datasets: [
              {
                label: 'Cantidad de PAF por Unidad Menor',
                data: Object.values(unidadesData),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          modalAbierto.value = 1;
          mostrarModal.value = true;
          const response1 = await $axios.get(`/api/paf-en-linea/estadisticas/unidad-mayor/${label}/${semestreSeleccionado.value}`);
          cantidadPersonasSai.value = response1.data.total_profesores;
          totalPafPipelsoft.value = response1.data.total_pipelsoft;
          cantidadPafUnicas.value = response1.data.total_pipelsoft_unicos;
          const response2 = await $axios.get(`/api/paf-en-linea/estadisticas/pafActivas/unidad-mayor/${label}/${semestreSeleccionado.value}`);
          cantidadPafActivas.value = response2.data.totalRegistros;
          const estadoProcesoCount = response1.data.estado_proceso_count;

          const normalizedEstadoProcesoCount = Object.fromEntries(
            Object.entries(estadoProcesoCount).map(([key, value]) => [
              key.replace(/\s+/g, '-').replace(/\./g, '-'),
              value
            ])
          );

          const ordenCorrecto = [
            'Sin-Solicitar',
            'Enviada-al-Interesado',
            'Enviada-al-Validador',
            'Aprobada-por-Validador',
            'Rechazada-por-Validador',
            'Aprobada-por-Dir--Pregrado',
            'Rechazada-por-Dir--de-Pregrado',
            'Aprobada-por-RRHH',
            'Rechazada-por-RRHH',
            'Anulada'
          ];

          cantidadPafPorEstado.value = Object.fromEntries(
            ordenCorrecto.map(key => [key, normalizedEstadoProcesoCount[key]])
          );

          totalPaf.value = Object.values(cantidadPafPorEstado.value).reduce((a, b) => a + b, 0);
          totalPorcPaf.value = Object.values(cantidadPafPorEstado.value).map((value) =>
            ((value / totalPaf.value) * 100).toFixed(2)
          );
          configurarGraficos();
        },
      },
    };
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF por unidad mayor:', error);
    if (error.message === 'El label está vacío. No se puede realizar la consulta.') {
      // Maneja el error si el label está vacío (opcionalmente mostrar un mensaje al usuario)
      alert(error.message);
    }
  }
};

const fetchCantidadPafPorEstado = async () => {
  try {
    const response = await $axios.get(`/api/paf-en-linea/estadisticas/${semestreSeleccionado.value}`);
    const estadoProcesoCount = response.data.estado_proceso_count;

    const normalizedEstadoProcesoCount = Object.fromEntries(
      Object.entries(estadoProcesoCount).map(([key, value]) => [
        key.replace(/\s+/g, '-').replace(/\./g, '-'),
        value
      ])
    );

    const ordenCorrecto = [
      'Sin-Solicitar',
      'Enviada-al-Interesado',
      'Enviada-al-Validador',
      'Aprobada-por-Validador',
      'Rechazada-por-Validador',
      'Aprobada-por-Dir--Pregrado',
      'Rechazada-por-Dir--de-Pregrado',
      'Aprobada-por-RRHH',
      'Rechazada-por-RRHH',
      'Anulada'
    ];

    cantidadPafPorEstado.value = Object.fromEntries(
      ordenCorrecto.map(key => [key, normalizedEstadoProcesoCount[key]])
    );

    totalPaf.value = Object.values(cantidadPafPorEstado.value).reduce((a, b) => a + b, 0);
    totalPorcPaf.value = Object.values(cantidadPafPorEstado.value).map((value) =>
      ((value / totalPaf.value) * 100).toFixed(2)
    );
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF por estado:', error);
  }
};

const fetchPromedioTiempoPorEstado = async () => {
  try {
    promedioTiempoPorEstado.value = {
      "A1": 10.5,
      "A3": 8.3,
      "F1": 12.1,
    };
  } catch (error) {
    console.error('Error al obtener el promedio de tiempo por estado:', error);
  }
};

const obtenerListaProfesoresSinPaf = async () => {
  try {
    const response = await $axios.get('/api/paf-en-linea/profesores/NoContrato');
    console.log("asdasdasdasdasd2", response);
    listaProfesores.value = response.data.profesores_sin_contrato;
    console.log('Lista de profesores sin PAF:', listaProfesores.value);
    valor.value = true;
    modalAbierto.value = 2;
  } catch (error) {
    console.error('Error al obtener la lista de profesores sin PAF:', error);
  }
};

const configurarGraficos = () => {
  const commonDatalabelsOptions = {
    formatter: (value) => (parseFloat(value) > 0 ? `${value}%` : ''),
    color: '#ffffff',
    font: {
      weight: 'bold',
    },
    align: 'center',
    anchor: 'center',
  };
  profesoresChartData.value = {
    labels: [
    `Profesores con PAF (${cantidadPafUnicas.value})`,
    `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`
  ],
    datasets: [
      {
        label: 'Porcentaje de PAF', 
        data: [
          ((cantidadPafUnicas.value / cantidadPersonasSai.value) * 100).toFixed(2),
          ((cantidadPersonasSai.value - cantidadPafUnicas.value) / cantidadPersonasSai.value * 100).toFixed(2),
        ],
        backgroundColor: ['#42A5F5', '#EF5350'],
      },
    ],
    options: {
        responsive: true,
        plugins: {
          datalabels: commonDatalabelsOptions,
          tooltip: {
            enabled: true,
          },
        },
        onClick: async (event, elements) => {
          if (elements.length === 0) return;

          const index = elements[0].index;
          const label = profesoresChartData.value.labels[index];
          let mostrarBoton = false;
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          console.log("label", label);
          let response = null;
          let unidadesData = null;
          let labelNuevo = "";
          if(unidadSeleccionada.value === null) {
            if (label === `Profesores con PAF (${cantidadPafUnicas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-mayores/cant_profesores/${semestreSeleccionado.value}`);
            console.log("response", response);
            unidadesData = response.data.unidadesMayores;
            labelNuevo = 'Cantidad de profesores con PAF por Unidad Mayor';
            } else if (label === `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-mayores/sin_profesores/${semestreSeleccionado.value}`);
            console.log("response", response);
            unidadesData = response.data;
            labelNuevo = 'Cantidad de profesores sin PAF por Unidad Mayor';
            mostrarBoton = true;
            }
          } else {
            if (label === `Profesores con PAF (${cantidadPafUnicas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-con-profesores-activos/8_1/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            unidadesData = response.data;
            labelNuevo = 'Cantidad de profesores con PAF por Unidad Menor';
            } else if (label === `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-sin-profesores-8-2/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            unidadesData = response.data.unidades;
            labelNuevo = 'Cantidad de profesores sin PAF por Unidad Menor';
            mostrarBoton = true;
            }
          }
          graficoModalData.value = {
            labels: Object.keys(unidadesData),
            datasets: [
              {
                label: labelNuevo,
                data: Object.values(unidadesData),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          mostrarModal.value = true;
          modalAbierto.value = 1;
          mostrarBotonModal.value = mostrarBoton;
        },
      },
  };

  pafPorEstadoChartData.value = {
    labels: Object.keys(cantidadPafPorEstado.value),
    datasets: [
      {
        label: 'Porcentaje de PAF por estado',
        data: totalPorcPaf.value,
        backgroundColor: ['#66BB6A', '#FFA726', '#AB47BC', '#EA7600', '#C8102E', '#42A5F5', '#0db58b', '#6d8a0c', '#76095b', '#394049'],
      },
    ],
    options: {
        responsive: true,
        plugins: {
          datalabels: commonDatalabelsOptions,
          tooltip: {
            enabled: true,
          },
        },
        onClick: async (event, elements) => {
          if (elements.length === 0) return;
          mostrarBotonModal.value = false;
          const index = elements[0].index;
          const label = pafPorEstadoChartData.value.labels[index].replace(/--/g, '. ').replace(/-/g, ' ');
          let labelNuevo = "";
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          let response1 = null;
          if(unidadSeleccionada.value === null) {
              response1 = await $axios.get(`/api/paf-en-linea/estadisticas/profesores/estado/${encodeURIComponent(label)}/${semestreSeleccionado.value}`);
              labelNuevo = "Cantidad de PAF por Unidad Mayor";
            } else {
              response1 = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores/${encodeURIComponent(label)}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
              labelNuevo = "Cantidad de PAF por Unidad Menor";
            }
          const unidadesData1 = response1.data;
          graficoModalData.value = {
            labels: Object.keys(unidadesData1),
            datasets: [
              {
                label: labelNuevo,
                data: Object.values(unidadesData1),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          modalAbierto.value = 1;
          mostrarModal.value = true;
        },
      },
  };
  pafChartData.value = {
    labels: [`Profesores con PAF activas (${cantidadPafActivas.value})`, 
             `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`],
    datasets: [
      {
        label: 'Porcentaje de Profesores por estado',
        data: [
          ((cantidadPafActivas.value / totalPafPipelsoft.value) * 100).toFixed(2),
          ((totalPafPipelsoft.value - cantidadPafActivas.value) / totalPafPipelsoft.value * 100).toFixed(2),
        ],
        backgroundColor: ['#42A5F5', '#EF5350'],
      },
    ],
    options: {
        responsive: true,
        plugins: {
          datalabels: commonDatalabelsOptions,
          tooltip: {
            enabled: true,
          },
        },
        onClick: async (event, elements) => {
          if (elements.length === 0) return;

          const index = elements[0].index;
          const label = pafChartData.value.labels[index];
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          mostrarBotonModal.value = false;
          let response2 = null;
          let unidadesData2 = null;
          let labelNuevo = "";
          if(unidadSeleccionada.value === null) {
            if (label === `Profesores con PAF activas (${cantidadPafActivas.value})`) {

            response2 = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-mayores/profesores-filtrados/${semestreSeleccionado.value}`);
            labelNuevo = 'Cantidad de profesores con PAF activas por Unidad Mayor';
            } else if (label === `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`) {
            response2 = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-mayores/profesores-codestado/${semestreSeleccionado.value}`);
            labelNuevo = 'Cantidad de profesores sin PAF activas por Unidad Mayor';
            } 
            unidadesData2 = response2.data;
            } else {
            if (label === `Profesores con PAF activas (${cantidadPafActivas.value})`) {
            response2 = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-sin-profesores/8_3/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            unidadesData2 = response2.data.unidades;
            labelNuevo = 'Cantidad de profesores con PAF activas por Unidad Menor';
            } else if (label === `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`) {
            response2 = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-con-profesores-paf-activos/8_4/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            unidadesData2 = response2.data;
            labelNuevo = 'Cantidad de profesores sin PAF activas por Unidad Menor';
          }
          }
          graficoModalData.value = {
            labels: Object.keys(unidadesData2),
            datasets: [
              {
                label: labelNuevo,
                data: Object.values(unidadesData2),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          modalAbierto.value = 1;
          mostrarModal.value = true;
        },
      },
  };
};

onMounted(async () => {
  const response = await $axios.get(`/api/paf-en-linea/pipelsoft/contratos`);
    // Extraer los semestres únicos de la respuesta
    const semestres = response.data.map(item => item.PipelsoftData.Semestre);
    // Filtrar y ordenar los semestres en base al año (YY) y mes (MM)
    const semestresUnicos = [...new Set(semestres)].sort((a, b) => {
      const [monthA, yearA] = a.split('-'); // Obtener mes y año de "MM-YY"
      const [monthB, yearB] = b.split('-');
      
      // Primero ordenar por año (YY) de forma ascendente
      if (yearA !== yearB) {
        return yearA.localeCompare(yearB);
      }
      // Si los años son iguales, ordenar por mes (MM) de forma ascendente
      return monthA.localeCompare(monthB);
    });

    semestresDisponibles.value = semestresUnicos;

    // Establecer semestre por defecto como el más reciente, solo si no hay una selección previa
    if (semestreSeleccionado.value === '') {
      semestreSeleccionado.value = semestresUnicos[semestresUnicos.length - 1];  // Seleccionar el semestre más reciente (último semestre)
    }
  await Promise.all([

    fetchCantidadPersonasSai(),
    fetchCantidadPafSai(),
    fetchCantidadPafPorEstado(),
    fetchPromedioTiempoPorEstado(),
    fetchPafPorUnidadMayor(),
  ]);

  configurarGraficos();
});
</script>

<style scoped>
/* Estilo de textos */
.titulo-principal {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: #EA7600;
  font-family: "Bebas Neue Pro", sans-serif;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  max-width: 900px;
  width: 90%;
  text-align: center;
}

.modal-title {
  font-size: 1.5em;
  margin-bottom: 10px;
}

.modal-close-button {
  margin-top: 15px;
  padding: 10px 20px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.modal-close-button:hover {
  background: #d32f2f;
}

.subtitulo {
  margin-top: 1.5rem;
  color: #394049;
  font-family: "Bebas Neue Pro", sans-serif;
  text-align: center;
}

.cantidad-text {
  font-size: 1.2rem;
  color: #394049;
  font-family: "Helvetica Neue LT", sans-serif;
}

.detalle-text {
  font-size: 1rem;
  color: #394049;
  text-align: center;
}

/* Contenedor de gráficos */
.grafico-container {
  display: flex;
  justify-content: space-between;
  gap: 2rem;
  flex-wrap: wrap;
}

.pie-chart {
  margin: 2rem;
  max-width: 500px;
  height: 300px;
}
.grafico {
  display: flex;
  justify-content: center;  
  align-items: center;    
  height: 500px;            
  width: 100%;              
}

.pie-chart1 {
  margin: 2rem;
  width: 80%; 
  max-width: 500px; 
  height: auto;
}

.bar-chart {
  margin: 2rem;
  max-width: 100%;  
}

.estado-linea {
  display: flex;
  gap: 0.2rem;
  margin-bottom: 2.5rem;
  max-width: 100%;
  overflow-x: auto;
}

.estado-rectangulo {
  cursor: pointer;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  color: #ffffff;
  background-color: #394049;
}

.estado-seleccionado {
  border: 2px solid white;
  background-color: #333333;
}

.estado-Anulada {
  background-color: #394049; 
}

.estado-Aprobada-por-Dir--Pregrado {
  background-color: #76095b; 
}

.estado-Aprobada-por-RRHH {
  background-color: #6d8a0c; 
}

.estado-Aprobada-por-Validador {
  background-color: #0db58b;
}

.estado-Enviada-al-Interesado {
  background-color: #42A5F5; 
}

.estado-Enviada-al-Validador {
  background-color: #C8102E; 
}

.estado-Rechazada-por-Dir--de-Pregrado {
  background-color: #EA7600; 
}

.estado-Rechazada-por-RRHH {
  background-color: #AB47BC;
}

.estado-Rechazada-por-Validador {
  background-color: #FFA726; 
}

.estado-Sin-Solicitar {
  background-color: #66BB6A; 
}

.detalles-container {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.estado-flecha {
  margin: 0 0rem;
  font-size: 1rem;
  color: #394049;
  align-self: center;
}

  .btn-recargar {
    background-color: #f07115;
    color: white;
    border: none;
    padding: 10px 20px;
    cursor: pointer;
    font-size: 16px;
    margin-top: 20px;
    border-radius: 5px;
    justify-content: center;
    text-align: center;
  }

  .btn-recargar:hover {
    background-color: #e51e1e;
  }
  .btn-lista {
  display: inline-block;
  padding: 8px 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-align: center;
  text-decoration: none;
  color: #fcfcfc;
  background-color: #4CAF50;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.btn-lista:hover {
  background-color: #307d33;
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

/* Ancho específico para columnas */
.col-small {
  width: 80px; /* Ancho menor para columnas pequeñas */
}

.col-medium {
  width: 150px; /* Ancho medio para columnas estándar */
}

.col-large {
  width: 350px; /* Ancho mayor para columnas grandes */
}

.table-container {
  width: 100%;
  padding: 20px;
  background-color: var(--background-color);
}

/* Estilos para la tabla */
table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #030101;
  background-color: #EA7600;
}

tbody td {
  padding: 12px;
  font-size: 0.875rem;
  color: var(--secondary-color);
}
</style>