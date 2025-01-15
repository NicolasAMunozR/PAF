<template>
  <!-- Botón para volver -->
  <div class="flex justify-between mt-4">
  <!-- Botón Volver alineado a la izquierda -->
  <button @click="volver" class="volver-button">Volver</button>
  
  <!-- Botón Enviar Selección alineado a la derecha, solo si hay una selección -->
  <button
    v-if="fichaSeleccionadaPAF && bloquesSeleccionados.length > 0"
    @click="enviarSeleccion"
    class="procesar-button"
  >
    Enviar Selección
  </button>
  <br>
  <button
    v-if="fichaSeleccionadaPAF && bloquesSeleccionados.length > 0"
    @click="enviarSeleccion1"
    class="procesar-button"
  >
    Enviar Selección Con Comentario
  </button>
  <div v-if="fichaSeleccionadaPAF && bloquesSeleccionados.length > 0 && mostrarDialogo" class="modal-overlay">
  <div class="modal">
    <h3>Ingrese su comentario</h3>
    <textarea v-model="comentario" class="textarea-comentario" placeholder="Escriba su comentario aquí..."></textarea>
    <div class="modal-actions">
      <button @click="cancelarEnvio" class="cancel-button">Cancelar</button>
      <button @click="enviarSeleccion" class="confirm-button">Enviar</button>
    </div>
  </div>
</div>
</div>
 <div class="container">
  <div class="table-container">
  <div class="flex flex-wrap" @click="cerrarTodosLosBloques">
    <!-- Tabla de horarios -->
    <div class="w-full md:w-2/3 mt-12 relative">
      <h1 class="section-title">
        Horario para: {{ persona[0]?.Nombres }} {{ persona[0]?.PrimerApellido }} {{ persona[0]?.SegundoApellido }}
      </h1>
      <h1>Run: {{ persona[0]?.Run }}</h1>

      <div v-if="persona.length > 0">
        <div class="mb-4">
          <label for="semestre">Seleccionar Semestre:</label>
          <select id="semestre" v-model="semestreSeleccionado" class="select-input" @click.stop @change="limpiarSelecciones">
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
                          :style="{ backgroundColor: bloque.color }"
                        >
                          <label>
                            <input
                              type="checkbox"
                              :value="`${dia}${index + 1}/${bloque.nombre}/${bloque.seccion}/${bloque.ID}/${bloque.cupo}/${bloque.codigo_asignatura}/${bloque.tipo}/${bloque.run}/${bloque.semestre}`"
                              v-model="bloquesSeleccionados"
                              :disabled="bloque.color === '#FFCC80'"
                            />
                            {{ bloque.nombre }} <br />
                            Sección: {{ bloque.seccion }}
                          </label>
                        </li>
                      </ul>
                    </div>
                  </div>

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
                          :disabled="bloque.color === '#FFCC80'"
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
      </div>
    </div>
  </div>
    <!-- Fichas y botón de envío -->
   <div class="fichas-container">
    <div class="w-full md:w-1/3 pl-4 mt-12">
    <br /><br /><br /><br />
      <h2 class="sub-title">Selecciona PAF y Asignatura</h2>

      <h2 class="sub-title">PAF con Asignatura</h2>
      <div v-if="fichasPAFMatch.length > 0" class="flex flex-wrap">
        <div
          v-for="(p, index) in fichasPAFMatch"
          :key="index"
          class="card"
          :style="{ backgroundColor: '#FFCC80', margin: '8px', flex: '1 1 48%' }"
        >
          <p><strong>Código PAF:</strong> {{ p.CodigoPaf }}</p>
          <p><strong>Semestre PAF:</strong> {{ p.SemestrePaf }}</p>
          <p><strong>Unidad Menor:</strong> {{ p.NombreUnidadMenor }}</p>
          <p><strong>Codigo de Asignatura:</strong> {{ p.CodigoA }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.cupo }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
          <p><strong>Estado del Proceso:</strong> {{ p.DesEstado }}</p>
        </div>
      </div>

      <h2 class="sub-title">PAF</h2>
      <div v-if="fichasPAF.length > 0" class="flex flex-wrap">
        <div
        v-if="Number(detalle) === 2"
          v-for="(p, index) in fichasPAF"
          :key="index"
          class="card"
          :style="{ backgroundColor: fichaSeleccionadaPAF === p ? '#B3E5FC' : coloresPAF[index % coloresPAF.length], margin: '8px', flex: '1 1 48%' }"
          @click="fichaSeleccionadaPAF = p"
        >
          <p><strong>Código PAF:</strong> {{ p.CodigoPaf }}</p>
          <p><strong>Unidad Menor:</strong> {{ p.NombreUnidadMenor }}</p>
          <p><strong>Codigo de Asignatura:</strong> {{ p.CodigoAsignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.NombreAsignatura }}</p>
          <p><strong>Cantidad de Horas PAF:</strong> {{ p.CantidadHorasPAF }}</p>
          <p><strong>Estado del Proceso:</strong> {{ p.DesEstado }}</p>
        </div>

  <div
    v-if="Number(detalle) === 1"
    v-for="(paf, index) in fichasAgrupadasPAF"
    :key="index"
    class="card"
    :style="{ backgroundColor: fichaSeleccionadaPAF?.CodigoPaf === paf.CodigoPaf ? '#B3E5FC' : coloresPAF[index % coloresPAF.length], margin: '8px', flex: '1 1 48%' }"
    @click="fichaSeleccionadaPAF = { ...fichaSeleccionadaPAF, ...paf }"
  >
    <p><strong>Código PAF:</strong> {{ paf.CodigoPaf }}</p>
    <p><strong>Unidad Menor:</strong> {{ paf.NombreUnidadMenor }}</p>
    <p><strong>Cantidad de Horas PAF:</strong> {{ paf.CantidadHorasPAF }}</p>
    <p><strong>Estado del Proceso:</strong> {{ paf.DesEstado }}</p>
    <div v-for="(asignatura, idx) in paf.Asignaturas" :key="idx">
      <p><strong>Asignatura {{ idx + 1 }}:</strong></p>
      <p><strong>Código de Asignatura:</strong> {{ asignatura.CodigoAsignatura }}</p>
      <p><strong>Nombre de Asignatura:</strong> {{ asignatura.NombreAsignatura }}</p>
    </div>
  </div>
</div>

      <h2 class="sub-title">Horario Asignatura</h2>
      <div v-if="fichasAsignaturas.length > 0" class="flex flex-wrap">
        <div
          v-for="(p, index) in fichasAsignaturas"
          :key="index"
          class="card"
          :style="{ backgroundColor: fichaSeleccionadaAsignatura === p ? '#B3E5FC' : obtenerColorAsignatura(p.codigo_asignatura, p.bloque || '', p.seccion || ''), margin: '8px', flex: '1 1 48%' }"
        >
          <p><strong>Código de Asignatura:</strong> {{ p.codigo_asignatura }}</p>
          <p><strong>Nombre de Asignatura:</strong> {{ p.nombre_asignatura }}</p>
          <p><strong>Sección:</strong> {{ p.seccion }}</p>
          <p><strong>Bloque:</strong> {{ p.bloque }}</p>
          <p><strong>Cupo:</strong> {{ p.Cupo }}</p>
          <p><strong>Semestre:</strong> {{ p.semestre }}</p>
        </div>
      </div>
    </div>
   </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useNuxtApp } from '#app'

const historialSeleccionado = computed(() => persona.value.filter((p) => p.ID !== 0) || null);
const historialNoSeleccionado = computed(() => persona.value.filter((p) => p.ID === 0) || null);

const detalle = ref('');
const mostrarDialogo = ref(false);
const comentario = ref('');

const cancelarEnvio = () => {
  mostrarDialogo.value = false;
  comentario.value = ''; // Limpia el comentario al cancelar
};
const fichasAgrupadasPAF = computed(() => {
  const agrupadas: { [key: number]: { CodigoPaf: number, NombreUnidadMenor: string, CantidadHorasPAF: number, DesEstado: string, Asignaturas: { CodigoAsignatura: string, NombreAsignatura: string}[] } } = {};

  fichasPAF.value.forEach((p) => {
    if (!agrupadas[p.CodigoPaf]) {
      agrupadas[p.CodigoPaf] = {
        CodigoPaf: p.CodigoPaf,
        NombreUnidadMenor: p.NombreUnidadMenor,
        CantidadHorasPAF: p.CantidadHorasPAF,
        DesEstado: p.DesEstado,
        Asignaturas: []
      };
    }

    agrupadas[p.CodigoPaf].Asignaturas.push({
      CodigoAsignatura: p.CodigoAsignatura,
      NombreAsignatura: p.NombreAsignatura
    });
  });

  return Object.values(agrupadas);
});

// Método para limpiar selecciones
const limpiarSelecciones = () => {
  bloquesSeleccionados.value = [];
  fichaSeleccionadaPAF.value = null;
};

// Limpiar al cambiar de página
onBeforeRouteLeave((to, from, next) => {
  limpiarSelecciones();
  next();
});

// Cerrar todos los popups
const cerrarTodosLosBloques = () => {
  Object.keys(visibleBloques.value).forEach((key) => {
    visibleBloques.value[key] = false;
  });
};
const coloresAsignaturas = [
  '#33FF57', '#FF69B4', '#FF33A1', 
  '#8E44AD', '#FF6347', '#ADFF2F', '#4682B4', '#32CD32', '#8A2BE2', 
  '#FF4500', '#20B2AA', '#FF1493', '#FFD700', '#8B0000', '#7FFF00',
  '#00CED1', '#800080', '#FF8C00', '#C71585', '#FFDAB9', '#FFA07A',
  '#D2691E', '#B0C4DE'
];

const asignaturaColoresMap = new Map(); // Mapa para almacenar los colores asignados

const obtenerColorAsignatura = (codigoAsignatura: string, bloqueId: string, seccion: string): string => {
  // Crear una clave única para la combinación
  const key = `${codigoAsignatura}-${bloqueId}-${seccion}`;

  // Verificar si ya se ha asignado un color para esta clave
  if (asignaturaColoresMap.has(key)) {
    return asignaturaColoresMap.get(key);
  }

  // Buscar un color disponible que aún no se haya utilizado
  const availableColors = coloresAsignaturas.filter(color => 
    ![...asignaturaColoresMap.values()].includes(color)
  );

  // Asignar el color
  let color;
  if (availableColors.length > 0) {
    color = availableColors[0]; // Usar el primer color disponible
  } else {
    // Si todos los colores ya están en uso, se pueden reciclar o manejar de otra manera
    color = coloresAsignaturas[asignaturaColoresMap.size % coloresAsignaturas.length];
  }

  // Guardar la asignación en el mapa
  asignaturaColoresMap.set(key, color);

  return color;
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
  persona.value.filter((p) => {
    // Extraemos el año y semestre de semestreSeleccionado.value (formato 2022-01)
    const [añoSeleccionado, semestreSeleccionadoNumber] = semestreSeleccionado.value.split('-').map(Number);

    // Extraemos el año y semestre de p.SemestrePaf (formato 2-22)
    const semestrePafAño = 2000 + parseInt(p.SemestrePaf.split('-')[1], 10); // Convertimos '22' a 2022
    const semestrePafNumber = parseInt(p.SemestrePaf.split('-')[0], 10);

    if(Number(detalle.value) === 1){
      return !historialSeleccionado.value.some(h => 
      h.CodigoPaf === p.CodigoPaf
    )&& semestrePafAño === añoSeleccionado && semestrePafNumber === semestreSeleccionadoNumber
    }
    // Comparamos los semestres y los años
    return !historialSeleccionado.value.some(h => 
      h.ID === p.ID
    ) && semestrePafAño === añoSeleccionado && semestrePafNumber === semestreSeleccionadoNumber;
  })
);

const fichasPAFMatch = computed(() => 
  persona.value.filter((p) => {
    // Extraemos el año y semestre de semestreSeleccionado.value (formato 2022-01)
    const [añoSeleccionado, semestreSeleccionadoNumber] = semestreSeleccionado.value.split('-').map(Number);

    // Extraemos el año y semestre de p.SemestrePaf (formato 2-22)
    const semestrePafAño = 2000 + parseInt(p.SemestrePaf.split('-')[1], 10); // Convertimos '22' a 2022
    const semestrePafNumber = parseInt(p.SemestrePaf.split('-')[0], 10);

    // Comparamos los semestres y los años
    return historialSeleccionado.value.some(h => 
      h.ID === p.ID
    ) && semestrePafAño === añoSeleccionado && semestrePafNumber === semestreSeleccionadoNumber;
  })
);


const fichasAsignaturas = computed(() =>
  persona1.value.filter((p) =>
    !historialSeleccionado.value.some((historial) => {
      const codigosA = historial.CodigoA ? historial.CodigoA.split(" / ") : []; // Dividir CodigoA en un arreglo
      const semestres1 = historial.semestre1 ? historial.semestre1.split(" / ") : []; // Dividir semestre1 en un arreglo
      const secciones = historial.seccion ? historial.seccion.split(" / ") : []; // Dividir seccion en un arreglo
 
      // Comparar cada par de elementos en codigosA y semestres1
      return codigosA.some((codigo, index) => {
        return (
          codigo === p.codigo_asignatura && // Comparar CodigoA
          semestres1[index] === p.semestre &&  // Comparar semestre correspondiente
          secciones[index] === p.seccion // Comparar sección correspondiente
        );
      });
    }) && p.semestre === semestreSeleccionado.value // Comparar con semestreSeleccionado
  )
);

const fichasAsignaturasNo = computed(() =>
  persona1.value.filter((p) =>
    historialSeleccionado.value.some((historial) => {
      const codigosA = historial.CodigoA ? historial.CodigoA.split(" / ") : []; // Dividir CodigoA en un arreglo
      const semestres1 = historial.semestre1 ? historial.semestre1.split(" / ") : []; // Dividir semestre1 en un arreglo
      const secciones = historial.seccion ? historial.seccion.split(" / ") : []; // Dividir seccion en un arreglo

      // Comparar cada par de elementos en codigosA y semestres1
      return codigosA.some((codigo, index) => {
        return (
          codigo === p.codigo_asignatura && // Comparar CodigoA
          semestres1[index] === p.semestre &&  // Comparar semestre correspondiente
          secciones[index] === p.seccion // Comparar sección correspondiente
        );
      });
    }) && p.semestre === semestreSeleccionado.value // Comparar con semestreSeleccionado
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
  DesEstado: string;
  bloque: string;
  cupo: number;
  seccion: string;
  CodigoA: string;
  bloques: string;
  semestre1: string;
  SemestrePaf: string;
  CantidadHorasPAF: number;
  Run: string;
  Asignaturas?: { CodigoAsignatura: string; NombreAsignatura: string }[]; // Add the 'Asignaturas' property
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
  Seccion: string;
}

const persona = ref<Persona[]>([])
const persona1 = ref<Horario[]>([])

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
const semestres = computed(() => {
  // Ordenar los semestres y eliminar duplicados
  const semestresOrdenados = [...new Set(persona1.value.map(p => p.semestre))]
    .sort((a, b) => (a || '').localeCompare(b || '')); // Comparación lexicográfica adecuada para "AAAA-MM"

  return semestresOrdenados; // O retorna [ultimoSemestre] si solo necesitas el último
});
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
            


      // Filtrar por día y módulo
      const bloques = p.bloque.split('-');
      return bloques.some((b) => {
        const diaBloque = b.charAt(0);
        const moduloBloque = b.slice(1);
        return inicialDia[dia] === diaBloque && parseInt(moduloBloque) === modulo;
      });
    })
    .map((p) => {
      // Determinar color basado en si el bloque es excluido
      const esExcluido = fichasAsignaturasNo.value.some((ficha) => 
        ficha.codigo_asignatura === p.codigo_asignatura && 
        ficha.semestre === p.semestre && 
        ficha.seccion === p.seccion
      );
      const color = esExcluido 
        ?  '#FFCC80' // Cambiar esto por el color deseado para bloques excluidos
        : obtenerColorAsignatura(p.codigo_asignatura, p.bloque || '', p.seccion || '');

      return {
      nombre: p.nombre_asignatura,
      seccion: p.seccion,
      color: color,
      ID: p.ID,
      cupo: p.Cupo,
      codigo_asignatura: p.codigo_asignatura,
      tipo: p.tipo,
      run: p.run,
      semestre: p.semestre,
      bloque: (p.bloque ?? '').split('-')
        .map(b => `${b.charAt(0)}${b.slice(1)}`) // Mapear para obtener "V3", "M2", etc.
        .join('-')
    };
    });
};

const obtenerDatosPersona = async (semestreGuardado: string | null) => {
  try {

    // NO DEVUELÑVE LAS PAF LISTAS

    const response = await $axios.get(`/api/paf-en-linea/pipelsoft/obtenerContratos/mostrarTodo/${run.value}`);
    const rut = typeof run.value === 'string' ? run.value.replace(/^0+(?=\d+-)/, '') : '';
    const response1 = await $axios.get(`/api/paf-en-linea/profesorDB/${rut.slice(0, -2)}`);
    persona1.value = response1.data;

    persona.value = response.data.map((item: any) => {
      const bloquesArray = item.HistorialPafData.Bloque || []; // Asegurar que Bloque sea un arreglo (vacío si es null o undefined)

      // Verificar si el arreglo no está vacío antes de hacer el map
      const bloque = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.bloques).join(" / ") : "";
      const CodigoA = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.codigoAsignatura).join(" / ") : "";
      const cupo = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.cupos).join(" / ") : "";
      const seccion = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.seccion).join(" / ") : "";
      const semestre1 = bloquesArray.length > 0 ? bloquesArray.map((bloque: any) => bloque.semestre).join(" / ") : "";

      return {
        CodigoPaf: item.PipelsoftData.IdPaf,
        CodigoAsignatura: item.PipelsoftData.CodigoAsignatura,
        Nombres: item.PipelsoftData.Nombres,
        NombreAsignatura: item.PipelsoftData.NombreAsignatura,
        PrimerApellido: item.PipelsoftData.PrimerApp,
        SegundoApellido: item.PipelsoftData.SegundoApp,
        NombreUnidadMenor: item.PipelsoftData.NombreUnidadMenor,
        SemestrePaf: item.PipelsoftData.Semestre,
        CantidadHorasPAF: item.PipelsoftData.CantidadHorasPaf,
        Run: item.PipelsoftData.RunEmpleado,
        DesEstado: item.PipelsoftData.DesEstado,
        Bloque: response1.data.bloque,
        Cupo: response1.data.cupo,
        Seccion: response1.data.seccion,
        Semestre: response1.data.semestre,
        ID: item.HistorialPafData.ID,
        semestre: item.HistorialPafData.semestre,
        codigo: item.HistorialPafData.CodigoAsignatura,
        paf: item.HistorialPafData.CodigoPaf,
        bloque, // Agregar las cadenas combinadas como propiedades
        CodigoA,
        cupo,
        seccion,
        semestre1,
      };
    });

    // Identificar el semestre más reciente
    const semestresDisponibles = persona1.value.map(p => p.semestre).filter(Boolean);
    const semestreReciente = semestresDisponibles.sort((a, b) => {
      if (a && b) {
        return a > b ? -1 : 1;
      }
      return 0;
    })[0];

    if(semestreGuardado !== null) {
      semestreSeleccionado.value = semestreGuardado;
    }
    // Establecer el semestre más reciente como seleccionado
     else if (semestreReciente) {
      semestreSeleccionado.value = semestreReciente;
    }
  } catch (error) {
    console.error("Error al obtener los datos:", error);
  }
};
const enviarSeleccion1 = async () => {
  if (!fichaSeleccionadaPAF || !fichaSeleccionadaAsignatura) {
    alert('Por favor selecciona una ficha de PAF y una de asignatura.');
    return;
  }
  mostrarDialogo.value = true;
};

const enviarSeleccion = async () => {
  if (!fichaSeleccionadaPAF || !fichaSeleccionadaAsignatura) {
    alert('Por favor selecciona una ficha de PAF y una de asignatura.');
    return;
  }
  if(comentario.value === '') {
    comentario.value = 'Sin comentarios';
  }
  try {
    const codigoPAF = fichaSeleccionadaPAF.value?.CodigoPaf; // Ajustar si el código es otro campo
    const bloquesSeleccionadosString = computed(() => bloquesSeleccionados.value.join(','))
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
    bloque: bloques.map(b => `${semestre} ${codigoAsignatura} ${seccion} ${cupo} ${bloques.join("-")}`)
  };
});
const result = resultado.map(item => {
  return item.bloque[0];
});

    const data = resultado[0];
    data.bloque = result;
    let Asignatura = fichaSeleccionadaPAF.value?.CodigoAsignatura;
    if (!Asignatura){
        Asignatura = fichaSeleccionadaPAF.value?.Asignaturas?.[0]?.CodigoAsignatura;        ;
    }
    await $axios.post(`/api/paf-en-linea/historial/post/${codigoPAF}/${Asignatura}/${comentario.value}`, data);
    alert('Datos enviados correctamente.');
    mostrarDialogo.value = false;
    await $axios.put(`/api/paf-en-linea/historial/${codigoPAF}/actualizarBanderaAceptacion`, {
      nuevaBanderaAceptacion: 1,
    });
  } catch (error) {
    console.error('Error al enviar los datos:', error);
    alert('Hubo un error al enviar los datos.');
  }
  obtenerDatosPersona(semestreSeleccionado.value);
};


const volver = () => {
  router.go(-1)
}

onMounted(() => {
  detalle.value = sessionStorage.getItem('detalle') || '';
  obtenerDatosPersona(null)
})
const coloresPAF = [
  '#FFCDD2', '#F8BBD0', '#E1BEE7', '#D1C4E9', '#C5CAE9', // Tonos originales suaves
  '#81D4FA', '#4FC3F7', '#29B6F6', '#0288D1', '#0277BD', // Tonos azules
  '#DCE775', '#C0CA33', '#9CCC65', '#8BC34A', '#558B2F', // Tonos verdes y amarillos
  '#FFC107', '#FF9800', '#FF7043', '#795548', '#607D8B'  // Tonos cálidos y neutros
];

</script>

<style scoped>
.container {
  display: grid;
  grid-template-columns: 70% 30%;
  gap: 10px;
  max-width: 100%;
}

.tabla-container {
  max-width: 100%;
}

.fichas-container {
  max-width: 150%;
}

.w-full {
  width: 100%;
}
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
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* Fondo semitransparente */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999; /* Asegúrate de que esté por encima de otros elementos */
}

.modal {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  width: 400px; /* Ajusta el tamaño según lo necesites */
  z-index: 10000; /* Asegúrate de que el modal esté encima del overlay */
}

.textarea-comentario {
  width: 100%;
  height: 100px;
}

.modal-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.cancel-button, .confirm-button {
  padding: 10px 20px;
  cursor: pointer;
}

.cancel-button {
  background-color: #f44336;
  color: white;
}

.confirm-button {
  background-color: #4CAF50;
  color: white;
}

</style>