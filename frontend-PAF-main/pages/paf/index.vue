<template>
  <div class="flex flex-col">
    <!-- Botón para volver debajo de la barra superior pero encima de la información de la PAF -->
    <div class="mt-4 ml-4">
      <button @click="volver" class="volver-button">
        Volver
      </button>
    </div>

    <!-- Información de la PAF -->
    <div class="info-container w-2/3 mt-6">
      <h1 v-if="paf.length > 0" class="section-title">Información de la PAF:</h1>

      <!-- Mostrar la lista de personas -->
      <div v-if="paf.length > 0">
        <div
          v-for="persona in paf"
          :key="persona.CodigoPaf"
          class="paf-container"
        >
          <p><strong>Código PAF:</strong> {{ persona.CodigoPaf }}</p>
          <p><strong>Código Asignatura:</strong> {{ persona.CodigoAsignatura }}</p>
          <p><strong>Nombre:</strong> {{ persona.Nombres }} {{ persona.PrimerApellido }} {{ persona.SegundoApellido }}</p>
          <p><strong>Asignatura:</strong> {{ persona.NombreAsignatura }}</p>
          <p><strong>Bloque:</strong> {{ persona.Bloque }}</p>
          <p><strong>Unidad Contratante:</strong> {{ persona.NombreUnidadContratante }}</p>

          <!-- Botón ubicado en la parte inferior -->
          <div class="flex justify-end mt-4">
            <button @click="dejarListaPaf(persona.CodigoPaf)" class="procesar-button">
              Dejar lista la PAF
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <p>Cargando datos o no se encontraron registros para la PAF.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { onMounted, ref } from 'vue';
import { useNuxtApp } from '#app';

const route = useRoute();
const router = useRouter();

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

const codigoPaf = ref(route.query.codigoPaf || '');
const paf = ref<any[]>([]);

const obtenerDatosPaf = async () => {
  try {
    if (!codigoPaf.value) return;
    const response = await $axios.get(`/contratos/codigo_paf/${codigoPaf.value}`);
    if (response.data) {
      paf.value = response.data.map((item: any) => ({
        CodigoPaf: item.PipelsoftData.IdPaf,
        CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
        Nombres: item.PipelsoftData.Nombres,
        NombreAsignatura: item.PipelsoftData.NombreAsignatura,
        PrimerApellido: item.PipelsoftData.PrimerApellido,
        SegundoApellido: item.PipelsoftData.SegundoApellido,
        CantidadHoras: item.PipelsoftData.CantidadHoras,
        NombreUnidadContratante: item.PipelsoftData.NombreUnidadContratante,
        Bloque: item.HistorialPafData.bloque,
      }));
    }
  } catch (error) {
    console.error('Error al obtener los datos de la PAF:', error);
  }
};

const volver = () => {
  router.go(-1);
};

const dejarListaPaf = async (codigoPaf: string) => {
  try {
    await $axios.put(`/historial/${codigoPaf}/actualizarBanderaAceptacion`, {
      nuevaBanderaAceptacion: 1,
    });
    router.push('/personas');
  } catch (error) {
    console.error('Error al procesar la PAF:', error);
  }
};

onMounted(() => {
  obtenerDatosPaf();
});
</script>

<style scoped>
/* Contenedor general */
.info-container {
  margin: auto;
  font-family: "Helvetica Neue LT", sans-serif;
}

/* Botón "Volver" */
.volver-button {
  background-color: #394049;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  font-family: "Bebas Neue Pro", sans-serif;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.volver-button:hover {
  background-color: #EA7600;
}

/* Contenedor de la información */
.paf-container {
  background-color: #f9f9f9;
  border: 1px solid #394049;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* Botón "Dejar lista la PAF" */
.procesar-button {
  background-color: #EA7600;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  font-family: "Bebas Neue Pro", sans-serif;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.procesar-button:hover {
  background-color: #C8102E;
}

/* Título de sección */
.section-title {
  font-size: 1.5rem;
  font-family: "Bebas Neue Pro", sans-serif;
  color: #EA7600;
  margin-bottom: 16px;
}
</style>
