<template>
  <div class="container">
    <h1>Información del Profesor</h1>

    <div v-if="contratos.length > 0" class="profesor">
      <p><strong>Nombre:</strong> {{ contratos[0].PipelsoftData.Nombres }}</p>
      <p><strong>Apellido:</strong> {{ contratos[0].PipelsoftData.PrimerApp }} {{ contratos[0].PipelsoftData.SegundoApp }}</p>
    </div>

    <div v-if="contratos.length > 0" class="contratos">
      <h2>Contratos Relacionados</h2>
      <table>
        <thead>
          <tr>
            <th>Código PAF</th>
            <th>Jerarquia</th>
            <th>Nombre de Asignatura</th>
            <th>Estado del Proceso</th>
            <th>Descripción del Proceso</th>
            <th>Fecha de la Última Actualización de Estado</th>
            <th>Historial de Estados</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(contrato, index) in contratos" :key="index">
            <td>{{ contrato.PipelsoftData.IdPaf }}</td>
            <td>{{ contrato.PipelsoftData.Jerarquia }}</td>
            <td>{{ contrato.PipelsoftData.NombreAsignatura }}</td>
            <td>{{ contrato.PipelsoftData.CodEstado }}</td>
            <td>{{ contrato.PipelsoftData.DesEstado }}</td>
            <td>
              {{ new Date(contrato.PipelsoftData.UpdatedAt).toLocaleDateString() }}
              {{ new Date(contrato.PipelsoftData.UpdatedAt).toLocaleTimeString() }} 
            </td>
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
    </div>

    <div v-else-if="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import("axios").default };

const route = useRoute();
const run = ref<string>("");
const contratos = ref<any[]>([]);
const errorMessage = ref<string>("");

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

const calcularTiempoEnEstado = (fechaInicio: string): number => {
  const fechaActual = new Date();
  const fechaInicioEstado = new Date(fechaInicio);
  const diferenciaTiempo = fechaActual.getTime() - fechaInicioEstado.getTime();
  return Math.floor(diferenciaTiempo / (1000 * 3600 * 24));
};

const fetchProfesorYContratos = async () => {
  if (!run.value) {
    errorMessage.value = "RUN no proporcionado.";
    return;
  }

  try {
    const contratosResponse = await $axios.get(`/pipelsoft/contratos-run/${run.value}`);
    console.log(contratosResponse.data);
    if (contratosResponse.data && Array.isArray(contratosResponse.data)) {
      contratos.value = contratosResponse.data;
      errorMessage.value = "";
    } else {
      errorMessage.value = "No se encontraron contratos para el RUN proporcionado.";
      contratos.value = [];
    }
  } catch (error) {
    errorMessage.value = "Hubo un error al obtener los datos.";
    console.error(error);
  }
};

onMounted(() => {
  const runFromQuery = route.query.run as string;
  if (runFromQuery) {
    run.value = runFromQuery;
    fetchProfesorYContratos();
  }
});
</script>

<style scoped>
.container {
  max-width: 100%;
  margin: 0 auto;
  padding: 20px;
  font-family: "Bebas Neue Pro", sans-serif;
  color: #394049;
}

h1 {
  font-size: 1.8rem;
  color: #EA7600;
  margin-bottom: 20px;
}

.profesor p {
  font-size: 16px;
  margin: 10px 0;
  font-family: "Helvetica Neue LT", sans-serif;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  font-family: "Helvetica Neue LT", sans-serif;
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
  color: #C8102E;
  font-weight: bold;
  margin-top: 20px;
}
</style>
