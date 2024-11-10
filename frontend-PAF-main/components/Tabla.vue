<template>
     <div class="table-container">
        <table class="w-full text-sm bg-white divide-y divide-gray-200">
        <thead>
          <tr>
            <th class="px-4 py-2 font-medium text-gray-900">Id</th>
            <th class="px-4 py-2 font-medium text-gray-900">Nombre Profesor</th>
            <th class="px-4 py-2 font-medium text-gray-900">Cupos</th>
            <th class="px-4 py-2 font-medium text-gray-900">Grupo</th>
            <th class="px-4 py-2 font-medium text-gray-900">Tipo</th>
            <th class="px-4 py-2 font-medium text-gray-900">Etapa</th>
            <th class="px-4 py-2 font-medium text-gray-900">Codigo Asignatura</th>
            <th class="px-4 py-2 font-medium text-gray-900">Opciones</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-500">
          <tr v-for="(row, index) in filteredData" :key="index">
            <td class="px-4 py-2 font-medium text-gray-900">{{ row.name }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.dob }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.role }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.salary }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.salary }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.salary }}</td>
            <td class="px-4 py-2 text-gray-700">{{ row.salary }}</td>
            <td class="px-4 py-2">
              <a href="#" class="inline-block rounded bg-indigo-600 px-4 py-2 text-xs font-medium text-white hover:bg-indigo-700">Ver PAF</a>
              <a href="#" class="inline-block rounded bg-indigo-600 px-4 py-2 text-xs font-medium text-white hover:bg-indigo-700">Ver Horarios</a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      data: Array,
      filters: Object
    },
    computed: {
      filteredData() {
        return this.data.filter(item => {
          // Filtra por disponibilidad y rango de precios (modifica según la lógica deseada)
          const inPriceRange = (!this.filters.price.from || item.salary >= this.filters.price.from) &&
                               (!this.filters.price.to || item.salary <= this.filters.price.to);
          const inAvailability = this.filters.availability.length === 0 || this.filters.availability.includes(item.availability);
          return inPriceRange && inAvailability;
        });
      }
    }
  };
  </script>

<style scoped>
.table-container {
  width: 100% /* El contenedor ocupa el ancho completo */
}

table {
  width: 100%; /* La tabla ocupa el ancho completo del contenedor */
  table-layout: auto; /* Las columnas se ajustan automáticamente */
}
</style>
  