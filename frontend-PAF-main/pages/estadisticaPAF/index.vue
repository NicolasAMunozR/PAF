<template>
      <!-- Detalles de la unidad seleccionada -->
      <div v-if="unidadSeleccionada" class="unidad-seleccionada">
      <h4 class="subtitulo">Unidad Seleccionada: {{ unidadSeleccionada }}</h4>
      <p style="text-align: center;">Cantidad de PAF en esta unidad: <strong>{{ detalleUnidadSeleccionada }}</strong></p>
      <!-- Botón para recargar los datos -->
      <button @click="recargarPagina" class="btn-recargar">Recargar Datos Iniciales</button>
    </div>
    <br />
    <br />
  <div>
    <h1 class="titulo-principal">Datos SAI y PAF</h1>

    <!-- Cantidades -->
    <p class="cantidad-text">Cantidad de docentes: <strong>{{ cantidadPersonasSai }}</strong></p>
    <p class="cantidad-text">Cantidad de PAF únicas: <strong>{{ cantidadPafUnicas }}</strong></p>
    <br />
    <br />
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
      <p class="detalle-text">Promedio de tiempo en estado:
        <strong>
          {{
            promedioTiempoPorEstado[estadoSeleccionado]
              ? promedioTiempoPorEstado[estadoSeleccionado].toFixed(2)
              : 'N/A'
          }}
        </strong>
        días
      </p>
    </div>

    <!-- Gráficos -->
    <div v-if="profesoresChartData && pafChartData" class="grafico-container">
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores con y sin PAF</h4>
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
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, BarElement } from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import { Bar } from 'vue-chartjs'; // Importar gráfico de barras

const { $axios } = useNuxtApp();

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, ChartDataLabels, BarElement);

const mostrarModal = ref(false); // Controla si el modal está visible
const graficoModalData = ref(null); // Datos para el gráfico en el modal
const cantidadPersonasSai = ref(0);
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

// Las demás funciones para obtener datos y configurar gráficos son iguales.
const fetchCantidadPersonasSai = async () => {
  try {
    const response = await $axios.get('/estadisticas');
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
    
    const response = await $axios.get('/estadisticas/PafActivas');
    cantidadPafActivas.value = response.data.conteo;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};

const cerrarModal = () => {
  mostrarModal.value = false;
  graficoModalData.value = null; // Limpiar datos del gráfico
};

const fetchPafPorUnidadMayor = async () => {
  try {
    const response = await $axios.get('/estadisticas/frecuencia-unidades-mayores');
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

          const response = await $axios.get(`/estadisticas/unidad-mayor/unidades-menores-frecuencia/${label}`);
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
          mostrarModal.value = true;
          const response1 = await $axios.get(`/estadisticas/unidad-mayor/${label}`);
          cantidadPersonasSai.value = response1.data.total_profesores;
          cantidadPafUnicas.value = response1.data.total_pipelsoft_unicos;
          const response2 = await $axios.get(`/estadisticas/pafActivas/unidad-mayor/${label}`);
          cantidadPafActivas.value = response2.data.totalRUNs;
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
    const response = await $axios.get('/estadisticas');
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
    labels: ['Profesores con PAF', 'Profesores sin PAF'],
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
          const value = profesoresChartData.value.datasets[0].data[index];
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          let response = null;
          let unidadesData = null;
          let labelNuevo = "";
          if(unidadSeleccionada.value === null) {
            if (label === 'Profesores con PAF') {
            response = await $axios.get(`/estadisticas/unidades-mayores/cant_profesores`);
            } else if (label === 'Profesores sin PAF') {
            response = await $axios.get(`/estadisticas/unidades-mayores/sin_profesores`);
            }
            unidadesData = response.data.unidadesMayores;
            labelNuevo = 'Cantidad de PAF por Unidad Mayor';
          } else {
            if (label === 'Profesores con PAF') {
            response = await $axios.get(`/estadistica/unidades-menores-con-profesores-activos/8_1/${unidadSeleccionada.value}`);
            } else if (label === 'Profesores sin PAF') {
            response = await $axios.get(`/estadistica/unidades-menores-sin-profesores-8-2/${unidadSeleccionada.value}`);
            }
            unidadesData = response.data;
            labelNuevo = 'Cantidad de PAF por Unidad Menor';
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

          const index = elements[0].index;
          const label = pafPorEstadoChartData.value.labels[index].replace(/--/g, '. ').replace(/-/g, ' ');
          const value = pafPorEstadoChartData.value.datasets[0].data[index];
          if (!label || label.trim() === '') {
            throw new Error('El label está vacío. No se puede realizar la consulta.');
          }
          let response1 = null;
          if(unidadSeleccionada.value === null) {
              response1 = await $axios.get(`/estadisticas/profesores/estado/${encodeURIComponent(label)}`);
          } else {
              response1 = await $axios.get(`/estadistica/unidades-menores/${encodeURIComponent(label)}/${unidadSeleccionada.value}`);
          }
          const unidadesData1 = response1.data;
          graficoModalData.value = {
            labels: Object.keys(unidadesData1),
            datasets: [
              {
                label: 'Cantidad de PAF por Unidad Estado',
                data: Object.values(unidadesData1),
                backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350', '#26C6DA', '#FFEE58', '#8D6E63', '#5C6BC0', '#EC407A', '#78909C', '#9CCC65'],
              },
            ],
          };
          mostrarModal.value = true;
        },
      },
  };

  pafChartData.value = {
    labels: ['Profesores con PAF activas', 'Profesores sin PAF activas'],
    datasets: [
      {
        label: 'Porcentaje de Profesores por estado',
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
          let response2 = null;
          let unidadesData2 = null;
          let labelNuevo = "";
          if(unidadSeleccionada.value === null) {
            if (label === 'Profesores con PAF activas') {

            response2 = await $axios.get(`/estadisticas/unidades-mayores/profesores-filtrados`);
            } else if (label === 'Profesores sin PAF activas') {

            response2 = await $axios.get(`/estadisticas/unidades-mayores/profesores-codestado`);
            }
            unidadesData2 = response2.data;
            labelNuevo = 'Cantidad de PAF por Unidad Mayor';
          } else {
            if (label === 'Profesores con PAF activas') {

            response2 = await $axios.get(`/estadistica/unidades-menores-sin-profesores/8_3/${unidadSeleccionada.value}`);
            unidadesData2 = response2.data.unidades;
            } else if (label === 'Profesores sin PAF activas') {

            response2 = await $axios.get(`/estadistica/unidades-menores-con-profesores-paf-activos/8_4/${unidadSeleccionada.value}`);
            unidadesData2 = response2.data;
            labelNuevo = 'Cantidad de PAF por Unidad Menor';
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
          mostrarModal.value = true;
        },
      },
  };
};

onMounted(async () => {
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

.unidad-seleccionada{
  justify-content: center;
  text-align: center;
}
</style>