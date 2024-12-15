<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-primary-color text-white">
        <tr>
          <th class="px-4 py-3 text-left font-semibold">Código de la PAF</th>
          <th class="px-4 py-3 text-left font-semibold">Código de la Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Run</th>
          <th class="px-4 py-3 text-left font-semibold">Nombre de Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Estado de Proceso</th>
          <th class="px-4 py-3 text-left font-semibold">Sección</th>
          <th class="px-4 py-3 text-left font-semibold">Cupos</th>
          <th class="px-4 py-3 text-left font-semibold">Semestre de PAF</th>
          <th v-if="showButtons" class="px-4 py-3 text-left font-semibold">Opciones</th>
          <th v-if="showButton" class="px-4 py-3 text-left font-semibold">Opciones</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="persona in paginatedData" :key="persona.Id" :class="persona.rowClass" class="hover:bg-gray-50 transition-colors">
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoPAF }} {{ persona.IdPaf }}</td>
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.NombreAsignatura }} {{ persona.nombre_asignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.DesEstado }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.seccion }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Cupo }} {{ persona.cupo }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Semestre }} {{ persona.SemestrePaf }}</td>
          <td v-if="showButtons" class="px-4 py-3">
            <NuxtLink :to="`/personas/paf?codigoPaf=${persona.CodigoPAF}`" class="button">Ver PAF</NuxtLink>
            <br>
            <br>
            <NuxtLink :to="`/personas/horario?run=${persona.Run}`" class="button">Ver Horarios</NuxtLink>
          </td>
          <td v-if="showButton" class="px-4 py-3">
            <NuxtLink :to="`/personas/paf?codigoPaf=${persona.IdPaf}`" class="button">Ver PAF</NuxtLink>
            <br>
            <br>
            <button @click="deletePAF(persona.IdPaf)" class="buttons">Eliminar</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Paginación -->
    <div class="pagination">
      <button 
        @click="goToPage(currentPage - 1)" 
        :disabled="currentPage === 1">
        Anterior
      </button>
      <span>Página {{ currentPage }} de {{ totalPages }}</span>
      <button 
        @click="goToPage(currentPage + 1)" 
        :disabled="currentPage === totalPages">
        Siguiente
      </button>
    </div>
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
      default: true
    },
    showButton: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      currentPage: 1,         // Página actual
      itemsPerPage: 10,       // Número de elementos por página
    };
  },
  computed: {
    // Paginación: calculamos los elementos que se mostrarán según la página
    paginatedData() {
      const startIndex = (this.currentPage - 1) * this.itemsPerPage;
      const endIndex = startIndex + this.itemsPerPage;
      return this.data.slice(startIndex, endIndex);
    },
    // Número total de páginas
    totalPages() {
      return Math.ceil(this.data.length / this.itemsPerPage);
    }
  },
  methods: {
  goToPage(page) {
    if (page < 1 || page > this.totalPages) return;
    this.currentPage = page;
  },
  async deletePAF(codigoPAF) {
    try {
      const confirmDelete = confirm(`¿Estás seguro de que deseas eliminar la PAF con código ${codigoPAF}?`);
      if (!confirmDelete) return;
      
      // Realiza la solicitud DELETE
      await this.$axios.delete(`/historial/${codigoPAF}`);
      
      // Muestra una notificación o mensaje de éxito
      alert(`PAF con código ${codigoPAF} eliminada con éxito.`);
      
    } catch (error) {
      // Maneja errores
      console.error("Error al eliminar la PAF:", error);
      alert("Ocurrió un error al intentar eliminar la PAF.");
    }
  }
}
}
</script>


<style scoped>
/* Colores institucionales */
:root {
  --primary-color: #EA7600; /* Color principal USACH */
  --secondary-color: #394049; /* Color secundario USACH */
  --accent-color: #C8102E; /* Complementario */
  --background-color: #1d558d; /* Fondo neutro */
  --button-background-color: #4CAF50; /* Color de fondo de los botones */
  --button-hover-color: #388E3C; /* Color de hover en los botones */
}

/* Fila modificada y eliminada */
.modified-row {
  background-color: yellow; /* Amarillo suave */
}

.deleted-row {
  background-color: red; /* Rojo suave */
}

/* Contenedor de la tabla */
.table-container {
  width: 100%;
  padding: 20px;
  overflow-x: auto;
  background-color: var(--background-color);
}

/* Estilos para la tabla */
table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #030101;
  background-color: var(--primary-color);
}

tbody td {
  padding: 12px;
  font-size: 0.875rem;
  color: var(--secondary-color);
}

/* Estilos para los botones */
.button {
  display: inline-block;
  padding: 8px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  color: #090000;
  background-color: #4CAF50;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.button:hover {
  background-color: var(--button-hover-color);
}

/* Estilos para los botones */
.buttons {
  display: inline-block;
  padding: 8px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  color: #090000;
  background-color: #e00606;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.buttons:hover {
  background-color: var(--button-hover-color);
}
.hover\:bg-gray-50:hover {
  background-color: #f1f5f9;
}
/* Paginación */
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination button {
  padding: 8px 16px;
  background-color: #EA7600;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  background-color: #f0f0f0;
  cursor: not-allowed;
}

.pagination span {
  font-size: 1rem;
}
</style>
