<template>
  <div class="container">
    <Filtros />
    <ListaPersonas />
  </div>
</template>

<script setup lang="ts">
import ListaPersonas from '../../components/ListaPersonas.vue'
import Filtros from '../../components/Filtros.vue'

</script>

<style scoped>
  .container {
    display: grid;
    grid-template-columns: 1fr 3fr; /* 1/4 para filtros, 3/4 para tabla */
    gap: 1rem; /* Espacio entre columnas */
  }
  </style>
  
  