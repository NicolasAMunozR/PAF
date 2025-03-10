<template>
  <div class="main-container">
    <!-- Barra superior fija -->
    <div class="top-bar">
      <button
        @click="toggleMenu"
        v-if="$route.path !== '/profesorPAF'"
        class="menu-button"
      >
        <svg v-if="!isMenuOpen" class="hamburger-icon">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg v-else class="hamburger-icon">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
      <div class="header-container">
        <h1 class="header-title">PAF - Sistema de Gestión</h1>
      </div>
      <button @click="logout" class="logout-button">Cerrar sesión</button>
    </div>

    <!-- Superposición detrás del menú -->
    <div v-if="isMenuOpen" class="overlay" @click="toggleMenu"></div>

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
          <li
            v-for="link in filteredMenu"
            :key="link.path"
            
            @mouseleave="hideDropdown(link.path)"
            class="menu-item"
          >
            <a @click.prevent="navigateTo(link.path)" class="menu-link">{{ link.label }}</a>
            <ul v-if="dropdownVisible[link.path]" class="dropdown-menu">
              <li v-for="page in link.pages" :key="page">
                <a @click.prevent="navigateTo(page)" class="dropdown-link">{{ page }}</a>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>

    <!-- Contenido paf-en-linea -->
    <div
      :class="{
        'ml-64': isMenuOpen && $route.path !== '/paf-en-linea/profesorPAF',
        'ml-0': !isMenuOpen || $route.path === '/paf-en-linea/profesorPAF'
      }"
      class="content"
    >
      <slot></slot>
    </div>
  </div>
</template>

<script>
import { reactive } from 'vue';

export default {
  data() {
    return {
      isMenuOpen: false,
      dropdownVisible: reactive({}),
      menuItems: [
        { path: "/paf-en-linea/personas", label: "Personas", pages: ["/paf-en-linea/personas", "/paf-en-linea/historyPAF", "/paf-en-linea/personas/horario", "/paf-en-linea/personas/paf"] },
        { path: "/paf-en-linea/historyPAF", label: "Historial de PAF", pages: ["/paf-en-linea/personas", "/paf-en-linea/historyPAF", "/paf-en-linea/personas/horario", "/paf-en-linea/personas/paf"] },
        { path: "/paf-en-linea/seguimientoPAF", label: "Seguimiento de la PAF", pages: ["/paf-en-linea/seguimientoPAF", "/paf-en-linea/estadisticaPAF", "/paf-en-linea/creacionContratoPAF"] },
        { path: "/paf-en-linea/estadisticaPAF", label: "Estadísticas de PAF", pages: ["/paf-en-linea/seguimientoPAF", "/paf-en-linea/estadisticaPAF", "/paf-en-linea/creacionContratoPAF"] },
        { path: "/paf-en-linea/unidadMayorPAF", label: "Gestión de Unidad Mayor", pages: ["/paf-en-linea/unidadMayorPAF", "/paf-en-linea/estadisticaUnidadPAF"] },
        { path: "/paf-en-linea/estadisticaUnidadPAF", label: "Estadísticas de Unidad", pages: ["/paf-en-linea/unidadMayorPAF", "/paf-en-linea/estadisticaUnidadPAF"] },
        { path: "/paf-en-linea/creacionContratoPAF", label: "creacion de contratos", pages: ["/paf-en-linea/seguimientoPAF", "/paf-en-linea/estadisticaPAF","/paf-en-linea/creacionContratoPAF" ] },
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
      if (this.$route.path === "/paf-en-linea/profesorPAF") {
        return; // Evitar que el menú se abra/cierre en ProfesorPAF
      }
      this.isMenuOpen = !this.isMenuOpen;
    },
    navigateTo(path) {
      this.$router.push(path);
      this.isMenuOpen = false; // Cierra el menú después de la navegación
    },
    logout() {
      localStorage.clear();
      sessionStorage.clear();
      this.$router.push('/');
    },
    showDropdown(path) {
      this.dropdownVisible[path] = true;
    },
    hideDropdown(path) {
      this.dropdownVisible[path] = false;
    },
  },
};
</script>

<style scoped>
.menu-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 10px;
}

.hamburger-icon {
  display: inline-block;
  width: 24px;
  height: 24px;
  position: relative;
  cursor: pointer;
}

.hamburger-icon path {
  stroke: white;
  stroke-width: 3;
}

.hamburger-icon path:nth-child(1) {
  transform: translateY(0px);
}

.hamburger-icon path:nth-child(2) {
  transform: translateY(6px); /* Ajusta la distancia entre las líneas */
}

.hamburger-icon path:nth-child(3) {
  transform: translateY(12px); /* Ajusta la distancia entre las líneas */
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  background: white;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  min-width: 150px;
}

.dropdown-menu ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.dropdown-menu li {
  padding: 10px;
}

.dropdown-menu li a {
  text-decoration: none;
  color: #333;
  display: block;
}

.dropdown-menu li a:hover {
  background-color: #f0f0f0;
}

.dropdown-link {
  text-decoration: none;
  color: #333;
  display: block;
  padding: 10px;
  cursor: pointer; /* Hace que el cursor cambie a una mano al pasar */
}

.dropdown-link:hover {
  background-color: #f0f0f0;
}

/* Barra superior */
.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background-color: #00a499;
  color: white;
  padding: 1rem;
  z-index: 50;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 4rem;
  position: relative;
  z-index: 1; /* Menor que el del menú */
}

.header-container {
  flex: 1;
  display: flex;
  justify-content: left;
}

.header-title {
  font-family: "Bebas Neue Pro", sans-serif;
  font-size: 1.5rem;
  font-weight: bold;
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

.menu-item {
  position: relative;
}

.menu-link {
  text-decoration: none;
  color: white;
  background-color: #00897b;
  display: block;
  padding: 10px;
  cursor: pointer; /* Hace que el cursor cambie a una mano al pasar */
  margin: 5px 0;
  border-radius: 10px;
}

.menu-link:hover {
  background-color: #026c62;
}
</style>