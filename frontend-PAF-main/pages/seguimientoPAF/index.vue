<template>
    <div class="container">
      <h1>Buscar Contratos por RUN</h1>
      <form @submit.prevent="fetchContratos">
        <label for="run">Ingrese el RUN del Profesor:</label>
        <input v-model="run" id="run" type="text" placeholder="Ej. 12345678-9" required />
        <button type="submit">Buscar</button>
      </form>
  
      <div v-if="contratos.length > 0" class="contratos">
        <h2>Contratos Relacionados</h2>
        <table>
          <thead>
            <tr>
              <th>Código PAF</th>
              <th>Jefatura</th>
              <th>Nombre de Asignatura</th>
              <th>Estado del Proceso</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(contrato, index) in contratos" :key="index">
              <td>{{ contrato.pipelsoft_data.CodigoPAF }}</td>
              <td>{{ contrato.pipelsoft_data.Jerarquia }}</td>
              <td>{{ contrato.pipelsoft_data.NombreAsignatura }}</td>
              <td>{{ estadoProceso(contrato.pipelsoft_data.EstadoProceso) }}</td>
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
  import { ref } from 'vue'
  const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }
  
  const run = ref('')
  const contratos = ref<any[]>([])
  const errorMessage = ref('')
  
  // Función para mapear los estados de proceso a nombres legibles
  const estadoProceso = (estado: number): string => {
    switch (estado) {
      case 1:
        return "Activo"
      case 2:
        return "Inactivo"
      case 3:
        return "Finalizado"
      default:
        return "Desconocido"
    }
  }
  
  const fetchContratos = async () => {
    if (!run.value) {
      errorMessage.value = 'Por favor, ingrese un RUN válido.'
      return
    }
  
    try {
      // Hacer la solicitud al backend con el RUN
      const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`)
      console.log(response)
  
      // Si la respuesta contiene data, se asigna a contratos
      if (response.data && Array.isArray(response.data)) {
        contratos.value = response.data
      } else {
        errorMessage.value = 'No se encontraron contratos para el RUN ingresado.'
      }
    } catch (error) {
      errorMessage.value = 'Hubo un error al obtener los datos.'
      console.error(error)
    }
  }
  </script>
  
  <style scoped>
  .container {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
  }
  
  form {
    margin-bottom: 20px;
  }
  
  input {
    padding: 8px;
    margin-right: 10px;
  }
  
  button {
    padding: 8px 16px;
    background-color: #4CAF50;
    color: white;
    border: none;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #45a049;
  }
  
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
    background-color: #f2f2f2;
  }
  
  .error {
    color: red;
    font-weight: bold;
  }
  </style>
  