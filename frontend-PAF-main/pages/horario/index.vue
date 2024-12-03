<template>
  <div class="flex" @click="cerrarTodosLosBloques">
    <!-- Botón para volver -->
    <div class="mt-4 ml-4">
      <button @click="volver" class="volver-button">Volver</button>
    </div>

    <!-- Tabla de horarios -->
    <div class="w-2/3 mt-12 relative">
      <h1 class="section-title">
        Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}
      </h1>

      <div v-if="persona.length > 0">
        <div class="mb-4">
          <label for="semestre">Seleccionar Semestre:</label>
          <select id="semestre" v-model="semestreSeleccionado" class="select-input" @click.stop>
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
            <tr
              v-for="(horario, index) in horarios"
              :key="index"
              :class="{ alternado: index % 2 === 0 }"
            >
              <td>{{ horario.modulo }}</td>
              <td v-for="dia in dias" :key="dia" class="relative">
                <div>
                  <!-- Mostrar "Tope Horario" si hay más de 2 bloques -->
                  <div v-if="bloquesPorDia(dia, index + 1).length > 2">
                    <button
                      class="tope-horario-button"
                      @click.stop="mostrarBloques(dia, index + 1)"
                    >
                      Tope Horario
                    </button>
                    <div
                      v-if="isBloqueVisible(dia, index + 1)"
                      class="detalle-horarios-popup"
                      @click.stop
                    >
                      <ul>
                        <li
                          v-for="(bloque, idx) in bloquesPorDia(dia, index + 1)"
                          :key="bloque.nombre"
                          class="bloque"
                          :style="{ backgroundColor: bloque.color, marginBottom: '10px' }"
                        >
                          <label>
                            <input
                              type="checkbox"
                              :value="`${dia}${index + 1}/${bloque.nombre}/${bloque.seccion}/${bloque.ID}/${bloque.cupo}/${bloque.codigo_asignatura}/${bloque.tipo}/${bloque.run}/${bloque.semestre}`"
                              v-model="bloquesSeleccionados"
                            />
                            {{ bloque.nombre }} <br />
                            Sección: {{ bloque.seccion }}
                          </label>
                        </li>
                      </ul>
                    </div>
                  </div>

                  <!-- Mostrar bloques directamente si son 2 o menos -->
                  <div v-else>
                    <div
                      v-for="(bloque, bloqueIndex) in bloquesPorDia(dia, index + 1)"
                      :key="`${bloque.nombre}-${bloqueIndex}`"
                      class="bloque"
                      :style="{ backgroundColor: bloque.color }"
                    >
                      <label>
                        <input
                          type="checkbox"
                          :value="`${dia}${index + 1}/${bloque.nombre}/${bloque.seccion}/${bloque.ID}/${bloque.cupo}/${bloque.codigo_asignatura}/${bloque.tipo}/${bloque.run}/${bloque.semestre}`"
                          v-model="bloquesSeleccionados"
                        />
                        {{ bloque.nombre }} <br /> Sección: {{ bloque.seccion }}
                      </label>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else>
        <p class="info-text">
          Cargando datos o no se encontraron registros para el RUN.
        </p>
      </div>
    </div>

    <!-- Fichas y botón de envío -->
    <div class="w-1/3 pl-4 mt-12">
      <h2 class="sub-title">Selecciona PAF y Asignatura</h2>

      <h2 class="sub-title">PAF con Asignatura</h2>
      <div v-if="fichasPAFMatch.length > 0">
        <div
          v-for="(p, index) in fichasPAFMatch"
          :key="index"
          class="card"
          :style="{ backgroundColor: '#FFCC80' }"
        >
          <p><strong>Código PAF:</strong> {{ p.CodigoPaf }}</p>
          <p><strong>Unidad Menor:</strong> {{ p.NombreUnidadMenor }}</p>
          <p><strong>Codigo de Asignatura:</strong> {{ p.CodigoAsignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.NombreAsignatura }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.cupo }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
        </div>
      </div>
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

        </div>
      </div>

      <h2 class="sub-title">Horario Asignatura</h2>
      <div v-if="fichasAsignaturas.length > 0">
        <div
          v-for="(p, index) in fichasAsignaturas"
          :key="index"
          class="card"
          :style="{ backgroundColor: fichaSeleccionadaAsignatura === p ? '#B3E5FC' : coloresAsignaturas[index % coloresAsignaturas.length] }"
          
        >
          <p><strong>Código de Asignatura:</strong> {{ p.codigo_asignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.nombre_asignatura }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.Cupo }}</p>
        </div>
      </div>

      <div class="flex justify-end mt-4">
        <button
          v-if="fichaSeleccionadaPAF && bloquesSeleccionados.length > 0"
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

const historialSeleccionado = computed(() => persona.value.filter((p) => p.ID !== 0) || null);

// Cerrar todos los popups
const cerrarTodosLosBloques = () => {
  Object.keys(visibleBloques.value).forEach((key) => {
    visibleBloques.value[key] = false;
  });
};


// Agregar y remover el evento global
onMounted(() => {
  document.addEventListener('click', cerrarTodosLosBloques);
});

onUnmounted(() => {
  document.removeEventListener('click', cerrarTodosLosBloques);
});

const visibleBloques = ref<Record<string, boolean>>({});

// Mostrar u ocultar los bloques según día y módulo
const mostrarBloques = (dia: string, modulo: number) => {
  const key = `${dia}-${modulo}`;
  visibleBloques.value[key] = !visibleBloques.value[key];
};

// Verificar si los bloques son visibles
const isBloqueVisible = (dia: string, modulo: number) => {
  const key = `${dia}-${modulo}`;
  return !!visibleBloques.value[key];
};

const fichasPAF = computed(() =>
  persona.value.filter((p) => 
    !historialSeleccionado.value.some(h => 
      h.codigo == p.CodigoAsignatura || h.paf === p.CodigoPaf || h.CodigoA === p.CodigoAsignatura
    )
  )
);

const fichasPAFMatch = computed(() =>
  persona.value.filter((p) => 
    historialSeleccionado.value.some(h => 
      h.codigo == p.CodigoAsignatura || h.paf === p.CodigoPaf || h.CodigoA === p.CodigoAsignatura
    )
  )
);

const fichasAsignaturas = computed(() =>
  persona1.value.filter((p) => 
    !historialSeleccionado.value.some((historial) => 
      historial.codigo === p.codigo_asignatura
    )
  )
);

const fichasAsignaturasNo = computed(() =>
  persona1.value.filter((p) => 
    historialSeleccionado.value.some((historial) => 
      historial.codigo === p.codigo_asignatura
    )
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
  NombreUnidadMenor: string;
  Nombres: string;
  PrimerApellido: string;
  SegundoApellido: string;
  ID: number;
  semestre?: string; // Add the 'semestre' property
  codigo: string;
  paf: number;
  bloque: string;
  cupo: number;
  seccion: string;
  CodigoA: string;
}

interface Horario {
  run: string;
  codigo_asignatura: string;
  nombre_asignatura: string;
  bloque?: string;
  Cupo?: number;
  seccion?: string;
  semestre?: string;
  ID: number;
  tipo?: string;
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
  const inicialDia: { [key: string]: string } = {
    Lunes: 'L',
    Martes: 'M',
    Miércoles: 'W',
    Jueves: 'J',
    Viernes: 'V',
    Sábado: 'S',
  };

  // Filtrar bloques horarios en función de la ficha seleccionada
  return personaFiltrada.value
    .filter((p) => {
      if (!p.bloque) return false;

      // Excluir bloques asociados con las fichas seleccionadas en 'fichasAsignaturasNo'
      if (fichasAsignaturasNo.value.some((ficha) => ficha.codigo_asignatura === p.codigo_asignatura)) {
        return false;
      }

      // Filtrar por día y módulo
      const bloques = p.bloque.split('-');
      return bloques.some((b) => {
        const diaBloque = b.charAt(0);
        const moduloBloque = b.slice(1);
        return inicialDia[dia] === diaBloque && parseInt(moduloBloque) === modulo;
      });
    })
    .map((p) => ({
      nombre: p.nombre_asignatura,
      seccion: p.seccion,
      color: colores[persona1.value.indexOf(p) % colores.length],
      ID: p.ID,
      cupo: p.Cupo,
      codigo_asignatura: p.codigo_asignatura,
      tipo: p.tipo,
      run: p.run,
      semestre: p.semestre,
      // Crear el formato "V3-M2" con la inicial del día y el número de módulo
      bloque: (p.bloque ?? '').split('-')
        .map(b => `${b.charAt(0)}${b.slice(1)}`) // Mapear para obtener "V3", "M2", etc.
        .join('-')
    }));
};


const obtenerDatosPersona = async () => {
  try {
    const response = await $axios.get(`/pipelsoft/contratos-run/${run.value}`);
    console.log("response", response.data)
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
      NombreUnidadMenor: item.PipelsoftData.NombreUnidadMenor,
      Bloque: response1.data.bloque,
      Cupo: response1.data.cupo,
      Seccion: response1.data.seccion,
      Semestre: response1.data.semestre,
      ID: item.HistorialPafData.ID,
      semestre: item.HistorialPafData.semestre,
      codigo: item.HistorialPafData.codigo_asignatura,
      paf: item.HistorialPafData.CodigoPaf,
      bloque: item.HistorialPafData.bloque,
      cupo: item.HistorialPafData.cupo,
      seccion: item.HistorialPafData.seccion,
      CodigoA: item.HistorialPafData.CodigoAsignatura,
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
    let lista = bloquesSeleccionadosString.value.split(',').map(item => {
    return item
        .replace("Lunes", "L")
        .replace("Martes", "M")
        .replace("Miércoles", "W")
        .replace("Jueves", "J")
        .replace("Viernes", "V")
        .replace("Sábado", "S")
        .replace("Domingo", "D");
});
console.log(lista)
const grouped: { [key: string]: string[] } = {};

lista.forEach(item => {
  const partes = item.split('/');

  const [bloque, nombreAsignatura, seccion, Id, cupo, codigoAsignatura, tipo, run, semestre] = partes;
  const rest = partes.slice(1).join('/');

  if (!grouped[rest]) {
    grouped[rest] = [];
  }
  grouped[rest].push(bloque);
});

// Formatear el resultado en el formato de la imagen
const resultado = Object.entries(grouped).map(([rest, bloques]) => {
  const [nombreAsignatura, seccion, Id, cupo, codigoAsignatura, tipo, run, semestre] = rest.split('/');

  return {
    profesor: {
      run,
      semestre,
      tipo,
      nombre_asignatura: nombreAsignatura,
      seccion,
      cupo: parseInt(cupo), // Convertir cupo a número
      bloque: `Bloque ${bloques.join(", ")}`, // Prefijo "Bloque" antes de los bloques agrupados
      codigo_asignatura: codigoAsignatura
    },
    bloque: bloques.map(b => `${codigoAsignatura} ${seccion} ${cupo} ${b}`)
  };
});
const result = resultado.map(item => {
  if (item.bloque.length >= 2) {
    // Une los elementos con un "-"
    const merged = item.bloque[0] + "-" + item.bloque.slice(1).join("-").split(" ")[3];
    return merged;
  }
  return item.bloque[0];
});

console.log(result);


// Mostrar el resultado en consola


    // si bloqueseleccionadosString tiene un elemento Miercoles Trasformarlo en W y los demas tipo Lunes o Matrtes dejarlos como L o M
    
    const data = resultado[0];
    data.bloque = result;
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

.tope-horario-button {
  background-color: #C8102E;
  color: white;
  padding: 5px 10px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  position: relative;
}

.tope-horario-button:hover {
  background-color: #A50D20;
}

.detalle-horarios-popup {
  position: absolute;
  top: -10px;
  left: 0;
  z-index: 10;
  width: 200px;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: #fff;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
  padding: 10px;
}

</style>