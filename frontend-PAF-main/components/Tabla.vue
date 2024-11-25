<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-gray-800 text-white">
        <tr>
          <th class="px-4 py-3 text-left font-semibold">Codigo de la PAF</th>
          <th class="px-4 py-3 text-left font-semibold">Codigo de la Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Run</th>
          <th class="px-4 py-3 text-left font-semibold">Nombre de asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Grupo</th>
          <th class="px-4 py-3 text-left font-semibold">Cupos</th>
          <th v-if="showButtons" class="px-4 py-3 text-left font-semibold">Opciones</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="persona in data" :key="persona.CodigoPAF" class="hover:bg-gray-50 transition-colors">
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoPAF }}</td>
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.NombreAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700"> </td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Cupo }} {{ persona.cupo }} </td>
          <td v-if="showButtons" class="px-4 py-3">
            <a :href="`/paf?codigoPaf=${persona.CodigoPAF}`" class="button">Ver PAF</a>
            <a :href="`/horario?run=${persona.Run}`" class="button">Ver Horarios</a>
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
    },
    showButtons: {
      type: Boolean,
      default: true // Por defecto, se muestran los botones
    }
  }
}
</script>



<style scoped>
.table-container {
  width: 100%;
  padding: 20px;
  overflow-x: auto;
  background-color: #f9fafb;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #ffffff;
  background-color: #4a5568;
}

tbody td {
  padding: 12px;
  font-size: 0.875rem;
  color: #4a5568;
}

.button {
  display: inline-block;
  padding: 8px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  color: #ffffff;
  background-color: #4f46e5;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.button:hover {
  background-color: #4338ca;
}

.hover\:bg-gray-50:hover {
  background-color: #f1f5f9;
}
</style>
