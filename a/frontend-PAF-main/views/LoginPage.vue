<template>
    <div class="login-page">
      <!-- Barra superior con imagotipo y título -->
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
        <form @submit.prevent="handleLogin" v-if="selectedRole" class="login-form">
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
          <p class="role-info">
            Iniciarás sesión como: <strong>{{ selectedRole }}</strong>
          </p>
          <button type="submit">Ingresar</button>
        </form>
  
        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    layout: false,
    data() {
      return {
        run: "",
        selectedRole: "",
        errorMessage: "",
        roleOptions: [
          //{ value: "profesor", label: "Profesor" },
          { value: "personal-dei", label: "Personal del Dei" },
          { value: "encargado", label: "Encargado" },
        ],
      };
    },
    methods: {
      selectRole(role) {
        this.selectedRole = role;
        this.errorMessage = "";
      },
      async handleLogin() {
        if (this.run && this.selectedRole) {
          sessionStorage.setItem("rut", this.run); // Guardar en sesión
          try {
            // Accede a $axios desde el contexto del componente
            const response = await this.$axios.get(`/usuario/rut/${this.run}`);
            if (!response.data) {
              this.errorMessage = "Usuario no encontrado.";
              return;
            }
            // Redirige según el rol
            //if (response.data.rol === "profesor") {
            //  this.$router.push(`/profesorPAF?run=${this.run}`);
            //} else 
            if (response.data.Rol === "encargado" && this.selectedRole === "encargado") {
              if (response.data.UnidadMayor == "RECTORIA" || response.data.UnidadMayor == "VR ACADEMICA") {
                this.$router.push("principal/seguimientoPAF");
              } else {
                this.$router.push(`principal/unidadMayorPAF?UnidadMayor=${response.data.UnidadMayor}`);
              }
            } else if (response.data.Rol === "personal-dei" && this.selectedRole === "personal-dei") {
              this.$router.push("/personas");
            } else {
              this.errorMessage = "El rol seleccionado no coincide con el usuario.";
            }
          } catch (error) {
            this.errorMessage = "Hubo un error al iniciar sesión.";
            console.error(error);
          }
        } else {
          this.errorMessage = "Por favor, completa todos los campos.";
        }
      },
    },
  };
  </script>
  
  
  <style scoped>
  /* Página principal */
  .login-page {
    max-width: 500px;
    margin: 20px auto;
    border: 1px solid #394049;
    border-radius: 8px;
    background-color: #f9f9f9;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    text-align: center;
  }
  
  /* Barra superior */
  .header-bar {
    background-color: #EA7600;
    color: #fff;
    padding: 10px;
    font-family: "Bebas Neue Pro", sans-serif;
    font-size: 1.5rem;
    font-weight: bold;
  }
  
  /* Contenido principal */
  .content {
    padding: 20px;
  }
  
  /* Selector de roles */
  .role-selector {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-bottom: 20px;
  }
  
  .role-selector button {
    background-color: #f0f0f0;
    border: 1px solid #394049;
    border-radius: 4px;
    padding: 10px 15px;
    font-family: "Helvetica Neue LT", sans-serif;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  .role-selector button.active {
    background-color: #00A499;
    color: #fff;
  }
  
  /* Formulario */
  .login-form {
    margin-top: 20px;
  }
  
  .form-group {
    margin-bottom: 15px;
  }
  
  .form-group label {
    display: block;
    font-family: "Helvetica Neue LT", sans-serif;
    font-size: 0.9rem;
    text-align: left;
    margin-bottom: 5px;
  }
  
  .form-group input {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.9rem;
  }
  
  button[type="submit"] {
    background-color: #EA7600;
    border: none;
    color: white;
    font-family: "Bebas Neue Pro", sans-serif;
    padding: 10px 20px;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  button[type="submit"]:hover {
    background-color: #C8102E;
  }
  
  /* Errores */
  .error {
    color: red;
    margin-top: 10px;
    font-size: 0.9rem;
  }
  </style>
  