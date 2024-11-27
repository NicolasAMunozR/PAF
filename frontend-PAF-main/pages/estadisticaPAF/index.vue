<template>
  <div>
    <h1 class="titulo-principal">Datos SAI y PAF</h1>

    <!-- Cantidades -->
    <p class="cantidad-text">Cantidad de personas en el SAI: <strong>{{ cantidadPersonasSai }}</strong></p>
    <p class="cantidad-text">Cantidad de PAF únicas: <strong>{{ cantidadPafUnicas }}</strong></p>

    <!-- Estado de avance -->
    <div class="estado-linea">
      <div
        v-for="(cantidad, estado, index) in cantidadPafPorEstado"
        :key="'estado-' + estado"
        :class="['estado-rectangulo', `estado-${estado}`, { 'estado-seleccionado': estado === estadoSeleccionado }]"
        @click="mostrarDetalles(estado)"
      >
        Estado {{ estado }}
      </div>
    </div>

    <!-- Detalles desplegables por estado -->
    <div v-if="estadoSeleccionado !== null" class="detalles-container">
      <h4 class="subtitulo">Detalles del Estado {{ estadoSeleccionado }}</h4>
      <p class="detalle-text">Cantidad de PAF en este estado: <strong>{{ cantidadPafPorEstado[estadoSeleccionado] }}</strong></p>
      <p class="detalle-text">Promedio de tiempo en estado: <strong>{{ promedioTiempoPorEstado[estadoSeleccionado].toFixed(2) }}</strong> días</p>
    </div>

    <!-- Gráficos -->
    <div v-if="profesoresChartData && pafPorEstadoChartData" class="grafico-container">
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores con y sin PAF</h4>
        <Pie :data="profesoresChartData" />
      </div>
      <div class="pie-chart">
        <h4 class="subtitulo">PAF por estado</h4>
        <Pie :data="pafPorEstadoChartData" />
      </div>
      <div class="pie-chart">
        <h4 class="subtitulo">PAF por Unidad Mayor</h4>
        <Pie :data="pafPorUnidadMayorChartData" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Pie } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale } from 'chart.js';

const { $axios } = useNuxtApp();

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale);

const cantidadPersonasSai = ref(0);
const cantidadPafUnicas = ref(0);
const cantidadPafPorEstado = ref({});
const promedioTiempoPorEstado = ref({});
const profesoresChartData = ref(null);
const pafPorEstadoChartData = ref(null);
const estadoSeleccionado = ref(null);
const totalPaf = ref(0);
const pafPorUnidadMayorChartData = ref(null);

const fetchCantidadPersonasSai = async () => {
  try {
    const response = await $axios.get('/estadisticas');
    cantidadPersonasSai.value = response.data.TotalProfesores;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};

const fetchCantidadPafUnicas = async () => {
  try {
    const response = await $axios.get('/estadisticas');
    cantidadPafUnicas.value = response.data.TotalPipelsoftUnicos;
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF únicas:', error);
  }
};

const fetchCantidadPafPorEstado = async () => {
  try {
    const response = await $axios.get('/estadisticas');
    cantidadPafPorEstado.value = response.data.EstadoProcesoCount;
    totalPaf.value = Object.values(response.data.EstadoProcesoCount).reduce((a, b) => a + b, 0);
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF por estado:', error);
  }
};

const fetchPromedioTiempoPorEstado = async () => {
  try {
    promedioTiempoPorEstado.value = {
      1: 10.5,
      2: 8.3,
      3: 12.1,
    };
  } catch (error) {
    console.error('Error al obtener el promedio de tiempo por estado:', error);
  }
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
          backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350'],
        },
      ],
    };
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF por unidad mayor:', error);
  }
};

const configurarGraficos = () => {
  profesoresChartData.value = {
    labels: ['Profesores con PAF', 'Profesores sin PAF'],
    datasets: [
      {
        label: 'Cantidad',
        data: [cantidadPafUnicas.value, cantidadPersonasSai.value - cantidadPafUnicas.value],
        backgroundColor: ['#42A5F5', '#EF5350'],
      },
    ],
  };

  pafPorEstadoChartData.value = {
    labels: Object.keys(cantidadPafPorEstado.value),
    datasets: [
      {
        label: 'Cantidad de PAF por estado',
        data: Object.values(cantidadPafPorEstado.value),
        backgroundColor: ['#66BB6A', '#FFA726', '#AB47BC', '#394049', '#EA7600', '#C8102E'],
      },
    ],
  };
};

const mostrarDetalles = (estado) => {
  estadoSeleccionado.value = estado;
};

onMounted(async () => {
  await Promise.all([
    fetchCantidadPersonasSai(),
    fetchCantidadPafUnicas(),
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
  max-width: 400px;
  height: auto;
}

/* Estados */
.estado-linea {
  display: flex;
  margin: 1rem 0;
  align-items: center;
  width: 100%;
}

.estado-rectangulo {
  height: 40px;
  text-align: center;
  color: white;
  font-weight: bold;
  cursor: pointer;
  line-height: 40px;
  margin: 0 5px;
  flex-grow: 1;
  border-radius: 5px;
}

.estado-seleccionado {
  border: 2px solid white;
  background-color: #333333;
}

.estado-1 {
  background-color: #66BB6A;
}

.estado-2 {
  background-color: #FFA726;
}

.estado-3 {
  background-color: #AB47BC;
}

.estado-4 {
  background-color: #394049;
}

.estado-5 {
  background-color: #EA7600;
}

.estado-6 {
  background-color: #C8102E;
}
</style>
