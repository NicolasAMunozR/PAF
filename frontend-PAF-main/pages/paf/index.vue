<template>
  <div class="main-container">
    <!-- Barra superior fija -->
    <div class="top-bar fixed top-0 left-0 w-full bg-gray-800 text-white flex items-center justify-between p-4 z-50">
      <h1 class="header-title">PAF - Sistema de Gestión</h1>
      <button @click="toggleMenu" class="p-2 text-gray-600 rounded-lg hover:bg-gray-100 focus:outline-none">
        <svg v-if="!isMenuOpen" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" color="white">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Menú desplegable como barra lateral fija -->
    <div :class="{'translate-x-0': isMenuOpen, '-translate-x-full': !isMenuOpen}" class="fixed top-16 left-0 h-full w-64 bg-white border-r border-gray-200 transition-transform duration-300 ease-in-out z-40">
      <div class="px-4 py-6">
        <ul class="mt-6 space-y-1">
          <li><a href="/" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Home</a></li>
          <li><a href="/personas" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Listado de Personas</a></li>
          <li><a href="/paf" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">PAF</a></li>
          <li><a href="/historyPAF" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Historial de PAF</a></li>
        </ul>
      </div>
    </div>

    <!-- Contenido principal -->
    <div :class="{'ml-64': isMenuOpen, 'ml-0': !isMenuOpen}" class="content transition-all duration-300 ease-in-out">
      <!-- Contenedor de información de la PAF -->
      <div class="paf-info-container">
        <h2 class="text-2xl font-bold mb-4">Información de la PAF</h2>
        <p class="mb-2">Código PAF: {{ paf.codigo_paf }}</p>
        <p class="mb-2">Run: {{ paf.run }}</p>
        <!-- Otros detalles de la PAF -->

        <!-- Botones de acción -->
        <div class="action-buttons mt-4">
          <button @click="goBack" class="bg-gray-500 text-white py-2 px-4 rounded mr-4">Volver</button>
          <button @click="markPafAsReady" class="bg-blue-500 text-white py-2 px-4 rounded">Dejar lista la PAF</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isMenuOpen: false,
      paf: {
        codigo_paf: 'PAF001',
        run: '11111111-1', // Este es el RUN que necesitas
      },
    };
  },
  methods: {
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
    },
    goBack() {
      this.$router.go(-1); // Vuelve a la página anterior
    },
    async markPafAsReady() {
      try {
        // Llamada al backend para crear el historial (POST)
        const response = await this.$axios.post('/historial', {
          Run: this.paf.run,
          CodigoPAF: this.paf.codigo_paf,
          FechaAceptacionPaf: new Date().toISOString(), // Fecha actual en formato ISO
        });

        console.log('Historial creado:', response.data);
        // Luego de completar el POST, puedes hacer algo, como redirigir o mostrar un mensaje
      } catch (error) {
        console.error('Error al crear el historial:', error);
      }
    },
  },
};
</script>

<style scoped>
.main-container {
  margin-top: 20px; /* Espacio superior para evitar que el contenido quede pegado al borde */
}

.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background-color: #1f2937; /* Fondo oscuro para la barra superior */
  padding: 1rem;
  z-index: 50;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-container {
  flex: 1; /* Permite que ocupe todo el espacio disponible */
  display: flex;
  justify-content: center; /* Centra el título horizontalmente */
}

.header-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
}

.content {
  padding: 20px;
  margin-top: 4rem; /* Espacio para que el contenido no quede pegado a la barra superior */
}

.paf-info-container {
  background-color: #f9fafb;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.action-buttons button {
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.action-buttons button:hover {
  opacity: 0.8;
}
</style>
