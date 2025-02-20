<template>
  <div class="table-container">
    <table class="w-full text-sm bg-white shadow-lg rounded-lg overflow-hidden">
      <thead class="bg-primary-color text-white">
        <tr>
          <th v-if="contratos" class="px-4 py-3 text-left font-semibold">Código de la PAF</th>
          <th class="px-4 py-3 text-left font-semibold">Código de la Asignatura</th>
          <th class="px-4 py-3 text-left font-semibold">Run</th>
          <th class="px-4 py-3 text-left font-semibold">Nombre de Asignatura</th>
          <th v-if="contratos" class="px-4 py-3 text-left font-semibold">Estado de Proceso</th>
          <th v-if="show" class="px-4 py-3 text-left font-semibold">Cantidad de horas</th>
          <th v-if="showButton || show" class="px-4 py-3 text-left font-semibold">Sección</th>
          <th v-if="showButton" class="px-4 py-3 text-left font-semibold">Cupos</th>
          <th class="px-4 py-3 text-left font-semibold">Semestre de PAF</th>
          <th v-if="showButton" class="px-4 py-3 text-left font-semibold">Comentarios</th>
          <th v-if="showButtons" class="px-4 py-3 text-left font-semibold">Opciones</th>
          <th v-if="showButton" class="px-4 py-3 text-left font-semibold">Opciones</th>
          <th v-if="show" class="px-4 py-3 text-left font-semibold">Opciones</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="persona in paginatedData" :key="persona.Id" :class="persona.rowClass" class="hover:bg-gray-50 transition-colors">
          <td v-if="contratos" class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoPAF }} {{ persona.IdPaf }}</td>
          <td class="px-4 py-3 text-gray-900 font-medium">{{ persona.CodigoAsignatura }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Run }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.NombreAsignatura }} {{ persona.nombre_asignatura }}</td>
          <td v-if="contratos" class="px-4 py-3 text-gray-700">{{ persona.DesEstado }}</td>
          <td v-if="show" class="px-4 py-3 text-gray-700">{{ persona.CantidadHoras }}</td>
          <td v-if="showButton || show" class="px-4 py-3 text-gray-700">{{ persona.seccion }}</td>
          <td v-if="showButton" class="px-4 py-3 text-gray-700">{{ persona.Cupo }} {{ persona.cupo }}</td>
          <td class="px-4 py-3 text-gray-700">{{ persona.Semestre }} {{ persona.SemestrePaf }}</td>
          <td v-if="showButton" class="px-4 py-3 text-gray-700">{{ persona.Comentario }}</td>
          <td v-if="showButtons" class="px-4 py-3">
            <NuxtLink :to="`/principal/personas/paf?codigoPaf=${persona.CodigoPAF}`" class="button">Ver PAF</NuxtLink>
            <br>
            <br>
            <NuxtLink 
    :to="`/principal/personas/horario?run=${persona.Run}`"
    class="button" 
    @click.prevent="storeDetalle(1)">
    Ver Horarios
  </NuxtLink>
  <br><br>
  <NuxtLink 
    :to="`/principal/personas/horario?run=${persona.Run}`"
    class="button"
    @click.prevent="storeDetalle(2)">
    Ver Horarios detallado
  </NuxtLink>
          </td>
          <td v-if="showButton" class="px-4 py-3">
            <NuxtLink :to="`/principal/personas/paf?codigoPaf=${persona.IdPaf}`" class="button">Ver PAF</NuxtLink>
            <br>
            <br>
            <button @click="deletePAF(persona.IdPaf)" class="buttons">Eliminar</button>
          </td>
          <td v-if="show" class="px-4 py-3">
            <NuxtLink :to="`/principal/creacionContratoPAF/formularioContrato`" class="button">Formulario</NuxtLink>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Paginación -->
    <div class="pagination">
  <button 
    @click="goToPage(1)" 
    :disabled="currentPage === 1">
    «
  </button>
  <button 
    v-for="page in pageNumbers" 
    :key="page" 
    :class="{ active: currentPage === page }"
    @click="goToPage(page)">
    {{ page }}
  </button>
  <button 
    @click="goToPage(totalPages)" 
    :disabled="currentPage === totalPages">
    »
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
    },
    show: {
      type: Boolean,
      default: false
    },
    contratos: {
      type: Boolean,
      default: true
    },
  },
  data() {
    return {
      currentPage: 1,         // Página actual
      itemsPerPage: 10,       // Número de elementos por página
    };
  },
  computed: {
    pageNumbers() {
    const total = this.totalPages;
    const current = this.currentPage;
    const range = 5; // Cantidad de páginas a mostrar alrededor de la actual
    const pages = [];

    if (total <= range * 2 + 1) {
      // Mostrar todas las páginas si son pocas
      for (let i = 1; i <= total; i++) {
        pages.push(i);
      }
    } else {
      // Mostrar solo un rango de páginas alrededor de la actual
      let start = Math.max(1, current - range);
      let end = Math.min(total, current + range);

      // Asegurar que el rango siempre muestra la cantidad correcta de páginas
      if (current <= range) {
        end = Math.min(total, range * 2 + 1);
      } else if (current > total - range) {
        start = Math.max(1, total - range * 2);
      }

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (start > 1) {
        pages.unshift(1, '...');
      }

      if (end < total) {
        pages.push('...', total);
      }
    }

    return pages;
  },
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
   // Función para almacenar el valor de detalle en sessionStorage
   storeDetalle(detalleValue) {
      sessionStorage.setItem('detalle', detalleValue); // Guarda el valor de 'detalle' en sessionStorage
    },
  async deletePAF(codigoPAF) {
    try {
      const confirmDelete = confirm(`¿Estás seguro de que deseas eliminar la PAF con código ${codigoPAF}?`);
      if (!confirmDelete) return;
      
      // Realiza la solicitud DELETE
      await this.$axios.delete(`/api/paf-en-linea/historial/${codigoPAF}`);
      
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
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 5px;
}

.pagination button {
  padding: 8px 12px;
  background-color: #f0f0f0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button.active {
  background-color: #EA7600;
  color: white;
}

.pagination button:disabled {
  background-color: #EA7600;
  cursor: not-allowed;
}
</style>
