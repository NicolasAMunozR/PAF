<template>
  <!-- Botón para volver debajo de la barra superior pero encima de la información de la PAF -->
  <div class="mt-4 ml-4">
    <button @click="volver" class="volver-button">
      Volver
    </button>
  </div>
  <div class="container">
    <!-- Mostrar la lista de personas -->
    <div v-if="filteredPersonas.length > 0">
      <div v-for="persona in filteredPersonas" :key="persona.ID" class="paf-container">
        <h1 class="section-title">Información de la PAF:</h1>
        <p><strong>ID:</strong> {{ persona.ID }}</p>
        <p><strong>Run:</strong> {{ persona.run }}</p>
        <p><strong>Semestre:</strong> {{ persona.semestre }}</p>
        <p><strong>Tipo:</strong> {{ persona.tipo }}</p>
        <p><strong>Nombre Asignatura:</strong> {{ persona.nombre_asignatura }}</p>
        <p><strong>Sección:</strong> {{ persona.seccion }}</p>
        <p><strong>Cupo:</strong> {{ persona.Cupo }}</p>
        <p><strong>Bloque:</strong> {{ persona.bloque }}</p>
        <p><strong>Código Asignatura:</strong> {{ persona.codigo_asignatura }}</p>
        <p><strong>Cantidad de Horas:</strong> {{ persona.CantidadHoras }}</p>
        <!-- Botón ubicado en la parte inferior -->
        <div class="flex justify-end mt-4">
          <button v-if="!persona.Aceptada" @click="dejarListaPaf(persona.ID)" class="procesar-button">Dejar lista la PAF</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Cargando datos o no se encontraron registros para la PAF.</p>
    </div>
  </div>
</template>

  <script>
 
  </script>
  
  <style scoped>
  /* Estilos para los botones */
  .buttons-container {
    margin-bottom: 20px;
  }
  
  button {
    padding: 8px 12px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-right: 10px;
  }
  
  button:hover {
    background-color: #45a049;
  }
  
  /* Estilos para la tabla */
  .table-container {
    overflow-x: auto;
    margin-top: 20px;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
  }
  
  th, td {
    padding: 12px 15px;
    text-align: left;
    border: 2px solid #ddd;
  }
  
  th {
    background-color: #4CAF50;
    color: white;
    font-size: 16px;
    font-weight: bold;
  }
  
  td {
    font-size: 14px;
  }
  
  tr:nth-child(even) {
    background-color: #f2f2f2;
  }
  
  tr:hover {
    background-color: #f1f1f1;
  }
  
  /* Estilo para los filtros */
  .filters-container {
    display: flex;
    margin-bottom: 20px;
    gap: 15px;
  }
  
  .filter {
    display: flex;
    flex-direction: column;
  }
  
  label {
    font-weight: bold;
    margin-bottom: 5px;
  }
  
  input[type="text"] {
    padding: 8px;
    font-size: 14px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  /* Estilos para el menú desplegable */
  .menu-container {
    margin-bottom: 20px;
  }
  
  select {
    padding: 8px;
    font-size: 14px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  /* Estilos para el botón debajo de la tabla */
  .button-container {
    margin-top: 20px;
    text-align: center;
  }
  </style>
  