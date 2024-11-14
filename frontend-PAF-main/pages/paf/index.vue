<template>
    <h1 class="text-3xl text-center">PAF</h1>
    <div class="container">
      <!-- Columna de Filtros -->
      <div class="filters">
        <Filtros @filter-changed="applyFilters" />
      </div>
  
      <!-- Columna de Tabla -->
      <div class="table">
        <Tabla :data="tableData" :filters="currentFilters" />
      </div>
    </div>
  </template>
  
  <script>
  import Filtros from '../components/Filtros.vue';
  import Tabla from '../components/Tabla.vue';
  
  export default {
    components: {
      Filtros,
      Tabla
    },
    data() {
      return {
        tableData: [
          { name: 'John Doe', dob: '24/05/1995', role: 'Web Developer', salary: 120000, availability: 'inStock' },
          { name: 'Jane Doe', dob: '04/11/1980', role: 'Web Designer', salary: 100000, availability: 'preOrder' },
          { name: 'Gary Barlow', dob: '24/05/1995', role: 'Singer', salary: 20000, availability: 'outOfStock' }
        ],
        currentFilters: { availability: [], price: { from: '', to: '' } }
      };
    },
    methods: {
      applyFilters(filters) {
        this.currentFilters = filters;
      }
    }
  };
  </script>
<style scoped>
  .container {
    display: grid;
    grid-template-columns: 1fr 3fr; /* 1/4 para filtros, 3/4 para tabla */
    gap: 1rem; /* Espacio entre columnas */
  }
  
  .filters {
    background-color: #f9f9f9; /* Color de fondo para diferenciar */
    padding: 1rem;
    border-radius: 8px;
  }
  
  .table {
    overflow-x: auto; /* Habilitar scroll si la tabla es muy grande */
  }
  </style>
  
  