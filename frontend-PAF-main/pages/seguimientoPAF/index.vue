<template>
  <div class="container">
    <h1>Buscar Contratos por RUN</h1>
    <form @submit.prevent="fetchContratos">
      <label for="run">Ingrese el RUN del Profesor:</label>
      <input v-model="run" id="run" type="text" placeholder="Ej. 12345678-9" required />
      <button type="submit">Buscar</button>
    </form>
    <h1>Buscar Contratos por Unidad Contratante</h1>
    <form @submit.prevent="fetchContratosUnidadContratante">
      <label for="unidad">Ingrese la Unidad Contratante:</label>
      <input v-model="Unidad" id="unidad" type="text" placeholder="" required />
      <button type="submit">Buscar</button>
    </form>
    <h1>Buscar Contratos por Unidad Mayor</h1>
    <form @submit.prevent="fetchContratosUnidadMayor">
      <label for="unidadMayor">Ingrese la Unidad:</label>
      <input v-model="NombreUnidadMayor" id="unidadMayor" type="text" placeholder="" required />
      <button type="submit">Buscar</button>
    </form>
    <div v-if="contratos.length > 0" class="contratos">
      <h2>Contratos Relacionados a: {{ contratos[0].PipelsoftData.Nombres }} {{ contratos[0].PipelsoftData.PrimerApp }} {{ contratos[0].PipelsoftData.SegundoApp }}</h2>
      <table>
        <thead>
          <tr>
            <th>Código PAF</th>
            <th>Jefatura</th>
            <th>Nombre de Asignatura</th>
            <th>Estado del Proceso</th>
            <th>Fecha de la última Actualización de Estado</th>
            <th>Historial de Estados</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(contrato, index) in contratos" :key="index">
            <td>{{ contrato.PipelsoftData.IdPaf }}</td>
            <td>{{ contrato.PipelsoftData.Jerarquia }}</td>
            <td>{{ contrato.PipelsoftData.NombreAsignatura }}</td>
            <td>{{ estadoProceso(contrato.PipelsoftData.CodEstado) }}</td>
            <td>
              {{ new Date(contrato.PipelsoftData.FechaUltimaModificacionProceso).toLocaleTimeString() }} 
              {{ new Date(contrato.PipelsoftData.FechaUltimaModificacionProceso).toLocaleDateString() }}
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
import { ref } from 'vue';
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

const run = ref('');
const Unidad = ref('');
const NombreUnidadMayor = ref('');
const contratos = ref<any[]>([]);
const errorMessage = ref('');

// Función para mapear los estados de proceso a nombres legibles
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

// Función para calcular el tiempo en cada estado en días
const calcularTiempoEnEstado = (fechaInicio: string): number => {
  const fechaActual = new Date();
  const fechaInicioEstado = new Date(fechaInicio);
  const diferenciaTiempo = fechaActual.getTime() - fechaInicioEstado.getTime();
  return Math.floor(diferenciaTiempo / (1000 * 3600 * 24)); // Devuelve los días
};

const fetchContratos = async () => {
  if (!run.value) {
    errorMessage.value = 'Por favor, ingrese un RUN válido.';
    return;
  }

  try {
    const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`);
    if (response.data && Array.isArray(response.data)) {
      contratos.value = response.data;
    } else {
      errorMessage.value = 'No se encontraron contratos para el RUN ingresado.';
    }
  } catch (error) {
    errorMessage.value = 'Hubo un error al obtener los datos.';
    console.error(error);
  }
};

const fetchContratosUnidadContratante = async () => {
  if (!Unidad.value) {
    errorMessage.value = 'Por favor, ingrese una unidad contratante válida.';
    return;
  }

  try {
    const response = await $axios.get(`/pipelsoft/contratos-nombreUnidadContratante/${Unidad.value}`);
    if (response.data && Array.isArray(response.data)) {
      contratos.value = response.data;
    } else {
      errorMessage.value = 'No se encontraron contratos para la unidad contratante ingresada.';
    }
  } catch (error) {
    errorMessage.value = 'Hubo un error al obtener los datos.';
    console.error(error);
  }
};

const fetchContratosUnidadMayor = async () => {
  if (!NombreUnidadMayor.value) {
    errorMessage.value = 'Por favor, ingrese una unidad mayor válida.';
    return;
  }

  try {
    const response = await $axios.get(`/pipelsoft/contratos-nombreUnidadMayor/${NombreUnidadMayor.value}`);
    if (response.data && Array.isArray(response.data)) {
      contratos.value = response.data;
    } else {
      errorMessage.value = 'No se encontraron contratos para la unidad mayor ingresada.';
    }
  } catch (error) {
    errorMessage.value = 'Hubo un error al obtener los datos.';
    console.error(error);
  }
};
</script>

<style scoped>
/* Contenedor principal */
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: "Bebas Neue Pro", sans-serif;
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
</style>
