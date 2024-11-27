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
      <p class="detalle-text">Promedio de tiempo en estado:<strong>
        {{
          promedioTiempoPorEstado[estadoSeleccionado]
            ? promedioTiempoPorEstado[estadoSeleccionado].toFixed(2)
            : 'N/A'
        }}
      </strong>días</p>
    </div>

    <!-- Gráficos -->
    <div v-if="profesoresChartData && pafChartData && pafPorEstadoChartData && pafPorUnidadMayorChartData" class="grafico-container">
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores con y sin PAF</h4>
        <Pie :data="profesoresChartData" />
      </div>
      <div class="pie-chart">
        <h4 class="subtitulo">Profesores con PAF y Profesores con PAF activas</h4>
        <Pie :data="pafChartData" />
      </div>
      <div class="pie-chart">
        <h4 class="subtitulo">Porcentaje de PAF por estado</h4>
        <Pie :data="pafPorEstadoChartData" />
      </div>
      <div class="bar-chart">
        <h4 class="subtitulo">PAF por Unidad Mayor</h4>
        <Bar :data="pafPorUnidadMayorChartData" />
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


const { $axios } = useNuxtApp();

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, ChartDataLabels, BarElement);

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

const fetchCantidadPersonasSai = async () => {
  try {
    const response = await $axios.get('/estadisticas');
    cantidadPersonasSai.value = response.data.TotalProfesores;
    cantidadPafUnicas.value = response.data.TotalPipelsoftUnicos;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};

const fetchCantidadPafSai = async () => {
  try {
    const response = await $axios.get('/estadisticas/PafActivas');
    console.log('response:', response);
    cantidadPafActivas.value = response.data.conteo;
  } catch (error) {
    console.error('Error al obtener la cantidad de personas del SAI:', error);
  }
};


const fetchCantidadPafPorEstado = async () => {
  try {
    const response = await $axios.get('/estadisticas');
    console.log('response:', response);
    cantidadPafPorEstado.value = response.data.EstadoProcesoCount;
    console.log('cantidadPafPorEstado:', cantidadPafPorEstado.value);
    totalPaf.value = Object.values(response.data.EstadoProcesoCount).reduce((a, b) => a + b, 0);
    totalPorcPaf.value = Object.values(response.data.EstadoProcesoCount).map((value) => ((value / totalPaf.value) * 100).toFixed(2));
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
    const response = await $axios.get('/estadisticas/frecuencia-unidades-mayores');
    const unidadesData = response.data;

    pafPorUnidadMayorChartData.value = {
      labels: Object.keys(unidadesData),
      datasets: [
        {
          label: 'Cantidad de PAF por Unidad Mayor',
          data: Object.values(unidadesData),
          backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350'], // Colores de las barras
        },
      ],
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
    labels: ['Profesores con PAF', 'Profesores sin PAF'],
    datasets: [
      {
        label: 'Cantidad',
        data: [
          ((cantidadPafUnicas.value / cantidadPersonasSai.value) * 100).toFixed(2),
          ((cantidadPersonasSai.value - cantidadPafUnicas.value) / cantidadPersonasSai.value * 100).toFixed(2),
        ],
        backgroundColor: ['#42A5F5', '#EF5350'],
      },
    ],
    plugins: {
      datalabels: commonDatalabelsOptions,
    },
  };

  // Gráfico de PAF por Estado
  pafPorEstadoChartData.value = {
    labels: Object.keys(cantidadPafPorEstado.value), // Estados
    datasets: [
      {
        label: 'Porcentaje de PAF por estado',
        data: totalPorcPaf.value, // Porcentajes calculados
        backgroundColor: ['#66BB6A', '#FFA726', '#AB47BC', '#394049', '#EA7600', '#C8102E', '#42A5F5', '#0db58b', '#f0f0f0', '#76095b'],
      },
    ],
    plugins: {
      datalabels: commonDatalabelsOptions,
    },
  };

  pafChartData.value = {
    labels: ['Profesores con PAF activa', 'Profesores sin PAF'],
    datasets: [
      {
        label: 'Cantidad',
        data: [
          ((cantidadPafActivas.value / cantidadPersonasSai.value) * 100).toFixed(2),
          ((cantidadPersonasSai.value - cantidadPafActivas.value) / cantidadPersonasSai.value * 100).toFixed(2),
        ],
        backgroundColor: ['#42A5F5', '#EF5350'],
      },
    ],
    plugins: {
      datalabels: commonDatalabelsOptions,
    },
  };
};


const mostrarDetalles = (estado) => {
  estadoSeleccionado.value = estado;
};

onMounted(async () => {
  await Promise.all([
    fetchCantidadPersonasSai(),
    fetchCantidadPafPorEstado(),
    fetchPromedioTiempoPorEstado(),
    fetchPafPorUnidadMayor(),
    fetchCantidadPafSai(),
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
  max-width: 500px;
  height: 300px;
}

.bar-chart {
  margin: 2rem;
  max-width: 600px;  /* Aumenté el ancho para las barras */
  height: auto;
}


/* Estados */
.estado-linea {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
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

.estado-A1 {
  background-color: #66BB6A;
}

.estado-A2 {
  background-color: #FFA726;
}

.estado-A3 {
  background-color: #AB47BC;
}

.estado-A9 {
  background-color: #394049;
}

.estado-B1 {
  background-color: #EA7600;
}

.estado-B9 {
  background-color: #C8102E;
}

.estado-C1D {
  background-color: #42A5F5;
}

.estado-C9D {
  background-color: #0db58b;
}

.estado-F1 {
  background-color: #f0f0f0;
}

.estado-F9 {
  background-color: #76095b;
}

.detalles-container {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
</style>
