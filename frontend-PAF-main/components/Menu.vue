<template>
  <div class="main-container">
    <!-- Barra superior fija -->
    <div class="top-bar">
      <!-- Botón para mostrar/ocultar el menú -->
      <button
        @click="toggleMenu"
        class="menu-button"
      >
        <!-- Icono de menú hamburguesa -->
        <svg
          v-if="!isMenuOpen"
          xmlns="http://www.w3.org/2000/svg"
          class="icon"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <!-- Icono de cerrar -->
        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          class="icon"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div class="header-container">
        <h1 class="header-title">PAF - Sistema de Gestión</h1>
      </div>

      <!-- Botón de cerrar sesión -->
      <button
        @click="logout"
        class="logout-button"
      >
        Cerrar sesión
      </button>
    </div>

    <!-- Menú desplegable como barra lateral fija -->
    <div
      :class="{'translate-x-0': isMenuOpen, '-translate-x-full': !isMenuOpen}"
      class="side-menu"
    >
      <div class="menu-content">
        <ul class="menu-list">
          <li v-for="link in filteredMenu" :key="link.path">
            <a :href="link.path" class="menu-link">
              {{ link.label }}
            </a>
          </li>
        </ul>
      </div>
    </div>

    <!-- Contenido principal -->
    <div :class="{'ml-64': isMenuOpen, 'ml-0': !isMenuOpen}" class="content">
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
        { path: "/personas", label: "Listado de Personas", pages: ["/personas", "/historyPAF", "/horario", "/paf"] },
        { path: "/historyPAF", label: "Historial de PAF", pages: ["/personas", "/historyPAF", "/horario", "/paf"] },
        { path: "/seguimientoPAF", label: "Seguimiento de la PAF", pages: ["/seguimientoPAF", "/estadisticaPAF"] },
        { path: "/estadisticaPAF", label: "Estadísticas de PAF", pages: ["/seguimientoPAF", "/estadisticaPAF"] },
        { path: "/unidadMayorPAF", label: "Gestión de Unidad Mayor", pages: ["/unidadMayorPAF", "/estadisticaUnidadMayorPAF"] },
        { path: "/estadisticaUnidadMayorPAF", label: "Estadísticas de Unidad Mayor", pages: ["/unidadMAyorPAF", "/estadisticaUnidadMayorPAF"] },
      ],
    };
  },
  computed: {
    filteredMenu() {
      const currentPage = this.$route.path;
      return this.menuItems.filter((item) => item.pages.includes(currentPage));
    },
  },
  methods: {
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
    },
    logout() {
      this.$router.push('/');
    },
  },
};
</script>

<style scoped>
/* Contenedor principal */
.main-container {
  margin-top: 20px;
}

/* Barra superior */
.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background-color: #EA7600;
  color: white;
  padding: 1rem;
  z-index: 50;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-container {
  flex: 1;
  display: flex;
  justify-content: center;
}

.header-title {
  font-family: "Bebas Neue Pro", sans-serif;
  font-size: 1.5rem;
  font-weight: bold;
}

/* Botones */
.menu-button,
.logout-button {
  background-color: #394049;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.menu-button:hover,
.logout-button:hover {
  background-color: #C8102E;
}

/* Menú lateral */
.side-menu {
  position: fixed;
  top: 4rem;
  left: 0;
  height: 100%;
  width: 16rem;
  background-color: #394049;
  border-right: 2px solid #EA7600;
  transition: transform 0.3s;
}

.menu-content {
  padding: 1rem;
}

.menu-list {
  list-style: none;
  padding: 0;
}

.menu-link {
  display: block;
  padding: 10px 15px;
  font-family: "Helvetica Neue LT", sans-serif;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  background-color: #00A499;
  transition: background-color 0.3s;
}

.menu-link:hover {
  background-color: #C8102E;
}

/* Contenido */
.content {
  padding: 20px;
  margin-top: 4rem;
}
</style>
