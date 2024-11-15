<template>
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
</template>

<script>
export default {
  data() {
    return {
      paf: {
        codigo_paf: 'PAF001',
        run: '11111111-1',
      },
    };
  },
  methods: {
    goBack() {
      this.$router.go(-1); // Vuelve a la página anterior
    },
    async markPafAsReady() {
      try {
        const response = await this.$axios.post('/historial', {
          codigo_paf: this.paf.codigo_paf,
        });

        console.log('Historial creado:', response.data);
      } catch (error) {
        console.error('Error al crear el historial:', error);
      }
    },
  },
};
</script>
