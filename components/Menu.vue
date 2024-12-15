<template>
  <div class="main-container">
    <!-- Barra superior fija -->
    <div class="top-bar">
      <button
        @click="toggleMenu"
        v-if="$route.path !== '/profesorPAF'"
        class="menu-button"
      >
        <svg
          v-if="!isMenuOpen"
          xmlns="http://www.w3.org/2000/svg"
          class="icon"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
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

      <button @click="logout" class="logout-button">Cerrar sesión</button>
    </div>

    <!-- Superposición detrás del menú -->
    <div
      v-if="isMenuOpen"
      class="overlay"
      @click="toggleMenu"
    ></div>

    <!-- Menú lateral -->
    <div
      v-if="$route.path !== '/profesorPAF'"
      :class="{
        'translate-x-0': isMenuOpen,
        '-translate-x-full': !isMenuOpen
      }"
      class="side-menu"
    >
      <div class="menu-content">
        <ul class="menu-list">
          <li v-for="link in filteredMenu" :key="link.path">
            <a :href="link.path" class="menu-link">{{ link.label }}</a>
          </li>
        </ul>
      </div>
    </div>

    <!-- Contenido principal -->
    <div
      :class="{
        'ml-64': isMenuOpen && $route.path !== '/principal/profesorPAF',
        'ml-0': !isMenuOpen || $route.path === '/principal/profesorPAF'
      }"
      class="content"
    >
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
        { path: "/principal/personas", label: "Personas", pages: ["/principal/personas", "/principal/historyPAF", "/personas/horario", "/personas/paf"] },
        { path: "/principal/historyPAF", label: "Historial de PAF", pages: ["/principal/personas", "/principal/historyPAF", "/personas/horario", "/personas/paf"] },
        { path: "/seguimientoPAF", label: "Seguimiento de la PAF", pages: ["/principal/seguimientoPAF", "/principal/estadisticaPAF"] },
        { path: "/principal/estadisticaPAF", label: "Estadísticas de PAF", pages: ["/principal/seguimientoPAF", "/principal/estadisticaPAF"] },
        { path: "/principal/unidadMayorPAF", label: "Gestión de Unidad Mayor", pages: ["/principal/unidadMayorPAF", "/principal/estadisticaUnidadPAF"] },
        { path: "/principal/estadisticaUnidadPAF", label: "Estadísticas de Unidad", pages: ["/principal/unidadMayorPAF", "/principal/estadisticaUnidadPAF"] },
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
      if (this.$route.path === "/principal/profesorPAF") {
        return; // Evitar que el menú se abra/cierre en ProfesorPAF
      }
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
  display: flex;
  flex-direction: column;
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
  height: 4rem;
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

/* Superposición detrás del menú */
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 900; /* Debajo del menú lateral */
}

/* Menú lateral */
.side-menu {
  position: fixed;
  top: 0;
  left: 0;
  height: 100%;
  width: 16rem;
  background-color: #394049;
  border-right: 2px solid #EA7600;
  transition: transform 0.3s;
  z-index: 1000; /* Siempre por encima de todo */
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
  margin-top: 4rem; /* Espaciado dinámico igual a la altura de la barra superior */
  flex: 1;
}
</style>