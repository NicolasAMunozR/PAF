<template>
  <div class="main-container">
    <!-- Barra superior fija -->
    <div class="top-bar fixed top-0 left-0 w-full bg-gray-800 text-white flex items-center justify-between p-4 z-50">
      <!-- Botón para mostrar/ocultar el menú -->
      <button
        @click="toggleMenu"
        class="p-2 text-gray-600 rounded-lg hover:bg-gray-100 focus:outline-none"
      >
        <!-- Icono de menú hamburguesa -->
        <svg
          v-if="!isMenuOpen"
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          color="white"
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

      <div class="header-container flex-1 text-center">
        <h1 class="header-title">PAF - Sistema de Gestión</h1>
      </div>

      <!-- Botón de cerrar sesión -->
      <button
        @click="logout"
        class="text-white bg-red-600 p-2 rounded-lg hover:bg-red-700 focus:outline-none"
      >
        Cerrar sesión
      </button>
    </div>

    <!-- Menú desplegable como barra lateral fija -->
    <div
      :class="{'translate-x-0': isMenuOpen, '-translate-x-full': !isMenuOpen}"
      class="fixed top-16 left-0 h-full w-64 bg-white border-r border-gray-200 transition-transform duration-300 ease-in-out z-40"
    >
      <div class="px-4 py-6">
        <ul class="mt-6 space-y-1">
          <li v-for="link in filteredMenu" :key="link.path">
            <a
              :href="link.path"
              class="block rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700"
            >
              {{ link.label }}
            </a>
          </li>
        </ul>
      </div>
    </div>

    <!-- Contenido principal -->
    <div :class="{'ml-64': isMenuOpen, 'ml-0': !isMenuOpen}" class="content transition-all duration-300 ease-in-out">
      <!-- Contenido de la página -->
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isMenuOpen: false,
      menuItems: [
        {
          path: "/personas",
          label: "Listado de Personas",
          pages: ["/personas", "/historyPAF", "/horario", "/paf"],
        },
        {
          path: "/historyPAF",
          label: "Historial de PAF",
          pages: ["/personas", "/historyPAF", "/horario", "/paf"],
        },
        {
          path: "/seguimientoPAF",
          label: "Seguimiento de la PAF",
          pages: ["/seguimientoPAF", "/estadisticaPAF"],
        },
        {
          path: "/profesorPAF",
          label: "Gestión de Profesor",
          pages: [],
        },
        {
          path: "/estadisticaPAF",
          label: "Estadísticas de PAF",
          pages: ["/seguimientoPAF", "/estadisticaPAF"],
        },
      ],
    };
  },
  computed: {
    filteredMenu() {
      const currentPage = this.$route.path; // Ruta actual
      return this.menuItems.filter((item) =>
        item.pages.includes(currentPage)
      );
    },
  },
  methods: {
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
    },
    logout() {
      // Redirigir al usuario a la página principal ("/") para cerrar sesión
      this.$router.push('/');
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
</style>
