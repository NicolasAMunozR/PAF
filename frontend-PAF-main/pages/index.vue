<template>
    <div class="login-page">
      <!-- Barra superior con título -->
      <header class="header-bar">
        <h1>PAF</h1>
      </header>
  
      <!-- Contenido principal -->
      <div class="content">
        <h2>Inicio de Sesión</h2>
  
        <!-- Selector de roles como secciones -->
        <div class="role-selector">
          <button
            v-for="option in roleOptions"
            :key="option.value"
            :class="{ active: selectedRole === option.value }"
            @click="selectRole(option.value)"
          >
            {{ option.label }}
          </button>
        </div>
  
        <!-- Formulario de inicio de sesión -->
        <form @submit.prevent="handleLogin" v-if="selectedRole">
          <div class="form-group">
            <label for="run">RUN:</label>
            <input
              type="text"
              id="run"
              v-model="run"
              placeholder="Ingresa tu RUN"
              required
            />
          </div>
          <div class="form-group">
            <label for="email">Correo Electrónico:</label>
            <input
              type="email"
              id="email"
              v-model="email"
              placeholder="Ingresa tu correo"
              required
            />
          </div>
          <p class="role-info">Iniciarás sesión como: <strong>{{ selectedRole }}</strong></p>
          <button type="submit">Ingresar</button>
        </form>
  
        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    layout: false, // Oculta el layout para esta página
    data() {
      return {
        run: "", // RUN ingresado
        email: "", // Correo electrónico ingresado
        selectedRole: "", // Rol seleccionado
        errorMessage: "", // Mensaje de error
        roleOptions: [
          { value: "profesor", label: "Profesor" },
          { value: "personal-dei", label: "Personal del DEI" },
          { value: "encargado", label: "Encargado" },
        ],
      };
    },
    methods: {
      selectRole(role) {
        this.selectedRole = role;
        this.errorMessage = ""; // Limpia el mensaje de error al cambiar de rol
      },
      handleLogin() {
        if (this.run && this.email && this.selectedRole) {
          // Redirigir según el rol seleccionado
          if (this.selectedRole === "profesor") {
            this.$router.push(`/profesorPAF?run=${this.run}`);
          } else if (this.selectedRole === "personal-dei") {
            this.$router.push("/personas");
          } else if (this.selectedRole === "encargado") {
            this.$router.push("/seguimientoPAF");
          }
        } else {
          this.errorMessage = "Por favor, completa todos los campos.";
        }
      },
    },
  };
  </script>
  
  <style scoped>
  /* Estilo general de la página */
  .login-page {
    max-width: 400px;
    margin: 0 auto;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #f9f9f9;
    overflow: hidden; /* Para evitar desbordes */
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  }
  
  /* Barra superior */
  .header-bar {
    background: var(--primary-color, #4A90E2);
    color: #fff;
    padding: 10px;
    text-align: center;
    font-size: 1.5rem;
    font-weight: bold;
    border-bottom: 4px solid var(--accent-color, #50E3C2);
    text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.2);
  }
  
  /* Contenido principal */
  .content {
    padding: 20px;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 20px;
  }
  
  /* Selector de roles */
  .role-selector {
    display: flex;
    justify-content: space-around;
    margin-bottom: 20px;
  }
  
  .role-selector button {
    padding: 10px 15px;
    background-color: #f0f0f0;
    border: 1px solid #ccc;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  .role-selector button.active {
    background-color: #007bff;
    color: #fff;
  }
  
  .role-selector button:hover:not(.active) {
    background-color: #e0e0e0;
  }
  
  /* Estilo del formulario */
  .form-group {
    margin-bottom: 15px;
  }
  
  label {
    display: block;
    font-weight: bold;
    margin-bottom: 5px;
  }
  
  input {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  button[type="submit"] {
    width: 100%;
    padding: 10px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button[type="submit"]:hover {
    background-color: #0056b3;
  }
  
  /* Información del rol seleccionado */
  .role-info {
    margin: 10px 0;
    font-style: italic;
    color: #555;
    text-align: center;
  }
  
  /* Mensaje de error */
  .error {
    color: red;
    margin-top: 10px;
    text-align: center;
  }
  </style>
  