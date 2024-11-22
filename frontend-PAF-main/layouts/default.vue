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
      return !noMenuRoutes.includes(this.$route.path); // Retorna true si la ruta no está en la lista
    },
  },
};
</script>

<style scoped>
/* Estilos existentes */
:root {
  --primary-color: #4A90E2;
  --secondary-color: #50E3C2;
  --accent-color: #B8E986;
  --text-color: #ffffff;
  --background-gradient: linear-gradient(135deg, #4A90E2 0%, #50E3C2 100%);
}

.header-container {
  background: var(--background-gradient);
  padding: 20px;
  border-bottom: 4px solid var(--accent-color);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.header-title {
  font-size: 2rem;
  color: var(--text-color);
  margin: 10px 0;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
}

.main-content {
  padding: 20px;
  background-color: #f4f7f6;
}
</style>

