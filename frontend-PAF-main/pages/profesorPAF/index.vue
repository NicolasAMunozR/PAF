<template>
    <div class="container">
      <h1>Información del Profesor</h1>
  
      <div v-if="contratos.length > 0">
        <p><strong>Nombre:</strong> {{ contratos[0].PipelsoftData.Nombres }}</p>
        <p><strong>Apellido:</strong> {{ contratos[0].PipelsoftData.PrimerApellido }}</p>
      </div>
  
      <div v-if="contratos.length > 0" class="contratos">
        <h2>Contratos Relacionados</h2>
        <table>
          <thead>
            <tr>
              <th>Código PAF</th>
            <th>Jefatura</th>
            <th>Nombre de Asignatura</th>
            <th>Estado del Proceso</th>
            <th>Fecha de la ultima Actualización de Estado</th>
            <th>Historial de Estados</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(contrato, index) in contratos" :key="index">
            <td>{{ contrato.PipelsoftData.CodigoPAF }}</td>
            <td>{{ contrato.PipelsoftData.Jerarquia }}</td>
            <td>{{ contrato.PipelsoftData.NombreAsignatura }}</td>
            <td>{{ estadoProceso(contrato.PipelsoftData.EstadoProceso) }}</td>
            <td>{{ new Date(contrato.PipelsoftData.FechaUltimaModificacionProceso).toLocaleTimeString() }} {{ new Date(contrato.PipelsoftData.FechaUltimaModificacionProceso).toLocaleDateString() }}</td>
            <td>
                <div v-for="(estado, idx) in contrato.PipelsoftData.historialEstados" :key="idx">
                  <p><strong>{{ estadoProceso(estado.estado) }}</strong>: 
                  {{ calcularTiempoEnEstado(estado.fechaInicio) }} días</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
  
      <div v-else-if="errorMessage" class="error">
        <p>{{ errorMessage }}</p>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted } from "vue";
  import { useRoute } from "vue-router"; // Para acceder a los query params
  const { $axios } = useNuxtApp() as unknown as { $axios: typeof import("axios").default };
  
  const route = useRoute(); // Acceso a la ruta actual
  const run = ref<string>(""); // RUN del profesor
  const profesor = ref<{ nombre: string; apellido: string } | null>(null); // Datos del profesor
  const contratos = ref<any[]>([]); // Lista de contratos
  const errorMessage = ref<string>("");
  
  // Mapeo de estados del proceso a nombres legibles
  const estadoProceso = (estado: number): string => {
    switch (estado) {
      case 1:
        return "Estado 1";
      case 2:
        return "Estado 2";
      case 3:
        return "Estado 3";
      case 4:
        return "Estado 4"
      case 5:
        return "Estado 5"
      case 6:
        return "Estado 6"
      default:
        return "Desconocido";
    }
  };
  
  // Función para calcular el tiempo en cada estado en días
  const calcularTiempoEnEstado = (fechaInicio: string): number => {
    const fechaActual = new Date();
    const fechaInicioEstado = new Date(fechaInicio);
    const diferenciaTiempo = fechaActual.getTime() - fechaInicioEstado.getTime();
    return Math.floor(diferenciaTiempo / (1000 * 3600 * 24)); // Devuelve los días
  }
  
  // Función para buscar la información del profesor y sus contratos
  const fetchProfesorYContratos = async () => {
    if (!run.value) {
      errorMessage.value = "RUN no proporcionado.";
      return;
    }
  
    try {
  
      // Solicitar contratos relacionados
      const contratosResponse = await $axios.get(`/pipelsoft/contratos-run/${run.value}`);
      console.log(contratosResponse);
      if (contratosResponse.data && Array.isArray(contratosResponse.data)) {
        contratos.value = contratosResponse.data;
        errorMessage.value = ""; // Limpiar mensaje de error
      } else {
        errorMessage.value = "No se encontraron contratos para el RUN proporcionado.";
        contratos.value = [];
      }
    } catch (error) {
      errorMessage.value = "Hubo un error al obtener los datos.";
      console.error(error);
    }
  };
  
  // Ejecutar al cargar la página
  onMounted(() => {
    const runFromQuery = route.query.run as string; // Obtener el RUN de los query params
    if (runFromQuery) {
      run.value = runFromQuery; // Asignar el valor al RUN
      fetchProfesorYContratos(); // Buscar información del profesor y contratos
    }
  });
  </script>
  
  <style scoped>
  .container {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
  }
  
  .profesor p {
    font-size: 18px;
    margin: 10px 0;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
  }
  
  th,
  td {
    padding: 10px;
    text-align: left;
    border: 1px solid #ddd;
  }
  
  th {
    background-color: #f2f2f2;
  }
  
  .error {
    color: red;
    font-weight: bold;
  }
  </style>
  