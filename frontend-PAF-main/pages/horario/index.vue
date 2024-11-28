<template>
  <div class="flex">
    <!-- Botón para volver -->
    <div class="mt-4 ml-4">
      <button @click="volver" class="volver-button">Volver</button>
    </div>

    <!-- Tabla de horarios -->
    <div class="w-2/3 mt-12">
      <h1 class="section-title">
        Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}
      </h1>

      <div v-if="persona.length > 0">
        <div class="mb-4">
          <label for="semestre">Seleccionar Semestre:</label>
          <select id="semestre" v-model="semestreSeleccionado" class="select-input">
            <option v-for="sem in semestres" :key="sem" :value="sem">{{ sem }}</option>
          </select>
        </div>

        <table class="tabla-horarios">
          <thead>
            <tr>
              <th>Módulo</th>
              <th v-for="dia in dias" :key="dia">{{ dia }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(horario, index) in horarios" :key="index" :class="{ alternado: index % 2 === 0 }">
              <td>{{ horario.modulo }}</td>
              <td v-for="dia in dias" :key="dia">
                <div v-for="bloque in bloquesPorDia(dia, index + 1)" :key="bloque.nombre" class="bloque" :style="{ backgroundColor: bloque.color }">
                  <label>
                    <input
                      type="checkbox"
                      :value="`${dia}-${index + 1}`"
                      v-model="bloquesSeleccionados"
                    />
                    {{ bloque.nombre }}
                  </label>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p class="info-text">Cargando datos o no se encontraron registros para el RUN.</p>
      </div>
    </div>

    <!-- Fichas y botón de envío -->
    <div class="w-1/3 pl-4 mt-12">
      <!-- Fichas de historial, PAF y asignaturas -->
      <h2 class="sub-title">Selecciona PAF y Asignatura</h2>

      <!-- Fichas de PAF -->
      <h2 class="sub-title">PAF</h2>
      <div v-if="fichasPAF.length > 0">
        <div
          v-for="(p, index) in fichasPAF"
          :key="index"
          class="card"
          :style="{ backgroundColor: fichaSeleccionadaPAF === p ? '#B3E5FC' : coloresPAF[index % coloresPAF.length] }"
          @click="fichaSeleccionadaPAF = p"
        >
          <p><strong>Código PAF:</strong> {{ p.CodigoPaf }}</p>
          <p><strong>Unidad Menor:</strong> {{ p.NombreUnidadMenor }}</p>
          <p><strong>Cantidad de Horas:</strong> {{ p.CantidadHoras }}</p>
        </div>
      </div>

      <!-- Fichas de asignaturas -->
      <h2 class="sub-title">Horario Asignatura</h2>
      <div v-if="fichasAsignaturas.length > 0">
        <div
          v-for="(p, index) in fichasAsignaturas"
          :key="index"
          class="card"
          :style="{ backgroundColor: fichaSeleccionadaAsignatura === p ? '#B3E5FC' : coloresAsignaturas[index % coloresAsignaturas.length] }"
          @click="fichaSeleccionadaAsignatura = p"
        >
          <p><strong>Código de Asignatura:</strong> {{ p.codigo_asignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.nombre_asignatura }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
        </div>
      </div>

      <!-- Botón para enviar selección -->
      <div class="flex justify-end mt-4">
        <button
          v-if="fichaSeleccionadaPAF && fichaSeleccionadaAsignatura && bloquesSeleccionados.length > 0"
          @click="enviarSeleccion"
          class="procesar-button"
        >
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
const bloquesSeleccionados = ref<string[]>([]);
const route = useRoute()
const router = useRouter()

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const run = ref(route.query.run || '')

interface Persona {
  CodigoPaf: number;
  CodigoAsignatura: string;
  NombreAsignatura: string;
  CantidadHoras: number;
  NombreUnidadMenor: string;
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

      // Separar bloques (por ejemplo: M2-M5-V1 se divide en ["M2", "M5", "V1"]).
      const bloques = p.bloque.split('-');

      // Buscar un bloque que coincida con el día y el módulo actual
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
      // Crear el formato "V3-M2" con la inicial del día y el número de módulo
      bloque: (p.bloque ?? '').split('-')
        .map(b => `${b.charAt(0)}${b.slice(1)}`) // Mapear para obtener "V3", "M2", etc.
        .join('-')
    }));
};


const obtenerDatosPersona = async () => {
  try {
    const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`);
    const response1 = await $axios.get(`/profesorDB/${run.value}`);
    persona1.value = response1.data;
    console.log("persona1", persona1.value)
    persona.value = response.data.map((item: any) => ({
      CodigoPaf: item.PipelsoftData.IdPaf,
      CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
      Nombres: item.PipelsoftData.Nombres,
      NombreAsignatura: item.PipelsoftData.NombreAsignatura,
      PrimerApellido: item.PipelsoftData.PrimerApp,
      SegundoApellido: item.PipelsoftData.SegundoApp,
      CantidadHoras: item.PipelsoftData.CantidadHorasPaf,
      NombreUnidadMenor: item.PipelsoftData.NombreUnidadMenor,
      Bloque: response1.data.bloque,
      Cupo: response1.data.cupo,
      Seccion: response1.data.seccion,
      Semestre: response1.data.semestre,
      ID: item.HistorialPafData.ID,
      semestre: item.HistorialPafData.semestre,
    }))
    console.log("persona", persona.value);

    // Identificar el semestre más reciente
    const semestresDisponibles = persona1.value.map(p => p.semestre).filter(Boolean);
    const semestreReciente = semestresDisponibles.sort((a, b) => {
      if (a && b) {
        return a > b ? -1 : 1;
      }
      return 0;
    })[0];
    
    // Establecer el semestre más reciente como seleccionado
    if (semestreReciente) {
      semestreSeleccionado.value = semestreReciente;
    }
  } catch (error) {
    console.error('Error al obtener los datos:', error);
  }
};

const enviarSeleccion = async () => {
  if (!fichaSeleccionadaPAF || !fichaSeleccionadaAsignatura) {
    alert('Por favor selecciona una ficha de PAF y una de asignatura.');
    return;
  }

  try {
    const codigoPAF = fichaSeleccionadaPAF.value?.CodigoPaf; // Ajustar si el código es otro campo
    console.log('Código PAF:', codigoPAF);
    const bloquesSeleccionadosString = computed(() => bloquesSeleccionados.value.join(','))
    console.log(bloquesSeleccionadosString.value)
    // si bloqueseleccionadosString tiene un elemento Miercoles Trasformarlo en W y los demas tipo Lunes o Matrtes dejarlos como L o M
    
    const data = {
      run: fichaSeleccionadaAsignatura?.value?.run || '', // Cambiar según lo que necesitas enviar
      semestre: fichaSeleccionadaAsignatura?.value?.semestre || '',
      codigo_asignatura: fichaSeleccionadaAsignatura?.value?.codigo_asignatura || '',
      nombre_asignatura: fichaSeleccionadaAsignatura?.value?.nombre_asignatura || '',
      seccion: fichaSeleccionadaAsignatura?.value?.seccion || '',
      cupo: fichaSeleccionadaAsignatura?.value?.cupo || 0, 
      bloque: bloquesSeleccionados.value, 
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
/* Estilo general */
.tabla-horarios {
  width: 100%;
  border-collapse: collapse;
}

.tabla-horarios th,
.tabla-horarios td {
  padding: 8px;
  border: 1px solid #ccc;
}

.tabla-horarios th {
  background-color: #394049;
  color: white;
}

.bloque {
  padding: 4px;
  border-radius: 4px;
  color: black;
}

.card {
  padding: 10px;
  border: 1px solid #394049;
  border-radius: 8px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

.volver-button {
  background-color: #EA7600;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  font-family: "Bebas Neue Pro", sans-serif;
  cursor: pointer;
}

.volver-button:hover {
  background-color: #C8102E;
}

.procesar-button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border-radius: 4px;
  font-family: "Bebas Neue Pro", sans-serif;
}

.procesar-button:hover {
  background-color: #388E3C;
}

.sub-title {
  color: #EA7600;
  font-family: "Bebas Neue Pro", sans-serif;
  font-size: 1.2rem;
  margin-bottom: 10px;
}

.info-text {
  color: #394049;
  font-family: "Helvetica Neue LT", sans-serif;
}
</style>