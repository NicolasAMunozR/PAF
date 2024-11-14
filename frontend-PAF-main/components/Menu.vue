<template>
  <div class="main-container">
    <!-- Botón para mostrar/ocultar el menú -->
    <button
      @click="toggleMenu"
      class="p-2 m-4 text-gray-600 rounded-lg hover:bg-gray-100 focus:outline-none fixed top-4 left-4 z-50"
    >
      <!-- Icono de menú hamburguesa -->
      <svg
        v-if="!isMenuOpen"
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
      <!-- Icono de cerrar -->
      <svg
        v-else
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>

    <!-- Menú desplegable como barra lateral fija -->
    <div
      :class="{'translate-x-0': isMenuOpen, '-translate-x-full': !isMenuOpen}"
      class="fixed top-0 left-0 h-full w-64 bg-white border-r border-gray-200 transition-transform duration-300 ease-in-out z-40"
    >
      <div class="px-4 py-6">
        <span class="grid h-10 w-32 place-content-center rounded-lg bg-gray-100 text-xs text-gray-600">
          Logo
        </span>

        <ul class="mt-6 space-y-1">
          <li>
            <a href="/" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Home</a>
          </li>
          <li>
            <a href="/personas" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Listado de Personas</a>
          </li>
          <li>
            <a href="/paf" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">PAF</a>
          </li>
          <li>
            <a href="/historyPAF" class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700">Historial de PAF</a>
          </li>
          <li>
            <details class="group [&_summary::-webkit-details-marker]:hidden">
              <summary class="flex cursor-pointer items-center justify-between rounded-lg px-4 py-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700">
                <span class="text-sm font-medium">Teams</span>
                <span class="shrink-0 transition duration-300 group-open:-rotate-180">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a 1 1 0 01-1.414 0l-4-4a 1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </span>
              </summary>
              <ul class="mt-2 space-y-1 px-4">
                <li>
                  <a href="#" class="block rounded-lg px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700">Banned Users</a>
                </li>
                <li>
                  <a href="#" class="block rounded-lg px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700">Calendar</a>
                </li>
              </ul>
            </details>
          </li>
          <!-- Otros elementos del menú -->
        </ul>
      </div>
    </div>

    <!-- Contenido principal -->
    <div :class="{'ml-64': isMenuOpen, 'ml-0': !isMenuOpen}" class="content transition-all duration-300 ease-in-out">
      <!-- Contenido de la página -->
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isMenuOpen: false,
    };
  },
  methods: {
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
    },
  },
};
</script>

<style scoped>
.main-container {
  margin-top: 20px; /* Espacio superior para evitar que el botón quede pegado al borde */
}

.content {
  padding: 20px;
}
</style>
