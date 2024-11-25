<template>
    <div>
      <h1>Datos SAI y PAF</h1>
      
      <!-- Cantidades -->
      <p>Cantidad de personas en el SAI: {{ cantidadPersonasSai }}</p>
      <p>Cantidad de PAF únicas: {{ cantidadPafUnicas }}</p>
      
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
      <div v-if="estadoSeleccionado !== null">
        <h4>Detalles del Estado {{ estadoSeleccionado }}</h4>
        <p>Cantidad de PAF en este estado: {{ cantidadPafPorEstado[estadoSeleccionado] }}</p>
        <p>Promedio de tiempo en estado (días): {{ promedioTiempoPorEstado[estadoSeleccionado].toFixed(2) }} días</p>
      </div>
    
      <!-- Gráficos -->
      <div v-if="profesoresChartData && pafPorEstadoChartData" class="grafico-container">
        <div class="pie-chart">
          <h4>Profesores con y sin PAF</h4>
          <Pie :data="profesoresChartData" />
        </div>
        <div class="pie-chart">
          <h4>PAF por estado</h4>
          <Pie :data="pafPorEstadoChartData" />
        </div>
        <div class="pie-chart">
          <h4>PAF por Unidad Mayor</h4>
          <Pie :data="pafPorUnidadMayorChartData" />
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue';
  import { Pie } from 'vue-chartjs';
  import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale } from 'chart.js';
  
  const { $axios } = useNuxtApp()
  
  ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale);
  
  // Variables reactivas
  const cantidadPersonasSai = ref(0);
  const cantidadPafUnicas = ref(0);
  const cantidadPafPorEstado = ref({});
  const promedioTiempoPorEstado = ref({});
  const profesoresChartData = ref(null);
  const pafPorEstadoChartData = ref(null);
  const estadoSeleccionado = ref(null);
  const totalPaf = ref(0);
  const pafPorUnidadMayorChartData = ref(null);
  
  // Funciones para obtener datos
  const fetchCantidadPersonasSai = async () => {
    try {
      const response = await $axios.get('/estadisticas')
      cantidadPersonasSai.value = response.data.TotalProfesores;
    } catch (error) {
      console.error('Error al obtener la cantidad de personas del SAI:', error);
    }
  };
  
  const fetchCantidadPafUnicas = async () => {
    try {
      const response = await $axios.get('/estadisticas')
      console.log(response.data)
      cantidadPafUnicas.value = response.data.TotalPipelsoftUnicos
      ;
    } catch (error) {
      console.error('Error al obtener la cantidad de PAF únicas:', error);
    }
  };
  
  const fetchCantidadPafPorEstado = async () => {
    try {
      const response = await $axios.get('/estadisticas')
      cantidadPafPorEstado.value = response.data.EstadoProcesoCount;
      totalPaf.value = Object.values(response.data).reduce((a, b) => a + b, 0);
      
    } catch (error) {
      console.error('Error al obtener la cantidad de PAF por estado:', error);
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
          backgroundColor: ['#42A5F5', '#66BB6A', '#FFA726', '#AB47BC', '#EF5350'], // Colores para las barras
        },
      ],
    };
  } catch (error) {
    console.error('Error al obtener la cantidad de PAF por unidad mayor:', error);
  }
};
  
  const fetchPromedioTiempoPorEstado = async () => {
    try {
      const response = {
        data: {
          1: 10.5, // Promedio en estado 1
          2: 8.3,  // Promedio en estado 2
          3: 12.1, // Promedio en estado 3
        },
      };
      promedioTiempoPorEstado.value = response.data;
    } catch (error) {
      console.error('Error al obtener el promedio de tiempo por estado:', error);
    }
  };
  
  // Configurar datos para los gráficos
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
          backgroundColor: ['#66BB6A', '#FFA726', '#AB47BC', 'black', 'red', 'brown'],
        },
      ],
    };

    
  };
  
  // Función para mostrar detalles del estado seleccionado
  const mostrarDetalles = (estado) => {
    estadoSeleccionado.value = estado;
  };
  
  // Ejecutar lógica al montar el componente
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
  h1 {
    font-size: 2rem;
    margin-bottom: 1rem;
  }
  
  h3 {
    margin-top: 1.5rem;
  }
  
  .pie-chart {
    margin: 2rem;
    max-width: 400px;
    height: auto;
  }
  
  .grafico-container {
    display: flex;
    justify-content: space-between;
    gap: 2rem;
  }
  
  .estado-linea {
    display: flex;
    margin: 1rem 0;
    align-items: center;
    position: relative;
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
    flex-grow: 1; /* Hace que todas las barras sean del mismo tamaño */
    width: 33.33%; /* Asigna un tamaño fijo para cada barra */
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
    background-color: black;
  }

  .estado-5 {
    background-color: red;
  }

  .estado-6 {
    background-color: brown;
  }
  .estado-seleccionado {
    border: 2px solid #FFFFFF;
    background-color: #333333;
  }
  
  .flecha {
    font-size: 2rem;
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    color: #000;
    left: 33.33%; /* Coloca la flecha en el medio de cada estado */
  }
  
  .flecha:nth-child(2) {
    left: 33.33%;
  }
  
  .flecha:nth-child(3) {
    left: 33.33%;
  }
  
  .avance-linea {
    position: relative;
    width: 100%;
    height: 10px;
    background-color: #E0E0E0;
    margin-bottom: 1rem;
  }
  
  .avance-indicator {
    height: 100%;
    background-color: #42A5F5;
    transition: width 0.3s ease;
  }
  </style>
  