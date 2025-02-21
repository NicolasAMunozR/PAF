<template>
  <div class="mt-4 ml-4">
    <button @click="volver" class="volver-button">
       Volver
    </button>
  </div>
  <div class="form-container">
    <h1>Formulario de Contrato</h1>
    <form @submit.prevent="submitForm">
      <div class="form-group">
        <label for="unidadMayor">Unidad Mayor</label>
        <input v-model="formData.UnidadMayor" type="text" id="unidadMayor" class="form-control" />
      </div>
      <div class="form-group">
        <label for="unidadMenor">Unidad Menor</label>
        <input v-model="formData.UnidadMenor" type="text" id="unidadMenor" class="form-control" />
      </div>
      <div class="form-group">
        <label for="numeroCentroDeCostos">Número Centro de Costos</label>
        <input v-model="formData.NumeroCentroDeCostos" type="text" id="numeroCentroDeCostos" class="form-control" />
      </div>
      <div class="form-group">
        <label for="celulaIdentidad">Cédula de Identidad</label>
        <input v-model="formData.CelulaIdentidad" type="text" id="celulaIdentidad" class="form-control" />
      </div>
      <div class="form-group">
        <label for="lugarNacimiento">Lugar de Nacimiento</label>
        <input v-model="formData.LugarNacimiento" type="text" id="lugarNacimiento" class="form-control" />
      </div>
      <div class="form-group">
        <label for="fechaYHoraNacimiento">Fecha y Hora de Nacimiento</label>
        <input v-model="formData.FechaYHoraNacimiento" type="datetime-local" id="fechaYHoraNacimiento" class="form-control" />
      </div>
      <div class="form-group">
        <label for="apellidoP">Apellido Paterno</label>
        <input v-model="formData.ApellidoP" type="text" id="apellidoP" class="form-control" />
      </div>
      <div class="form-group">
        <label for="apellidoM">Apellido Materno</label>
        <input v-model="formData.ApellidoM" type="text" id="apellidoM" class="form-control" />
      </div>
      <div class="form-group">
        <label for="nombres">Nombres</label>
        <input v-model="formData.Nombres" type="text" id="nombres" class="form-control" />
      </div>
      <div class="form-group">
        <label for="nacionalidad">Nacionalidad</label>
        <input v-model="formData.Nacionalidad" type="text" id="nacionalidad" class="form-control" />
      </div>
      <div class="form-group">
        <label for="domicilio">Domicilio</label>
        <input v-model="formData.Domicio" type="text" id="domicilio" class="form-control" />
      </div>
      <div class="form-group">
        <label for="correo">Correo</label>
        <input v-model="formData.Correo" type="email" id="correo" class="form-control" />
      </div>
      <div class="form-group">
        <label for="titulo">Título</label>
        <input v-model="formData.Titulo" type="text" id="titulo" class="form-control" />
      </div>
      <div class="form-group">
        <label for="institucion">Institución</label>
        <input v-model="formData.Institucion" type="text" id="institucion" class="form-control" />
      </div>
      <div class="form-group">
        <label for="fechaObtencion">Fecha de Obtención</label>
        <input v-model="formData.FechaObtencion" type="date" id="fechaObtencion" class="form-control" />
      </div>
      <div class="form-group">
        <label for="numeroSemestre">Número de Semestre</label>
        <input v-model="formData.NumeroSemestre" type="number" id="numeroSemestre" class="form-control" />
      </div>
      <div class="form-group">
        <label for="gradoAcademico">Grado Académico</label>
        <input v-model="formData.GradoAcademico" type="text" id="gradoAcademico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="institucionGradoAcademico">Institución del Grado Académico</label>
        <input v-model="formData.InstitucionGradoAcademico" type="text" id="institucionGradoAcademico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="fechaObtencionGradoAcademico">Fecha de Obtención del Grado Académico</label>
        <input v-model="formData.FechaObtencionGradoAcademico" type="date" id="fechaObtencionGradoAcademico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="tipoIngreso">Tipo de Ingreso</label>
        <input v-model="formData.TipoIngreso" type="text" id="tipoIngreso" class="form-control" />
      </div>
      <div class="form-group">
        <label for="cargo">Cargo</label>
        <input v-model="formData.Cargo" type="text" id="cargo" class="form-control" />
      </div>
      <div class="form-group">
        <label for="nivel">Nivel</label>
        <input v-model="formData.Nivel" type="text" id="nivel" class="form-control" />
      </div>
      <div class="form-group">
        <label for="grado">Grado</label>
        <input v-model="formData.Grado" type="text" id="grado" class="form-control" />
      </div>
      <div class="form-group">
        <label for="rango">Rango</label>
        <input v-model="formData.Rango" type="text" id="rango" class="form-control" />
      </div>
      <div class="form-group">
        <label for="funcion">Función</label>
        <input v-model="formData.Funcion" type="text" id="funcion" class="form-control" />
      </div>
      <div class="form-group">
        <label for="jerarquia">Jerarquía</label>
        <input v-model="formData.Jerarquia" type="text" id="jerarquia" class="form-control" />
      </div>
      <div class="form-group">
        <label for="asignatura">Asignatura</label>
        <input v-model="formData.Asignatura" type="text" id="asignatura" class="form-control" />
      </div>
      <div class="form-group">
        <label for="numeroHoras">Número de Horas</label>
        <input v-model="formData.NumeroHoras" type="number" id="numeroHoras" class="form-control" />
      </div>
      <div class="form-group">
        <label for="categoria">Categoría</label>
        <input v-model="formData.Categoria" type="text" id="categoria" class="form-control" />
      </div>
      <div class="form-group">
        <label for="calidad">Calidad</label>
        <input v-model="formData.Calidad" type="text" id="calidad" class="form-control" />
      </div>
      <div class="form-group">
        <label for="lugarDesempeño">Lugar de Desempeño</label>
        <input v-model="formData.LugarDesempeño" type="text" id="lugarDesempeño" class="form-control" />
      </div>
      <div class="form-group">
        <label for="cargoOtroPublico">Cargo en Otro Público</label>
        <input v-model="formData.CargoOtroPublico" type="text" id="cargoOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="gradoOtroPublico">Grado en Otro Público</label>
        <input v-model="formData.GradoOtroPublico" type="text" id="gradoOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="nivelOtroPublico">Nivel en Otro Público</label>
        <input v-model="formData.NivelOtroPublico" type="text" id="nivelOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="rangoOtroPublico">Rango en Otro Público</label>
        <input v-model="formData.RangoOtroPublico" type="text" id="rangoOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="numeroHorasOtroPublico">Número de Horas en Otro Público</label>
        <input v-model="formData.NumeroHorasOtroPublico" type="number" id="numeroHorasOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="calidadOtroPublico">Calidad en Otro Público</label>
        <input v-model="formData.CalidadOtroPublico" type="text" id="calidadOtroPublico" class="form-control" />
      </div>
      <div class="form-group">
        <label for="fechaInicioContrato">Fecha de Inicio del Contrato</label>
        <input v-model="formData.FechaInicioContrato" type="date" id="fechaInicioContrato" class="form-control" />
      </div>
      <div class="form-group">
        <label for="fechaFinContrato">Fecha de Fin del Contrato</label>
        <input v-model="formData.FechaFinContrato" type="date" id="fechaFinContrato" class="form-control" />
      </div>
      <div class="form-group button-container">
        <button type="submit" class="btn btn-primary">Enviar</button>
      </div>
      <!--<div v-if="downloadUrl" class="download-container">
      <a :href="downloadUrl" download="contrato.pdf" class="btn btn-secondary">Descargar Contrato</a>
    </div>-->
    </form>
  </div>
  <div v-if="downloadUrl" class="preview-container">
      <iframe :src="downloadUrl" width="80%" height="1000px"></iframe>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const { $axios } = useNuxtApp() as unknown as { $axios: typeof import('axios').default }

const router = useRouter();

const volver = () => {
  router.go(-1);
};

const formData = ref({
  UnidadMayor: "Facultad de Ingeniería",
  UnidadMenor: "Departamento de Informática",
  NumeroCentroDeCostos: "12345",
  CelulaIdentidad: "12345678-9",
  LugarNacimiento: "Santiago",
  FechaYHoraNacimiento: "1990-01-01T00:00",
  ApellidoP: "Pérez",
  ApellidoM: "Gómez",
  Nombres: "Juan Carlos",
  Nacionalidad: "Chilena",
  Domicio: "Av. Siempre Viva 742",
  Correo: "juan.perez@example.com",
  Titulo: "Ingeniero en Informática",
  Institucion: "Universidad XYZ",
  FechaObtencion: "2015-12-15",
  NumeroSemestre: "10",
  GradoAcademico: "Magíster en TI",
  InstitucionGradoAcademico: "Universidad ABC",
  FechaObtencionGradoAcademico: "2018-06-20",
  TipoIngreso: "Contrato",
  Cargo: "Profesor Asociado",
  Nivel: "2",
  Grado: "3",
  Rango: "Senior",
  Funcion: "Docencia",
  Jerarquia: "Profesor",
  Asignatura: "Programación Avanzada",
  NumeroHoras: "20",
  Categoria: "Académico",
  Calidad: "Titular",
  LugarDesempeño: "Campus Norte",
  CargoOtroPublico: "",
  GradoOtroPublico: "",
  NivelOtroPublico: "",
  RangoOtroPublico: "",
  NumeroHorasOtroPublico: 0,
  CalidadOtroPublico: "",
  FechaInicioContrato: "2023-03-01",
  FechaFinContrato: "2024-02-28"
})

const downloadUrl = ref<string | null>(null)

const submitForm = async () => {
        try {
    console.log(formData.value)
    await $axios.post('/api/paf-en-linea/archivo/generarContratos/SinDatos/', formData.value)

    const response = await $axios.get('/api/paf-en-linea/archivo/archivo/3/descargar', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/pdf' }))
    downloadUrl.value = url
  } catch (error) {
    console.error('Error al generar o descargar el contrato:', error)
  }
}
</script>

<style scoped>
.form-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f9fafb;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-control {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.btn {
  display: inline-block;
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-primary {
  background-color: #007bff;
}

.btn-primary:hover {
  background-color: #0056b3;
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

.download-container {
  text-align: center;
  margin-top: 20px;
}

.preview-container {
  display: flex;
  text-align: center;
  justify-content: center;
  margin-top: 20px;
  max-height: 10000px;
  max-width: 10000px;
}

.button-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>