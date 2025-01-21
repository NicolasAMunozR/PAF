<template>
    <div>
      <h1 class="titulo-principal">Datos SAI y PAF</h1>
      <h1>Semestre seleccionado: <select v-model="semestreSeleccionado" @change="obtenerSemestres">
      <option v-for="semestre in semestresDisponibles" :key="semestre" :value="semestre">
        {{ semestre }}
      </option>
    </select>
    </h1>
    <h1 class="subtitulo"> Unidad Mayor: {{ UnidadMayor }}</h1>
    <div v-if="unidadSeleccionada" class="unidad-seleccionada">
      <h4 class="subtitulo">Unidad Menor Seleccionada: {{ unidadSeleccionada }}</h4>
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
        <h4 class="detalle-text">Detalles del Estado {{ estadoSeleccionado }}</h4>
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
        <Bar :data="pafPorUnidadMayorChartData" :options="pafPorUnidadMayorChartData.options" />
      </div>
      <!-- Modal para el gráfico dinámico -->
      <div v-if="mostrarModal" class="modal-overlay">
      <div class="modal-content">
        <h3 class="modal-title">Detalles de {{ unidadSeleccionada }}</h3>
        <Bar v-if="graficoModalData" :data="graficoModalData" />
        <button @click="cerrarModal" class="modal-close-button">Cerrar</button>
      </div>
    </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue';
  import { Pie } from 'vue-chartjs';
  import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale,  BarElement} from 'chart.js';
  import ChartDataLabels from 'chartjs-plugin-datalabels';
  import { Bar } from 'vue-chartjs'; // Importar gráfico de barras

  
  const rut = ref(''); // Sin tipo
  const UnidadMayor = ref('');
  const UnidadMenor = ref('');
  ///pipelsoft/contratos
  //item.PipelsoftData.Semestre,
  const semestreSeleccionado = ref(''); // Almacenar semestre seleccionado
const semestresDisponibles = ref([]); // Semestres disponibles de la API

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
    if (!semestreSeleccionado.value) {
      semestreSeleccionado.value = semestresUnicos[semestresUnicos.length - 1];  // Seleccionar el semestre más reciente (último semestre)
    }
    if(UnidadMenor.value !== ''){
      await Promise.all([
      fetchCantidadPersonasSai1(),
      fetchCantidadPafSai1(),
      fetchCantidadPafPorEstado1(),
      ]);
      configurarGraficos1();
    }
    else{
      await Promise.all([
      fetchCantidadPersonasSai(),
      fetchCantidadPafSai(),
      fetchCantidadPafPorEstado(),
      fetchPafPorUnidadMayor(),
      ]);
      configurarGraficos();
    }
      
  } catch (error) {
    console.error('Error al obtener los semestres:', error);
  }
};
  const { $axios } = useNuxtApp();
  
  ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, ChartDataLabels, BarElement);
  
  const mostrarModal = ref(false); // Controla si el modal está visible
  const graficoModalData = ref(null); // Datos para el gráfico en el modal
  const cantidadPersonasSai = ref(0);
  const cantidadPafActivas = ref(0);
  const cantidadPafUnicas = ref(0);
  const totalPafPipelsoft = ref(0);
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
  const valor = ref(null);

  const fetchCantidadPersonasSai = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      //const response = await $axios.get(`/contratos/${response1.data.unidadMayor}`);
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/unidad-mayor/${UnidadMayor.value}/${semestreSeleccionado.value}`);
      console.log(response.data);
      console.log(UnidadMayor.value);
      cantidadPersonasSai.value = response.data.total_profesores;
      cantidadPafUnicas.value = response.data.total_pipelsoft_unicos;
    } catch (error) {
      console.error('Error al obtener la cantidad de personas del SAI:', error);
    }
  };

  const fetchCantidadPersonasSai1 = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      //const response = await $axios.get(`/contratos/${response1.data.unidadMayor}`);
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/6/${UnidadMayor.value}/${UnidadMenor.value}/${semestreSeleccionado.value}`);
      totalPafPipelsoft.value = response.data.total_pipelsoft;
      cantidadPersonasSai.value = response.data.total_profesores;
      cantidadPafUnicas.value = response.data.total_pipelsoft_unicos;
    } catch (error) {
      console.error('Error al obtener la cantidad de personas del SAI:', error);
    }
  };

const recargarPagina = () => {
  window.location.reload(); // Recarga la página completa
};

const cerrarModal = () => {
  mostrarModal.value = false;
  graficoModalData.value = null; // Limpiar datos del gráfico
};

  const fetchCantidadPafSai = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/pafActivas/unidad-mayor/${UnidadMayor.value}/${semestreSeleccionado.value}`);
      cantidadPafActivas.value = response.data.totalRegistros;
    } catch (error) {
      console.error('Error al obtener la cantidad de personas del SAI:', error);
    }
  };
  const fetchCantidadPafSai1 = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/7/${UnidadMayor.value}/${UnidadMenor.value}/${semestreSeleccionado.value}`);
      cantidadPafActivas.value = response.data.totalRegistros;
    } catch (error) {
      console.error('Error al obtener la cantidad de personas del SAI:', error);
    }
  };
  
  const fetchCantidadPafPorEstado = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      //const response = await $axios.get(`/contratos/${response1.data.unidadMayor}`);
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/unidad-mayor/${UnidadMayor.value}/${semestreSeleccionado.value}`);
      // Ordenar el objeto EstadoProcesoCount para que Sin Solicitar sea el último
      const estadoProcesoCount = response.data.estado_proceso_count;
      // Normalizar las claves
      const normalizedEstadoProcesoCount = Object.fromEntries(
        Object.entries(estadoProcesoCount).map(([key, value]) => [
          key.replace(/\s+/g, '-').replace(/\./g, '-'), // Reemplaza espacios por guiones y puntos por guiones bajos
          value
        ])
      );
      // Definir el orden de las claves
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
      // Ordenar las claves del objeto según el arreglo 'ordenCorrecto'
      const orderedEstadoProcesoCount = Object.fromEntries(
        ordenCorrecto.map(key => [key, normalizedEstadoProcesoCount[key]])
      );
      cantidadPafPorEstado.value = orderedEstadoProcesoCount;
      totalPaf.value = Object.values(cantidadPafPorEstado.value).reduce((a, b) => a + b, 0);
      totalPorcPaf.value = Object.values(cantidadPafPorEstado.value).map((value) =>
        ((value / totalPaf.value) * 100).toFixed(2)
      );
    } catch (error) {
      console.error('Error al obtener la cantidad de PAF por estado:', error);
    }
  };

  const fetchCantidadPafPorEstado1 = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      //const response = await $axios.get(`/contratos/${response1.data.unidadMayor}`);
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/6/${UnidadMayor.value}/${UnidadMenor.value}/${semestreSeleccionado.value}`);
      // Ordenar el objeto EstadoProcesoCount para que Sin Solicitar sea el último
      const estadoProcesoCount = response.data.estado_proceso_count;
      // Normalizar las claves
      const normalizedEstadoProcesoCount = Object.fromEntries(
        Object.entries(estadoProcesoCount).map(([key, value]) => [
          key.replace(/\s+/g, '-').replace(/\./g, '-'), // Reemplaza espacios por guiones y puntos por guiones bajos
          value
        ])
      );
      // Definir el orden de las claves
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
      // Ordenar las claves del objeto según el arreglo 'ordenCorrecto'
      const orderedEstadoProcesoCount = Object.fromEntries(
        ordenCorrecto.map(key => [key, normalizedEstadoProcesoCount[key]])
      );
      cantidadPafPorEstado.value = orderedEstadoProcesoCount;
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
  
  const fetchPafPorUnidadMayor = async () => {
    try {
      if(!semestreSeleccionado.value){
        obtenerSemestres();
      }
      const response = await $axios.get(`/api/paf-en-linea/estadisticas/unidad-mayor/unidades-menores-frecuencia/${UnidadMayor.value}/${semestreSeleccionado.value}`);
      const unidadesData = response.data;
      pafPorUnidadMayorChartData.value = {
        labels: Object.keys(unidadesData),
        datasets: [
          {
            label: 'Cantidad de PAF por Unidad Menor',
            data: Object.values(unidadesData),
            backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350'], // Colores de las barras
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
          const response1 = await $axios.get(`/api/paf-en-linea/estadisticas/6/${UnidadMayor.value}/${label}/${semestreSeleccionado.value}`);
          cantidadPersonasSai.value = response1.data.total_profesores;
          cantidadPafUnicas.value = response1.data.total_pipelsoft_unicos;
          const response2 = await $axios.get(`/api/paf-en-linea/estadisticas/7/${UnidadMayor.value}/${label}/${semestreSeleccionado.value}`);
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
    }
  };
  
  const configurarGraficos = () => {
    const commonDatalabelsOptions = {
      formatter: (value) => (parseFloat(value) > 0 ? `${value}%` : ''), // Mostrar solo si el porcentaje es mayor que 0
      color: '#ffffff', // Color del texto
      font: {
        weight: 'bold',
      },
      align: 'center', // Alinear al centro del sector
      anchor: 'center',
    };
  
    // Gráfico de Profesores con y sin PAF
    profesoresChartData.value = {
      labels: [`Profesores con PAF (${cantidadPafUnicas.value})`,
    `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`],
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
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          let response = null;
          let unidadesData = null;
          let labelNuevo = "";
          if(unidadSeleccionada.value === null) {
            if (label === `Profesores con PAF (${cantidadPafUnicas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-con-profesores-activos/8_1/${UnidadMayor.value}/${semestreSeleccionado.value}`);
            } else if (label === `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`) {
              response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-sin-profesores-8-2/${UnidadMayor.value}/${semestreSeleccionado.value}`);
            }
            unidadesData = response.data;
            labelNuevo = "Cantidad de PAF por Unidad Menor";
          } else {
            if (label === `Profesores con PAF (${cantidadPafUnicas.value})`) {
            // CAMBIAR AQUÍ
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidadesmenores/profesores/${UnidadMayor.value}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            } else if (label === `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`) {
              // CAMBIAR AQUÍ
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidadesmenores/sinprofesores/${UnidadMayor.value}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            }

            unidadesData = response.data;
            labelNuevo = "Cantidad de PAF por Unidad Mayor y Menor"; 
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
        },
      },
    };
  
    // Gráfico de PAF por Estado
    pafPorEstadoChartData.value = {
      labels: Object.keys(cantidadPafPorEstado.value), // Estados
      datasets: [
        {
          label: 'Porcentaje de PAF por estado',
          data: totalPorcPaf.value, // Porcentajes calculados
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

          const index = elements[0].index;
          const label = pafPorEstadoChartData.value.labels[index].replace(/--/g, '. ').replace(/-/g, ' ');
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          let response = null;
          let unidadesData = null;
          if(unidadSeleccionada.value === null) {
              response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores/${encodeURIComponent(label)}/${UnidadMayor.value}/${semestreSeleccionado.value}`);
          } else {
              // CAMBIAR AQUÍ
              response = await $axios.get(`/api/paf-en-linea/estadisticas/unidadesmenores/porcodestadopaf/${encodeURIComponent(label)}/${UnidadMayor.value}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
          }
          unidadesData = response.data;
          graficoModalData.value = {
            labels: Object.keys(unidadesData),
            datasets: [
              {
                label: 'Cantidad de PAF por Unidad Estado',
                data: Object.values(unidadesData),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          mostrarModal.value = true;
        },
      },
    };
  
    pafChartData.value = {
      labels: [`Profesores con PAF activas (${cantidadPafActivas.value})`, 
             `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`],
      datasets: [
        {
          label: 'Porcentaje de PAF',
          data: [
            ((cantidadPafActivas.value / cantidadPersonasSai.value) * 100).toFixed(2),
            ((cantidadPersonasSai.value - cantidadPafActivas.value) / cantidadPersonasSai.value * 100).toFixed(2),
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
          let response = null;
          let unidadesData = null;
          let labelNuevo = "";

          if(unidadSeleccionada.value === null) {
            if (label === `Profesores con PAF activas (${cantidadPafActivas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-sin-profesores/8_3/${UnidadMayor.value}/${semestreSeleccionado.value}`);
            unidadesData = response.data.unidades;
            } else if (label === `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`) {
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidades-menores-con-profesores-paf-activos/8_4/${UnidadMayor.value}/${semestreSeleccionado.value}`);
            unidadesData = response.data;
            } 
            labelNuevo = "Cantidad de PAF por Unidad Menor";
          } else {
            if (label === `Profesores con PAF activas (${cantidadPafActivas.value})`) {
            // CAMBIAR AQUÍ
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidadesmayores/filtradospafactivos/${UnidadMayor.value}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            } else if (label === `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`) {
              // CAMBIAR AQUÍ
            response = await $axios.get(`/api/paf-en-linea/estadisticas/unidadesmenores/filtradospafactivos/${UnidadMayor.value}/${unidadSeleccionada.value}/${semestreSeleccionado.value}`);
            }
            unidadesData = response.data;
            labelNuevo = "Cantidad de PAF por Unidad Mayor y Menor";
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
        },
      },
    };
  };
  
  const configurarGraficos1 = () => {
    const commonDatalabelsOptions = {
      formatter: (value) => (parseFloat(value) > 0 ? `${value}%` : ''), // Mostrar solo si el porcentaje es mayor que 0
      color: '#ffffff', // Color del texto
      font: {
        weight: 'bold',
      },
      align: 'center', // Alinear al centro del sector
      anchor: 'center',
    };
  
    // Gráfico de Profesores con y sin PAF
    profesoresChartData.value = {
      labels: [`Profesores con PAF (${cantidadPafUnicas.value})`,
    `Profesores sin PAF (${cantidadPersonasSai.value - cantidadPafUnicas.value})`],
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
      },
    };
  
    // Gráfico de PAF por Estado
    pafPorEstadoChartData.value = {
      labels: Object.keys(cantidadPafPorEstado.value), // Estados
      datasets: [
        {
          label: 'Porcentaje de PAF por estado',
          data: totalPorcPaf.value, // Porcentajes calculados
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
      },
    };
  
    pafChartData.value = {
      labels: [`Profesores con PAF activas (${cantidadPafActivas.value})`, 
             `Profesores sin PAF activas (${totalPafPipelsoft.value - cantidadPafActivas.value})`],
      datasets: [
        {
          label: 'Porcentaje de PAF',
          data: [
            ((cantidadPafActivas.value / cantidadPersonasSai.value) * 100).toFixed(2),
            ((cantidadPersonasSai.value - cantidadPafActivas.value) / cantidadPersonasSai.value * 100).toFixed(2),
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
      },
    };
  };

  const mostrarDetalles = (estado) => {
    estadoSeleccionado.value = estado;
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
    if (!semestreSeleccionado.value) {
      semestreSeleccionado.value = semestresUnicos[semestresUnicos.length - 1];  // Seleccionar el semestre más reciente (último semestre)
    }

    UnidadMayor.value = sessionStorage.getItem("unidadMayor") || "";
    UnidadMenor.value = sessionStorage.getItem("unidadMenor") || "";
    if(UnidadMenor.value !== ''){
      rut.valueOf = sessionStorage.getItem('rut') || '';
      valor.value = rut.valueOf;
      await Promise.all([
        fetchCantidadPersonasSai1(),
        fetchCantidadPafPorEstado1(),
        fetchCantidadPafSai1(),
      ]);
      configurarGraficos1();
    }
    else{
      rut.valueOf = sessionStorage.getItem('rut') || '';
      valor.value = rut.valueOf;
      await Promise.all([
        fetchCantidadPersonasSai(),
        fetchCantidadPafPorEstado(),
        fetchPromedioTiempoPorEstado(),
        fetchPafPorUnidadMayor(),
        fetchCantidadPafSai(),
      ]);
      configurarGraficos();
    }
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
  
  .subtitulo {
    margin-top: 1.5rem;
    color: #394049;
    font-family: "Bebas Neue Pro", sans-serif;
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
  justify-content: center;  /* Centra horizontalmente */
  align-items: center;      /* Centra verticalmente */
  height: 500px;            /* Hace que el contenedor ocupe el 100% de la altura de la ventana */
  width: 100%;              /* Asegura que el contenedor ocupe el 100% del ancho */
}

.pie-chart1 {
  margin: 2rem;
  width: 80%;  /* Ajusta el ancho para que sea más flexible */
  max-width: 500px;  /* Establece un tamaño máximo */
  height: auto;  /* Asegura que el gráfico se ajuste proporcionalmente */
}
  
  .bar-chart {
    margin: 2rem;
    max-width: 100%;  /* Aumenté el ancho para las barras */
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

  /* Estados */
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
  background-color: #394049; /* Verde claro */
}

.estado-Aprobada-por-Dir--Pregrado {
  background-color: #76095b; /* Naranja claro */
}

.estado-Aprobada-por-RRHH {
  background-color: #6d8a0c; /* Morado claro */
}

.estado-Aprobada-por-Validador {
  background-color: #0db58b; /* Gris oscuro */
}

.estado-Enviada-al-Interesado {
  background-color: #42A5F5; /* Naranja oscuro */
}

.estado-Enviada-al-Validador {
  background-color: #C8102E; /* Rojo */
}

.estado-Rechazada-por-Dir--de-Pregrado {
  background-color: #EA7600; /* Azul claro */
}

.estado-Rechazada-por-RRHH {
  background-color: #AB47BC; /* Verde agua */
}

.estado-Rechazada-por-Validador {
  background-color: #FFA726; /* Verde oliva */
}

.estado-Sin-Solicitar {
  background-color: #66BB6A; /* Morado oscuro */
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
</style>  