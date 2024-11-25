<template>
  <div class="flex">
    <!-- Botón para volver -->
    <div class="absolute top-4 left-4">
      <button @click="volver" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        Volver
      </button>
    </div>

    <!-- Tabla de horarios -->
    <div class="w-2/3 mt-12">
      <h1>Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}</h1>

      <div v-if="persona.length > 0">
        <div class="mb-4">
          <label for="semestre">Seleccionar Semestre:</label>
          <select id="semestre" v-model="semestreSeleccionado" class="ml-2">
            <option v-for="sem in semestres" :key="sem" :value="sem">{{ sem }}</option>
          </select>
        </div>

        <table class="w-full text-sm bg-white divide-y divide-gray-300 border">
          <thead>
            <tr class="bg-gray-200">
              <th class="px-4 py-2 font-medium text-gray-900 border-r">Módulo</th>
              <th v-for="dia in dias" :key="dia" class="px-4 py-2 font-medium text-gray-900 border-r">{{ dia }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(horario, index) in horarios" :key="index" :class="index % 2 === 0 ? 'bg-gray-100' : 'bg-white'">
              <td class="px-4 py-2 text-gray-700 border-r">{{ horario.modulo }}</td>
              <td v-for="dia in dias" :key="dia" class="px-4 py-2 text-gray-700 border-r">
                <div v-for="bloque in bloquesPorDia(dia, index + 1)" :key="bloque.nombre" :style="{ backgroundColor: bloque.color }" class="rounded p-1 text-black">
                  {{ bloque.nombre }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p>Cargando datos o no se encontraron registros para el RUN.</p>
      </div>
    </div>

    <!-- Fichas de historial, PAF y asignaturas -->
    <div class="w-1/3 pl-4 mt-12">
      <!-- Ficha del historial -->
      <div v-if="historialSeleccionado">
        <h2 class="font-semibold text-lg">PAF macheada</h2>
        <div class="p-4 mb-4 rounded shadow bg-yellow-200">
          <p><strong>Código PAF:</strong> {{ historialSeleccionado?.CodigoPaf }}</p>
          <p><strong>Código Asignatura:</strong> {{ historialSeleccionado?.CodigoAsignatura }}</p>
          <p><strong>Nombre Asignatura:</strong> {{ historialSeleccionado?.NombreAsignatura }}</p>
          <p><strong>Semestre:</strong> {{ historialSeleccionado?.semestre }}</p>
        </div>
      </div>

      <!-- Fichas de PAF -->
      <h2 class="font-semibold text-lg mt-8">PAF</h2>
      <div v-if="fichasPAF.length > 0">
        <div
          v-for="(p, index) in fichasPAF"
          :key="index"
          :style="{ 
            backgroundColor: fichaSeleccionadaPAF === p 
              ? 'lightblue' 
              : coloresPAF[index % coloresPAF.length] 
          }"
          class="p-4 mb-4 rounded shadow cursor-pointer"
          @click="fichaSeleccionadaPAF = p"
        >
          <p><strong>Nombres:</strong> {{ p.Nombres }}</p>
          <p><strong>Apellidos:</strong> {{ p.PrimerApellido }} {{ p.SegundoApellido }}</p>
          <p><strong>Jefatura:</strong> {{ p.NombreUnidadContratante }}</p>
          <p><strong>Cantidad de Horas:</strong> {{ p.CantidadHoras }}</p>
        </div>
      </div>
      <div v-else>
        <p>No se encontraron registros de PAF.</p>
      </div>

      <!-- Fichas de asignaturas -->
      <h2 class="font-semibold text-lg mt-8">Horario Asignatura</h2>
      <div v-if="fichasAsignaturas.length > 0">
        <div
          v-for="(p, index) in fichasAsignaturas"
          :key="index"
          :style="{ 
            backgroundColor: fichaSeleccionadaAsignatura === p 
              ? 'lightblue' 
              : coloresAsignaturas[index % coloresAsignaturas.length] 
          }"
          class="p-4 mb-4 rounded shadow cursor-pointer"
          @click="fichaSeleccionadaAsignatura = p"
        >
          <p><strong>Código de Asignatura:</strong> {{ p.codigo_asignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.nombre_asignatura }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.cupo }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
          <p><strong>Semestre:</strong> {{ p.semestre }}</p>
        </div>
      </div>
      <div v-else>
        <p>No se encontraron asignaturas filtradas.</p>
      </div>
            <!-- Botón ubicado en la parte inferior -->
            <div class="flex justify-end mt-4">
            <button v-if="fichaSeleccionadaPAF && fichaSeleccionadaAsignatura"
            @click="enviarSeleccion" class="bg-blue-500 text-white py-2 px-4 rounded">
            Enviar Selección
            </button>
          </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useNuxtApp } from '#app'

const historialSeleccionado = computed(() => persona.value.find((p) => p.ID !== 0) || null);

const fichasPAF = computed(() =>
  persona.value.filter((p) => 
    p.CodigoPaf !== historialSeleccionado.value?.CodigoPaf
  )
);

const fichasAsignaturas = computed(() =>
  persona1.value.filter((p) => 
    p.codigo_asignatura !== historialSeleccionado.value?.CodigoAsignatura
  )
);


const fichaSeleccionadaPAF = ref<Persona | null>(null);
const fichaSeleccionadaAsignatura = ref<Horario | null>(null);

const route = useRoute()
const router = useRouter()

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const run = ref(route.query.run || '')

interface Persona {
  CodigoPaf: string;
  CodigoAsignatura: string;
  NombreAsignatura: string;
  CantidadHoras: number;
  NombreUnidadContratante: string;
  Nombres: string;
  PrimerApellido: string;
  SegundoApellido: string;
  ID: number;
  semestre?: string; // Add the 'semestre' property
}

interface Horario {
  run: string;
  codigo_asignatura: string;
  nombre_asignatura: string;
  bloque?: string;
  cupo?: number;
  seccion?: string;
  semestre?: string;
}

const persona = ref<Persona[]>([])
const persona1 = ref<Horario[]>([])
console.log("persona1", persona1.value)
const colores = ['#C8E6C9', '#A5D6A7', '#81C784', '#66BB6A', '#4CAF50']
const dias = ['Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado']
const horarios = ref([
  { modulo: '08:15 - 09:35' },
  { modulo: '09:50 - 11:10' },
  { modulo: '11:25 - 12:45' },
  { modulo: '13:45 - 15:05' },
  { modulo: '15:20 - 16:40' },
  { modulo: '16:55 - 18:15' },
  { modulo: '18:45 - 20:05' },
  { modulo: '20:05 - 21:25' },
  { modulo: '21:25 - 22:45' }
])

const semestreSeleccionado = ref('')
const semestres = computed(() => [...new Set(persona1.value.map(p => p.semestre))])
const personaFiltrada = computed(() => persona1.value.filter(p => p.semestre === semestreSeleccionado.value))

const bloquesPorDia = (dia: string, modulo: number) => {
  // Mapear iniciales de días con sus nombres completos
  const inicialDia: { [key: string]: string } = {
    Lunes: 'L',
    Martes: 'M',
    Miércoles: 'W', // 'W' para distinguir de miércoles
    Jueves: 'J',
    Viernes: 'V',
    Sábado: 'S',
  };

  // Filtrar asignaturas para el día específico y el módulo actual
  return personaFiltrada.value
    .filter((p) => {
      if (!p.bloque) return false;

      // Separar bloques (por ejemplo: M2-M5-V1 se divide en ["M2", "M5", "V1"])
      const bloques = p.bloque.split('-');

      // Buscar un bloque que coincida con el día y el módulo
      return bloques.some((b) => {
        const diaBloque = b.charAt(0); // Primera letra es el día (M, V, etc.)
        const moduloBloque = b.slice(1); // Resto es el número de módulo

        // Comparar con el día actual y el módulo actual
        return inicialDia[dia] === diaBloque && parseInt(moduloBloque) === modulo;
      });
    })
    .map((p) => ({
      nombre: p.nombre_asignatura,
      seccion: p.seccion,
      color: colores[persona1.value.indexOf(p) % colores.length],
    }));
};


const obtenerDatosPersona = async () => {
  try {
    const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`)
    console.log('Datos de la persona:', response.data)
    const response1 = await $axios.get(`/profesorDB/${run.value}`)
    persona1.value = response1.data
    console.log("aaaa", persona1.value)
    console.log("aaaa", response1.data)
    console.log("bbbb", response1.data.bloque)
    persona.value = response.data.map((item: any) => ({
      CodigoPaf: item.PipelsoftData.CodigoPAF,
      CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
      Nombres: item.PipelsoftData.Nombres,
      NombreAsignatura: item.PipelsoftData.NombreAsignatura,
      PrimerApellido: item.PipelsoftData.PrimerApellido,
      SegundoApellido: item.PipelsoftData.SegundoApellido,
      CantidadHoras: item.PipelsoftData.CantidadHoras,
      NombreUnidadContratante: item.PipelsoftData.NombreUnidadContratante,
      Bloque: response1.data.bloque,
      Cupo: response1.data.cupo,
      Seccion: response1.data.seccion,
      Semestre: response1.data.semestre,
      ID: item.HistorialPafData.ID,
      semestre: item.HistorialPafData.semestre
    }))
  } catch (error) {
    console.error('Error al obtener los datos:', error)
  }
}

const enviarSeleccion = async () => {
  if (!fichaSeleccionadaPAF || !fichaSeleccionadaAsignatura) {
    alert('Por favor selecciona una ficha de PAF y una de asignatura.');
    return;
  }

  try {
    const codigoPAF = fichaSeleccionadaPAF.value?.CodigoPaf; // Ajustar si el código es otro campo
    const data = {
      run: fichaSeleccionadaAsignatura?.value?.run || '', // Cambiar según lo que necesitas enviar
      semestre: fichaSeleccionadaAsignatura?.value?.semestre || '',
      codigo_asignatura: fichaSeleccionadaAsignatura?.value?.codigo_asignatura || '',
      nombre_asignatura: fichaSeleccionadaAsignatura?.value?.nombre_asignatura || '',
      seccion: fichaSeleccionadaAsignatura?.value?.seccion || '',
      cupo: fichaSeleccionadaAsignatura?.value?.cupo || 0, 
      bloque: fichaSeleccionadaAsignatura?.value?.bloque || '', 
    };
    console.log('Datos a enviar:', data);
    await $axios.post(`/historial/post/${codigoPAF}`, data);
    alert('Datos enviados correctamente.');
  } catch (error) {
    console.error('Error al enviar los datos:', error);
    alert('Hubo un error al enviar los datos.');
  }
};


const volver = () => {
  router.go(-1)
}

onMounted(() => {
  obtenerDatosPersona()
})
const coloresPAF = ['#FFCDD2', '#F8BBD0', '#E1BEE7', '#D1C4E9', '#C5CAE9'];
const coloresAsignaturas = ['#C8E6C9', '#A5D6A7', '#81C784', '#66BB6A', '#4CAF50'];
</script>


<style scoped>
table {
  width: 100%;
  table-layout: auto;
  border-collapse: collapse;
}

th,
td {
  border: 1px solid #d1d5db; /* Borde para diferenciar módulos */
}

button {
  position: absolute;
  top: 4rem; /* Ajusta la posición según lo necesario */
  left: 1rem; /* Ajusta la distancia desde la izquierda */
  z-index: 10; /* Asegura que el botón esté por encima de otros elementos */
}

.flex.justify-end {
  margin-top: auto;
}
</style>
