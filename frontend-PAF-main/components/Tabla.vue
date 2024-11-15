<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-gray-800 text-white">
        <tr>
          <th class="px-4 py-3 text-left font-semibold">Codigo de la PAF</th>
          <th class="px-4 py-3 text-left font-semibold">Codigo de la Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Run</th>
          <th class="px-4 py-3 text-left font-semibold">Nombre</th>
          <th class="px-4 py-3 text-left font-semibold">Apellido</th>
          <th class="px-4 py-3 text-left font-semibold">Grupo</th>
          <th class="px-4 py-3 text-left font-semibold">Cupos</th>
          <th class="px-4 py-3 text-left font-semibold"> </th>
          <th class="px-4 py-3 text-left font-semibold">Opciones</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="persona in data" :key="persona.id" class="hover:bg-gray-50 transition-colors">
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoPAF }}</td>
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Nombres }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.PrimerApellido }}</td>
          <td class="px-4 py-3 text-gray-700"> </td>
          <td class="px-4 py-3 text-gray-700"> </td>
          <td class="px-4 py-3 text-gray-700"> </td>
          <td class="px-4 py-3">
            <a href="#" class="button">Ver PAF</a>
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
