<template>
  <div>
    <!-- Mostrar el menú solo en las páginas permitidas -->
    <Menu v-if="showMenu" />
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>

<script>
export default {
  computed: {
    showMenu() {
      // Rutas donde NO quieres que aparezca el menú
      const noMenuRoutes = ["/"];
      return !noMenuRoutes.includes(this.$route.path);
    },
  },
};
</script>

<style scoped>
/* Variables con colores institucionales */
:root {
  --primary-color: #EA7600; /* Color principal USACH */
  --secondary-color: #394049; /* Color secundario USACH */
  --accent-color: #C8102E; /* Complementario para énfasis */
  --text-color: #ffffff;
  --background-gradient: linear-gradient(135deg, #EA7600 0%, #394049 100%);
  --background-color: #f4f7f6; /* Fondo neutro */
}

/* Estilo del contenedor del encabezado */
.header-container {
  background: var(--background-gradient);
  padding: 20px;
  border-bottom: 4px solid var(--accent-color);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* Estilo del título del encabezado */
.header-title {
  font-size: 2rem;
  color: var(--text-color);
  margin: 10px 0;
  font-family: "Bebas Neue Pro", sans-serif;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
}

/* Estilo del contenido principal */
.main-content {
  padding: 20px;
  background-color: var(--background-color);
  font-family: "Helvetica Neue LT", sans-serif;
  color: var(--secondary-color);
}

/* Estilo del menú */
.menu {
  background-color: var(--primary-color);
  padding: 15px;
  color: var(--text-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.menu-item {
  margin: 10px 0;
  font-family: "Helvetica Neue LT", sans-serif;
  color: var(--text-color);
  text-decoration: none;
  font-weight: bold;
}

.menu-item:hover {
  color: var(--accent-color);
  text-decoration: underline;
}
</style>
