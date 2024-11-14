<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white divide-y divide-gray-200">
      <thead>
        <tr>
          <th class="px-4 py-2 font-medium text-gray-900">ID</th>
          <th class="px-4 py-2 font-medium text-gray-900">Run</th>
          <th class="px-4 py-2 font-medium text-gray-900">Nombres</th>
          <th class="px-4 py-2 font-medium text-gray-900">Correo</th>
          <th class="px-4 py-2 font-medium text-gray-900">Grupo</th>
          <th class="px-4 py-2 font-medium text-gray-900">Cupos</th>
          <th class="px-4 py-2 font-medium text-gray-900">  </th>
          <th class="px-4 py-2 font-medium text-gray-900">Opciones</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-500">
        <tr v-for="persona in data" :key="persona.id">
          <td class="px-4 py-2 font-medium text-gray-900">{{ persona.ID }}</td>
          <td class="px-4 py-2 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-2 text-gray-700">{{ persona.Nombres }}</td>
          <td class="px-4 py-2 text-gray-700">{{ persona.Correo }}</td>
          <td class="px-4 py-2 text-gray-700"> </td>
          <td class="px-4 py-2 text-gray-700"> </td>
          <td class="px-4 py-2 text-gray-700"> </td>
          <td class="px-4 py-2">
            <a href="#" class="inline-block rounded bg-indigo-600 px-4 py-2 text-xs font-medium text-white hover:bg-indigo-700 mr-2">Ver PAF</a>
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
    data: {
      type: Array,
      required: true
    }
  }
}
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
  