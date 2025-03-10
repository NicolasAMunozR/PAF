<template>
  <!-- Botón para volver debajo de la barra superior pero encima de la información de la PAF -->
  <div class="mt-4 ml-4">
    <button @click="volver" class="volver-button">
      Volver
    </button>
  </div>
  <div class="container">
    <!-- Mostrar la lista de personas -->
    <div v-if="filteredPersonas.length > 0">
      <div v-for="persona in filteredPersonas" :key="persona.CodigoPaf" class="paf-container">
        <h1 v-if="filteredPersonas.length > 0" class="section-title">Información de la PAF:</h1>
        <p><strong>Código PAF:</strong> {{ persona.CodigoPaf }}</p>
        <p><strong>Run:</strong> {{ persona.Run }}</p>
        <p><strong>Código Asignatura PAF:</strong> {{ persona.CodigoAsignatura }}</p>
        <p><strong>Nombre:</strong> {{ persona.Nombres }} {{ persona.PrimerApellido }} {{ persona.SegundoApellido }}</p>
        <p><strong>Asignatura:</strong> {{ persona.NombreAsignatura }}</p>
        <p><strong>Semestre PAF:</strong> {{ persona.SemestrePaf }}</p>
        <p><strong>Unidad Menor:</strong> {{ persona.NombreUnidadMenor }}</p>
        <p><strong>Unidad Mayor:</strong> {{ persona.NombreUnidadMayor }}</p>
        <p><strong>Horas Asignatura PAF:</strong> {{ persona.HorasAsignatura }}</p>
        <p><strong>Bloque:</strong> {{ persona.bloque }}</p>
        <p><strong>Código Asignatura Asociadas:</strong> {{ persona.CodigoA }}</p>
        <p><strong>Cupo:</strong> {{ persona.cupo }}</p>
        <p><strong>Sección:</strong> {{ persona.seccion }}</p>
        <p><strong>Semestre Asignatura:</strong> {{ persona.semestre1 }}</p>
        <!-- Botón ubicado en la parte inferior -->
        <div class="flex justify-end mt-4">
          <button v-if="persona.Aceptada === 0" @click="dejarListaPaf(persona.CodigoPaf)" class="procesar-button">Dejar lista la PAF</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>Cargando datos o no se encontraron registros para la PAF.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { ref, computed, onMounted } from 'vue';
import { useNuxtApp } from '#app';

const route = useRoute();
const router = useRouter();
const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default };

const codigoPaf = ref(route.query.codigoPaf || '');
const paf = ref<any[]>([]);

const filtros = ref({ semestre: '' });
const sortBy = ref('nombres');
const sortOrder = ref('asc');

const filteredPersonas = computed(() => {
  let filtered = paf.value.filter(contrato => 
    (contrato.SemestrePaf || '').toLowerCase().includes((filtros.value.semestre || '').toLowerCase())
  );

  if (sortBy.value) {
    filtered = filtered.sort((a, b) => {
      const compareA = a[sortBy.value];
      const compareB = b[sortBy.value];
      if (compareA < compareB) return sortOrder.value === 'asc' ? -1 : 1;
      if (compareA > compareB) return sortOrder.value === 'asc' ? 1 : -1;
      return 0;
    });
  }

  return filtered;
});

const obtenerDatosPaf = async () => {
  try {
    if (!codigoPaf.value) return;
    const response = await $axios.get(`/api/paf-en-linea/pipelsoft/obtenerContratos/mostrarTodo/idPaf/${codigoPaf.value}`);
    if (response.data) {
      // Agrupar datos por Run para combinar información
      const groupedData: Record<string, any> = {};

      response.data.forEach((item: any) => {
        const run = item.PipelsoftData.RunEmpleado;
        if (!groupedData[run]) {
          groupedData[run] = {
            CodigoPaf: item.PipelsoftData.IdPaf,
            Run: item.PipelsoftData.RunEmpleado,
            Nombres: item.PipelsoftData.Nombres,
            PrimerApellido: item.PipelsoftData.PrimerApp,
            SegundoApellido: item.PipelsoftData.SegundoApp,
            NombreUnidadMenor: item.PipelsoftData.NombreUnidadMenor,
            NombreUnidadMayor: item.PipelsoftData.NombreUnidadMayor,
            HorasAsignatura: item.PipelsoftData.HorasAsignatura,
            SemestrePaf: item.PipelsoftData.Semestre,
            CodigoAsignatura: [],
            NombreAsignatura: [],
            bloque: [],
            CodigoA: [],
            cupo: [],
            seccion: [],
            semestre1: []
          };
        }

        const bloquesArray = item.HistorialPafData.Bloque || [];
        groupedData[run].CodigoAsignatura.push(item.PipelsoftData.CodigoAsignatura);
        groupedData[run].NombreAsignatura.push(item.PipelsoftData.NombreAsignatura);
        if (bloquesArray.length > 0) {
          groupedData[run].bloque.push(...bloquesArray.map((bloque: any) => bloque.bloques));
          groupedData[run].CodigoA.push(...bloquesArray.map((bloque: any) => bloque.codigoAsignatura));
          groupedData[run].cupo.push(...bloquesArray.map((bloque: any) => bloque.cupos));
          groupedData[run].seccion.push(...bloquesArray.map((bloque: any) => bloque.seccion));
          groupedData[run].semestre1.push(...bloquesArray.map((bloque: any) => bloque.semestre));
        }
      });

      // Transformar los datos agrupados en un array para usar en la vista
      paf.value = Object.values(groupedData).map((persona: any) => ({
        ...persona,
        CodigoAsignatura: persona.CodigoAsignatura.join(', '),
        NombreAsignatura: persona.NombreAsignatura.join(', '),
        bloque: persona.bloque.join('/'),
        CodigoA: persona.CodigoA.join('/'),
        cupo: persona.cupo.join('/'),
        seccion: persona.seccion.join('/'),
        semestre1: persona.semestre1.join('/')
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
    await $axios.put(`/api/paf-en-linea/historial/${codigoPaf}/actualizarBanderaAceptacion`, {
      nuevaBanderaAceptacion: 1,
    });
    router.push('/paf-en-linea/personas');
  } catch (error) {
    console.error('Error al procesar la PAF:', error);
  }
};

onMounted(() => {
  obtenerDatosPaf();
});
</script>

<style scoped>
.container {
  display: -ms-grid;
  gap: 1rem;
  max-width: 70%;
  justify-content: center;
  margin: auto;
}

.info-container {
  margin: auto;
  font-family: "Helvetica Neue LT", sans-serif;
}

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

.paf-container {
  background-color: #f9f9f9;
  border: 1px solid #394049;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

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

.section-title {
  font-size: 1.5rem;
  font-family: "Bebas Neue Pro", sans-serif;
  color: #EA7600;
  margin-bottom: 16px;
}
</style>
